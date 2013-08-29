package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"gobus"
	"io/ioutil"
	"logger"
	"model"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ProcessTimeout = 60 * time.Second
	NetworkTimeout = 30 * time.Second
)

var internalError = errors.New("internal error")

type ErrorType string
type Error struct {
	Type    ErrorType
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("(%s)%s", e.Type, e.Message)
}

type Warning struct {
	Type ErrorType              `json:"type"`
	Vars map[string]interface{} `json:"message"`
}

func (w Warning) Error() string {
	return fmt.Sprintf("%s(%+v)", w.Type, w.Vars)
}

const (
	IDENTITY_NOT_FOUND ErrorType = "identity_not_found"
	CROSS_NOT_FOUND              = "cross_not_found"
	CROSS_NOT_MODIFIED           = "cross_not_modified"
	CROSS_FORBIDDEN              = "cross_forbidden"
	CROSS_ERROR                  = "cross_error"
	NOT_AUTHORIZED               = "not_authorized"
)

var client *http.Client

func init() {
	tran := &http.Transport{
		Proxy:               nil,
		Dial:                dial,
		TLSClientConfig:     nil,
		DisableKeepAlives:   true,
		DisableCompression:  false,
		MaxIdleConnsPerHost: 0,
	}
	client = &http.Client{
		Transport:     tran,
		CheckRedirect: nil,
		Jar:           nil,
	}
}

func dial(net_, addr string) (net.Conn, error) {
	conn, err := net.Dial(net_, addr)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(time.Now().Add(NetworkTimeout))
	return conn, nil
}

type Platform struct {
	dispatcher *gobus.Dispatcher
	config     *model.Config
	replacer   *strings.Replacer
}

func NewPlatform(config *model.Config) (*Platform, error) {
	table, err := gobus.NewTable(config.Dispatcher)
	if err != nil {
		return nil, err
	}
	dispatcher := gobus.NewDispatcher(table)
	return &Platform{
		dispatcher: dispatcher,
		config:     config,
		replacer:   strings.NewReplacer(`"place":{},`, "", `"time":{"begin_at":{}},`, ""),
	}, nil
}

func (p *Platform) Send(to model.Recipient, text string) (string, int64, bool, error) {
	url := fmt.Sprintf("http://%s:%d/v3/poster/message/%s/%s", p.config.ExfeService.Addr, p.config.ExfeService.Port, to.Provider, to.ExternalUsername)
	resp, err := Http("POST", url, "plain/text", []byte(text))
	if err != nil {
		logger.DEBUG("post %s error: %s with %s", url, err, text)
		return "", 0, false, internalError
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.DEBUG("read %s error: %s with %s", url, err, text)
		return "", 0, false, internalError
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return string(body), 0, true, nil
	case http.StatusAccepted:
		ontimeStr := resp.Header.Get("Ontime")
		ontime, err := strconv.ParseInt(ontimeStr, 10, 64)
		if err != nil {
			logger.DEBUG("can't parse ontime %s: %s", ontimeStr, err)
			return "", 0, false, err
		}
		defaultOK := resp.Header.Get("Default")
		return string(body), ontime, defaultOK == "true", nil
	}
	return "", 0, false, fmt.Errorf("(%s)%s", resp.Status, string(body))
}

func (p *Platform) FindIdentity(identity model.Identity) (model.Identity, error) {
	b, err := json.Marshal(identity)
	if err != nil {
		logger.ERROR("encode identity error: %s with %+v", err, identity)
		return identity, internalError
	}
	url := fmt.Sprintf("%s/v3/bus/revokeidentity", p.config.SiteApi)
	resp, err := Http("POST", url, "application/json", b)
	reader, err := HttpResponse(resp, err)

	if err != nil {
		switch resp.StatusCode {
		case 404:
			return identity, Error{IDENTITY_NOT_FOUND, err.Error()}
		}
		logger.ERROR("post %s error: %s with %s", url, err, string(b))
		return identity, internalError
	}

	defer reader.Close()
	var ret struct {
		Data model.Identity `json:"data"`
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&ret)
	if err != nil {
		logger.ERROR("decode %s error: %s", url, err, string(b))
		return identity, internalError
	}
	return ret.Data, nil
}

