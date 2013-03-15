package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gobus"
	"io"
	"io/ioutil"
	"model"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	ProcessTimeout = 60 * time.Second
	NetworkTimeout = 30 * time.Second
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
	}, nil
}

func (p *Platform) Send(to model.Recipient, text string) (string, error) {
	arg := model.ThirdpartSend{
		To:   to,
		Text: text,
	}

	var ids string
	err := p.dispatcher.DoWithTicket(to.Provider, "bus://exfe_service/thirdpart/message", "POST", &arg, &ids)

	if err != nil {
		return "", err
	}
	return ids, nil
}

func (p *Platform) GetHotRecipient(userID int64) ([]model.Recipient, error) {
	return nil, nil

	resp, err := http.Get(fmt.Sprintf("%s/v2/Gobus/HotRecipient?user=%d", p.config.SiteApi, userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("find identity failed: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	var ret []model.Recipient
	err = decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (p *Platform) FindIdentity(identity model.Identity) (model.Identity, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(identity)
	if err != nil {
		return identity, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/v2/Gobus/RevokeIdentity", p.config.SiteApi), "application/json", buf)
	if err != nil {
		return identity, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return identity, fmt.Errorf("find identity failed: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&identity)
	if err != nil {
		return identity, err
	}
	return identity, nil
}

func (p *Platform) FindCross(id uint64) (model.Cross, error) {
	var ret model.Cross
	resp, err := http.Get(fmt.Sprintf("%s/v2/Gobus/GetCrossById?id=%d", p.config.SiteApi, id))
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ret, fmt.Errorf("find cross failed: %s", resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (p *Platform) UploadPhoto(photoxID string, photos []model.Photo) error {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(photos)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/v2/Gobus/AddPhotos/%s", p.config.SiteApi, photoxID), "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("upload photo fail: %s", resp.Status)
	}
	return nil
}

func (p *Platform) BotCrossGather(cross model.Cross) (uint64, int, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(cross)
	if err != nil {
		return 0, 500, err
	}
	u := fmt.Sprintf("%s/v2/Gobus/Gather", p.config.SiteApi)
	p.config.Log.Debug("bot gather to: %s, cross: %s", u, buf.String())
	body, code, err := parseResp(client.Post(u, "application/json", buf))
	if err != nil {
		return 0, code, fmt.Errorf("error(%s) when send message(%s)", err, buf.String())
	}
	defer body.Close()
	decoder := json.NewDecoder(body)
	err = decoder.Decode(&cross)
	if err != nil {
		p.config.Log.Crit("can't parse gather return: %s", err)
		return 0, 500, err
	}

	return cross.ID, 200, nil
}

func (p *Platform) BotCrossUpdate(to, id string, cross model.Cross, by model.Identity) (int, error) {
	arg := make(map[string]interface{})
	arg[to] = id
	arg["cross"] = cross
	arg["by_identity"] = by
	fmt.Println(arg)

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(arg)
	if err != nil {
		return 500, err
	}

	u := fmt.Sprintf("%s/v2/Gobus/XUpdate", p.config.SiteApi)
	p.config.Log.Debug("bot invite to: %s, arg: %s", u, buf.String())
	body, code, err := parseResp(client.Post(u, "application/json", buf))
	if err != nil {
		return code, fmt.Errorf("error(%s) when send message(%s)", err, buf.String())
	}
	defer body.Close()

	return 200, nil
}

func (p *Platform) BotPostConversation(from, post, to, id string) (int, error) {
	u := fmt.Sprintf("%s/v2/Gobus/PostConversation", p.config.SiteApi)
	params := make(url.Values)
	params.Add(to, id)
	params.Add("content", post)
	params.Add("external_id", from)
	params.Add("provider", "email")
	p.config.Log.Debug("bot post to: %s, post content: %s\n", u, params.Encode())

	body, code, err := parseResp(client.PostForm(u, params))
	if err != nil {
		return code, fmt.Errorf("error(%s) when send message(%s)", err, params.Encode())
	}
	defer body.Close()

	return 200, nil
}

func parseResp(resp *http.Response, err error) (io.ReadCloser, int, error) {
	if err != nil {
		return nil, 500, err
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, 500, err
		}
		return nil, resp.StatusCode, fmt.Errorf("%s: %s", resp.Status, body)
	}
	return resp.Body, 200, nil
}