func (p *Platform) GetConversation(exfeeId int64, updatedAt string, clear bool, direction string, quantity int) ([]model.Post, error) {
	query := make(url.Values)
	query.Set("updated_at", updatedAt)
	query.Set("clear", fmt.Sprintf("%v", clear))
	query.Set("direction", direction)
	query.Set("quantity", fmt.Sprintf("%d", quantity))
	url := fmt.Sprintf("%s/v3/bus/conversation/%d?%s", p.config.SiteApi, exfeeId, query.Encode())

	resp, err := HttpResponse(Http("GET", url, "", nil))
	if err != nil {
		logger.ERROR("get %s error: %s", url)
		return nil, internalError
	}
	defer resp.Close()

	var ret struct {
		Data []model.Post `json:"data"`
	}
	decoder := json.NewDecoder(resp)
	err = decoder.Decode(&ret)
	if err != nil {
		logger.ERROR("decode %s error: %s", url, err)
		return nil, internalError
	}
	return ret.Data, nil
}

func (p *Platform) FindCross(id int64, query url.Values) (model.Cross, error) {
	url := fmt.Sprintf("%s/v3/bus/Crosses/%d?", p.config.SiteApi, id)
	if len(query) > 0 {
		url += query.Encode()
	}
	resp, err := Http("GET", url, "", nil)
	reader, err := HttpResponse(resp, err)

	var ret struct {
		Data model.Cross `json:"data"`
	}
	if err != nil {
		switch resp.StatusCode {
		case 304:
			logger.ERROR("get %s error: %s", url, err)
			return ret.Data, Error{CROSS_NOT_MODIFIED, err.Error()}
		case 403:
			logger.ERROR("get %s error: %s", url, err)
			return ret.Data, Error{CROSS_FORBIDDEN, err.Error()}
		case 404:
			logger.ERROR("get %s error: %s", url, err)
			return ret.Data, Error{CROSS_NOT_FOUND, err.Error()}
		}
		logger.ERROR("get %s error: %s", url, err)
		return ret.Data, internalError
	}

	defer reader.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ret)
	if err != nil {
		logger.ERROR("decode %s error: %s", url, err)
		return ret.Data, internalError
	}
	return ret.Data, nil
}

func (p *Platform) UploadPhoto(photoxID string, photos []model.Photo) error {
	b, err := json.Marshal(photos)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/v3/bus/addphotos/%s", p.config.SiteApi, photoxID)
	resp, err := HttpResponse(Http("POST", url, "application/json", b))
	if err != nil {
		return err
	}
	defer resp.Close()
	return nil
}

func (p *Platform) BotCrossGather(cross model.Cross) (model.Cross, error) {
	b, err := json.Marshal(cross)
	if err != nil {
		logger.ERROR("encode cross error: %s with %+v", err, cross)
		return model.Cross{}, internalError
	}
	b = []byte(p.replacer.Replace(string(b)))

	u := fmt.Sprintf("%s/v3/bus/gather", p.config.SiteApi)
	resp, err := Http("POST", u, "application/json", b)
	reader, err := HttpResponse(resp, err)

	if err != nil {
		switch resp.StatusCode {
		case 400:
			return model.Cross{}, Error{CROSS_ERROR, err.Error()}
		}
		logger.ERROR("post %s error: %s, with %s", u, err, string(b))
		return model.Cross{}, internalError
	}

	defer reader.Close()
	var ret struct {
		Data    model.Cross `json:"data"`
		Warning Warning     `json:"warning"`
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&ret)
	if err != nil {
		logger.ERROR("parse %s error: %s with %s", u, err, string(b))
		return model.Cross{}, internalError
	}

	if resp.StatusCode == 200 {
		return ret.Data, nil
	}
	return ret.Data, ret.Warning
}

func (p *Platform) BotCrossUpdate(objectType, objectId string, cross interface{}, by model.Identity) error {
	arg := make(map[string]interface{})
	arg[objectType] = objectId
	if cross != nil {
		arg["cross"] = cross
		arg["by_identity"] = by
	}

	b, err := json.Marshal(arg)
	if err != nil {
		logger.ERROR("encoding error: %s with %+v", err, arg)
		return internalError
	}
	b = []byte(p.replacer.Replace(string(b)))

	u := fmt.Sprintf("%s/v3/bus/xupdate", p.config.SiteApi)
	resp, err := Http("POST", u, "application/json", b)
	reader, err := HttpResponse(resp, err)
	if err != nil {
		switch resp.StatusCode {
		case 400:
			return Error{NOT_AUTHORIZED, err.Error()}
		case 404:
			return Error{CROSS_NOT_FOUND, err.Error()}
		}
		logger.ERROR("post %s error: %s with %s", u, err, string(b))
		return internalError
	}

	defer reader.Close()
	var ret struct {
		Warning Warning `json:"warning"`
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&ret)
	if err != nil {
		logger.ERROR("decode %s error: %s with %s", u, err, string(b))
		return internalError
	}
	if resp.StatusCode == 200 {
		return nil
	}
	return ret.Warning
}

func (p *Platform) BotPostConversation(from, post, createdAt string, exclude []*mail.Address, to, id string) error {
	u := fmt.Sprintf("%s/v3/bus/postconversation", p.config.SiteApi)
	params := make(url.Values)
	params.Add(to, id)
	params.Add("content", post)
	params.Add("external_id", from)
	params.Add("time", createdAt)
	params.Add("provider", "email")
	ex := make([]string, len(exclude))
	for i, addr := range exclude {
		ex[i] = fmt.Sprintf("%s@email", addr.Address)
	}
	params.Add("exclude", strings.Join(ex, ","))

	resp, err := HttpClient.PostForm(u, params)
	reader, err := HttpResponse(resp, err)
	if err != nil {
		logger.ERROR("post %s error: %s with %s", u, err, params.Encode())
		return internalError
	}
	defer reader.Close()

	return nil
}

func (p *Platform) GetIdentity(identities []model.Identity) ([]model.Identity, error) {
	arg := map[string]interface{}{
		"identities": identities,
	}
	b, err := json.Marshal(arg)
	if err != nil {
		logger.ERROR("encode error: %s with %+v", err, arg)
		return nil, err
	}
	u := fmt.Sprintf("%s/v2/identities/get", p.config.SiteApi)
	reader, err := HttpResponse(Http("POST", u, "application/json", b))
	if err != nil {
		logger.ERROR("post %s error: %s with %s", u, err, string(b))
		return nil, internalError
	}

	defer reader.Close()
	var ret struct {
		Meta struct {
			Code        int    `json:"code"`
			ErrorDetail string `json:"errorDetail"`
		} `json:"meta"`
		Response struct {
			Identities []model.Identity `json:"identities"`
		} `json:"response"`
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}

	if ret.Meta.Code != 200 {
		logger.ERROR("post %s error: %s with %s", u, ret.Meta.ErrorDetail, string(b))
		return nil, internalError
	}

	return ret.Response.Identities, nil
}

func (p *Platform) GetIcs(token string) (string, error) {
	url := fmt.Sprintf("%s/v2/ics/crosses?token=%s", p.config.SiteApi, token)
	reader, err := HttpResponse(HttpClient.Get(url))
	if err != nil {
		logger.ERROR("get %s error: %s", url, err)
		return "", internalError
	}
	defer reader.Close()
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.ERROR("get %s error: %s", url, err)
		return "", internalError
	}
	return string(b), nil
}

func (p *Platform) GetIdentityById(id int64) (model.Identity, error) {
	u := fmt.Sprintf("%s/v2/identities/%d", p.config.SiteApi, id)
	reader, err := HttpResponse(Http("GET", u, "applicatioin/json", nil))
	if err != nil {
		logger.ERROR("get %s error: %s", u, err)
		return model.Identity{}, err
	}

	defer reader.Close()
	var ret struct {
		Meta struct {
			Code        int    `json:"code"`
			ErrorDetail string `json:"errorDetail"`
		} `json:"meta"`
		Response struct {
			Identity model.Identity `json:"identity"`
		} `json:"response"`
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&ret)
	if err != nil || ret.Meta.Code != 200 {
		logger.ERROR("decode %s error: %s(%d %s)", u, err, ret.Meta.Code, ret.Meta.ErrorDetail)
		return model.Identity{}, err
	}
	return ret.Response.Identity, nil
}

func (p *Platform) GetRecipientsById(id string) ([]model.Recipient, error) {
	query := make(url.Values)
	query.Set("identity_id", id)
	u := fmt.Sprintf("%s/v3/bus/recipients?%s", p.config.SiteApi, query.Encode())
	reader, err := HttpResponse(Http("GET", u, "", nil))
	if err != nil {
		logger.ERROR("get %s error: %s", u, err)
		return nil, err
	}
	defer reader.Close()
	var ret struct {
		Data []model.Recipient
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&ret)
	if err != nil {
		logger.ERROR("decode %s error: %s", u, err)
		return nil, err
	}
	return ret.Data, nil
}

func (p *Platform) GetCrossByInvitationToken(token string) (model.Cross, error) {
	query := make(url.Values)
	query.Set("invitation_token", token)
	u := fmt.Sprintf("%s/v2/crosses/getcrossbyinvitationtoken", p.config.SiteApi)
	reader, err := HttpResponse(Http("POST", u, "application/x-www-form-urlencoded", []byte(query.Encode())))
	if err != nil {
		logger.ERROR("post %s error: %s with %s", u, err, query.Encode())
		return model.Cross{}, err
	}
	defer reader.Close()
	var ret struct {
		Meta struct {
			Code        int    `json:"code"`
			ErrorDetail string `json:"errorDetail"`
		} `json:"meta"`
		Response struct {
			Cross model.Cross `json:"cross"`
		} `json:"response"`
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&ret)
	if err != nil {
		logger.ERROR("decode %s error: %s with %s", u, err, query.Encode())
		return model.Cross{}, err
	}
	if ret.Meta.Code != 200 {
		return model.Cross{}, fmt.Errorf("%s", ret.Meta.ErrorDetail)
	}
	return ret.Response.Cross, nil
}

func (p *Platform) GetUserByIdentity(identity model.Identity) (model.User, string, error) {
	var resp struct {
		Data struct {
			Authorization struct {
				Token  string `json:"token"`
				UserId int64  `json:"user_id"`
			} `json:"authorization"`
			User model.User `json:"user"`
		} `json:"data"`
	}
	u := fmt.Sprintf("%s/v3/bus/users", p.config.SiteApi)
	_, err := RestHttp("POST", u, "application/json", identity, &resp)
	if err != nil {
		logger.ERROR("post %s error: %s with %+v", u, err, identity)
		return resp.Data.User, "", err
	}
	return resp.Data.User, resp.Data.Authorization.Token, nil
}

func (p *Platform) SetPassword(userId int64, password string) error {
	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", userId))
	query.Set("password", password)
	u := fmt.Sprintf("%s/v3/bus/setpassword", p.config.SiteApi)
	reader, err := HttpResponse(Http("POST", u, "application/x-www-form-urlencoded", []byte(query.Encode())))
	if err != nil {
		logger.ERROR("post %s error: %s with %s", u, err, query.Encode())
		return err
	}
	defer reader.Close()
	return nil
}

func (p *Platform) GetRouteXUrl(crossId uint64) (string, error) {
	query := make(url.Values)
	query.Set("cross_id", fmt.Sprintf("%d", crossId))
	u := fmt.Sprintf("%s/v3/bus/getroutexurl?%s", p.config.SiteApi, query.Encode())
	reader, err := HttpResponse(Http("GET", u, "", nil))
	if err != nil {
		logger.ERROR("get %s error: %s", u, err)
		return "", err
	}
	defer reader.Close()
	var ret struct {
		Data string `json:"data"`
	}
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&ret)
	if err != nil {
		logger.ERROR("decode %s error: %s", u, err)
		return "", err
	}
	return ret.Data, nil
}

func (p *Platform) GetWeatherIcon(lat, lng float64, date string) string {
	iconMap := map[string]string{
		"01d": "sun@2x.png",
		"02d": "sun@2x.png",
		"03d": "cloud_sun@2x.png",
		"04d": "cloud@2x.png",
		"09d": "drizzle_sun@2x.png",
		"10d": "rain@2x.png",
		"11d": "lightning@2x.png",
		"13d": "snow@2x.png",
		"50d": "haze@2x.png",
		"01n": "moon@2x.png",
		"02n": "moon@2x.png",
		"03n": "cloud_moon@2x.png",
		"04n": "cloud@2x.png",
		"09n": "drizzle_moon@2x.png",
		"10n": "rain_moon@2x.png",
		"11n": "lightning_moon@2x.png",
		"13n": "snow_moon@2x.png",
		"50n": "haze_moon@2x.png",
	}
	type Resp struct {
		List []struct {
			DtTxt   string `json:"dt_txt"`
			Weather []struct {
				Icon string `json:"icon"`
			} `json:"weather"`
		} `json:"list"`
	}
	u := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%.7f&lon=%.7f", lat, lng)
	var resp Resp
	_, err := RestHttp("GET", u, "application/json", nil, &resp)
	if err != nil {
		logger.ERROR("get weather %s failed: %s", u, err)
		return ""
	}
	icon := ""
	for _, list := range resp.List {
		if list.DtTxt <= date && len(list.Weather) > 0 {
			icon = list.Weather[0].Icon
		}
		if list.DtTxt > date {
			break
		}
	}
	ret, ok := iconMap[icon]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s/static/img/climacons/%s", p.config.SiteUrl, ret)
}

func (p *Platform) GetPlace(lat, lng float64, language string, radius int, query url.Values) ([]model.Place, error) {
	if query == nil {
		query = make(url.Values)
	}
	query.Set("key", p.config.Google.Key)
	query.Set("sensor", "false")
	query.Set("location", fmt.Sprintf("%07f,%07f", lat, lng))
	query.Set("radius", fmt.Sprintf("%d", radius))
	query.Set("language", language)
	u := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?%s", query.Encode())
	type Resp struct {
		Results []struct {
			Id       string
			Geometry struct {
				Location struct {
					Lat float64
					Lng float64
				}
			}
			Name     string
			Vicinity string
		}
		Status string
	}
	var resp Resp
	_, err := RestHttp("GET", u, "", nil, &resp)
	if err != nil {
		logger.ERROR("get place from google api %s failed: %s", u, err)
		return nil, err
	}
	if resp.Status != "OK" {
		logger.ERROR("get place from google api %s failed: %s", u, resp.Status)
		return nil, fmt.Errorf(resp.Status)
	}
	if len(resp.Results) == 0 {
		logger.ERROR("get place from google api %s failed: no place", u)
		return nil, nil
	}
	ret := make([]model.Place, len(resp.Results))
	for i, l := range resp.Results {
		ret[i] = model.Place{
			Title:       l.Name,
			Description: l.Vicinity,
			Lng:         fmt.Sprintf("%.7f", l.Geometry.Location.Lng),
			Lat:         fmt.Sprintf("%.7f", l.Geometry.Location.Lat),
			Provider:    "google",
			ExternalID:  l.Id,
		}
	}
	return ret, nil
}
