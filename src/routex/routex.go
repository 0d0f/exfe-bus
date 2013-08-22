package routex

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-broadcast"
	"github.com/googollee/go-rest"
	"logger"
	"math/rand"
	"model"
	"net/http"
	"net/url"
	"notifier"
	"os"
	"strconv"
	"sync"
	"time"
)

type RouteMap struct {
	rest.Service `prefix:"/v3/routex" mime:"application/json"`

	SearchRoutex rest.Processor `path:"/_inner/search/crosses" method:"POST"`
	GetRoutex    rest.Processor `path:"/_inner/users/:user_id/crosses/:cross_id" method:"GET"`
	SetUserInner rest.Processor `path:"/_inner/users/:user_id/crosses/:cross_id" method:"POST"`
	SetUser      rest.Processor `path:"/users/crosses/:cross_id" method:"POST"`

	UpdateBreadcrums       rest.Processor `path:"/breadcrumbs" method:"POST"`
	UpdateBreadcrumsInner  rest.Processor `path:"/_inner/breadcrumbs/users/:user_id" method:"POST"`
	GetBreadcrums          rest.Processor `path:"/breadcrumbs/crosses/:cross_id" method:"GET"`
	GetUserBreadcrums      rest.Processor `path:"/breadcrumbs/crosses/:cross_id/users/:user_id" method:"GET"`
	GetUserBreadcrumsInner rest.Processor `path:"/_inner/breadcrumbs/users/:user_id" method:"GET"`

	SearchGeomarks rest.Processor `path:"/_inner/geomarks/crosses/:cross_id" method:"GET"`
	GetGeomarks    rest.Processor `path:"/geomarks/crosses/:cross_id" method:"GET"`
	SetGeomark     rest.Processor `path:"/geomarks/crosses/:cross_id/:mark_type/:mark_id" method:"PUT"`
	DeleteGeomark  rest.Processor `path:"/geomarks/crosses/:cross_id/:mark_type/:mark_id" method:"DELETE"`

	Stream  rest.Streaming `path:"/crosses/:cross_id" method:"WATCH"`
	Options rest.Processor `path:"/crosses/:cross_id" method:"OPTIONS"`

	SendNotification rest.Processor `path:"/notification/crosses/:cross_id/:identity_id" method:"POST"`

	rand            *rand.Rand
	routexRepo      RoutexRepo
	breadcrumbCache BreadcrumbCache
	breadcrumbsRepo BreadcrumbsRepo
	geomarksRepo    GeomarksRepo
	conversion      GeoConversionRepo
	platform        *broker.Platform
	config          *model.Config
	tutorialDatas   map[int64][]TutorialData
	crossCast       map[int64]*broadcast.Broadcast
	castLocker      sync.RWMutex
}

func New(routexRepo RoutexRepo, breadcrumbCache BreadcrumbCache, breadcrumbsRepo BreadcrumbsRepo, geomarksRepo GeomarksRepo, conversion GeoConversionRepo, platform *broker.Platform, config *model.Config) (*RouteMap, error) {
	tutorialDatas := make(map[int64][]TutorialData)
	for _, userId := range config.TutorialBotUserIds {
		file := config.Routex.TutorialDataFile[fmt.Sprintf("%d", userId)]
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("can't find tutorial file %s for tutorial bot %d", file, userId)
		}
		var datas []TutorialData
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&datas)
		if err != nil {
			return nil, fmt.Errorf("invalid tutorial data %s for tutorial bot %d: %s", file, userId, err)
		}
		tutorialDatas[userId] = datas
	}
	ret := &RouteMap{
		rand:            rand.New(rand.NewSource(time.Now().Unix())),
		routexRepo:      routexRepo,
		breadcrumbCache: breadcrumbCache,
		breadcrumbsRepo: breadcrumbsRepo,
		geomarksRepo:    geomarksRepo,
		conversion:      conversion,
		platform:        platform,
		tutorialDatas:   tutorialDatas,
		config:          config,
		crossCast:       make(map[int64]*broadcast.Broadcast),
	}
	return ret, nil
}

func (m RouteMap) setTutorial(lat, lng float64, userId, crossId int64, locale string) (Geomark, error) {
	var ret Geomark
	query := make(url.Values)
	query.Set("keyword", "attractions")
	places, err := m.platform.GetPlace(lat, lng, locale, 10000, query)
	if err != nil {
		return ret, err
	}
	if len(places) == 0 {
		places, err = m.platform.GetPlace(lat, lng, locale, 50000, nil)
		if err != nil {
			return ret, err
		}
	}
	if len(places) == 0 {
		return ret, fmt.Errorf("can't find attraction place near %.7f,%.7f", lat, lng)
	}
	place := places[0]
	if lng, err = strconv.ParseFloat(place.Lng, 64); err != nil {
		return ret, err
	}
	if lat, err = strconv.ParseFloat(place.Lat, 64); err != nil {
		return ret, err
	}
	now := time.Now().Unix()
	ret = Geomark{
		Id:          fmt.Sprintf("%08d@location", m.rand.Intn(1e8)),
		Type:        "location",
		CreatedAt:   now,
		CreatedBy:   fmt.Sprintf("%d@exfe", userId),
		UpdatedAt:   now,
		UpdatedBy:   fmt.Sprintf("%d@exfe", userId),
		Tags:        []string{"destination"},
		Icon:        "http://panda.0d0f.com/static/img/map_pin_blue@2x.png",
		Title:       place.Title,
		Description: place.Description,
		Longitude:   lng,
		Latitude:    lat,
	}
	if err := m.geomarksRepo.Set(crossId, ret); err != nil {
		return ret, err
	}
	return ret, nil
}

type UserCrossSetup struct {
	SaveBreadcrumbs bool `json:"save_breadcrumbs,omitempty"`
	AfterInSeconds  int  `json:"after_in_seconds,omitempty"`
}

func (m RouteMap) HandleSetUser(setup UserCrossSetup) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	var token Token
	token, ok := m.auth(true)
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	m.Vars()["user_id"] = fmt.Sprintf("%d", token.UserId)
	m.HandleSetUserInner(setup)
}

func (m RouteMap) HandleSetUserInner(setup UserCrossSetup) {
	userIdStr := m.Vars()["user_id"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return
	}
	crossIdStr := m.Vars()["cross_id"]
	crossId, err := strconv.ParseInt(crossIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return
	}
	go func() {
		if setup.SaveBreadcrumbs {
			if setup.AfterInSeconds == 0 {
				setup.AfterInSeconds = 7200
			}
			if err := m.routexRepo.EnableCross(userId, crossId, setup.AfterInSeconds); err != nil {
				logger.ERROR("set user %d enable cross %d routex repo failed: %s", userId, crossId, err)
			}
			if err := m.breadcrumbsRepo.EnableCross(userId, crossId, setup.AfterInSeconds); err != nil {
				logger.ERROR("set user %d enable cross %d breadcrumbs repo failed: %s", userId, crossId, err)
			}
		} else {
			if err := m.routexRepo.DisableCross(userId, crossId); err != nil {
				logger.ERROR("set user %d disable cross %d routex repo failed: %s", userId, crossId, err)
			}
			if err := m.breadcrumbsRepo.DisableCross(userId, crossId); err != nil {
				logger.ERROR("set user %d disable cross %d breadcrumbs repo failed: %s", userId, crossId, err)
			}
		}
	}()

	if setup.SaveBreadcrumbs {
		if setup.AfterInSeconds == 0 {
			setup.AfterInSeconds = 7200
		}
		if err := m.breadcrumbCache.EnableCross(userId, crossId, setup.AfterInSeconds); err != nil {
			logger.ERROR("set user %d enable cross %d breadcrumb cache failed: %s", userId, crossId, err)
		}
	} else {
		if err := m.breadcrumbCache.DisableCross(userId, crossId); err != nil {
			logger.ERROR("set user %d disable cross %d breadcrumb cache failed: %s", userId, crossId, err)
		}
	}
}

func (m RouteMap) HandleSearchRoutex(crossIds []int64) []Routex {
	ret, err := m.routexRepo.Search(crossIds)
	if err != nil {
		logger.ERROR("search for route failed: %s with %+v", err, crossIds)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	return ret
}

func (m RouteMap) HandleGetRoutex() *bool {
	userIdStr, crossIdStr := m.Vars()["user_id"], m.Vars()["cross_id"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, fmt.Errorf("invalid user id %s", userIdStr))
		return nil
	}
	crossId, err := strconv.ParseInt(crossIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, fmt.Errorf("invalid user id %s", crossIdStr))
		return nil
	}
	route, err := m.routexRepo.Get(userId, crossId)
	if err != nil {
		logger.ERROR("get user %d cross %d routex failed: %s", userId, crossId, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	if route == nil {
		return nil
	}
	return &route.Enable
}

func (m RouteMap) HandleStream(stream rest.Stream) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth(true)
	if !ok {
		logger.DEBUG("invalid token")
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	m.castLocker.Lock()
	b, ok := m.crossCast[int64(token.Cross.ID)]
	if !ok {
		b = broadcast.NewBroadcast(-1)
		m.crossCast[int64(token.Cross.ID)] = b
	}
	m.castLocker.Unlock()
	c := make(chan interface{})
	b.Register(c)
	defer func() {
		b.Unregister(c)
		close(c)

		if b.Len() == 0 {
			m.castLocker.Lock()
			delete(m.crossCast, int64(token.Cross.ID))
			m.castLocker.Unlock()
		}
	}()

	toMars := m.Request().URL.Query().Get("coordinate") == "mars"
	isTutorial := false
	if token.Cross.By.ID == m.config.Routex.TutorialCreator {
		isTutorial = true
	}
	hasCreated := false

	m.WriteHeader(http.StatusOK)
	quit := make(chan int)
	defer func() { close(quit) }()

	for _, invitation := range token.Cross.Exfee.Invitations {
		userId := invitation.Identity.UserID
		route := Geomark{
			Id:   fmt.Sprintf("%d@exfe", userId),
			Type: "route",
			Tags: []string{"breadcrumbs"},
		}
		if route.Positions = m.getTutorialData(time.Now().UTC(), userId, 1); route.Positions != nil {
			go func() {
				for {
					select {
					case <-quit:
						return
					case <-time.After(time.Second * 10):
						route := Geomark{
							Id:        fmt.Sprintf("%d@exfe", userId),
							Type:      "route",
							Tags:      []string{"breadcrumbs"},
							Positions: m.getTutorialData(time.Now().UTC(), userId, 1),
						}
						c <- route
					}
				}
			}()
		} else {
			l, exist, err := m.breadcrumbCache.LoadCross(userId, int64(token.Cross.ID))
			if err != nil {
				logger.ERROR("can't get user %d breadcrumbs of cross %d: %s", userId, token.Cross.ID, err)
				continue
			}
			if !exist {
				continue
			}
			route.Positions = []SimpleLocation{l}
		}
		if toMars {
			route.ToMars(m.conversion)
		}
		err := stream.Write(route)
		if err != nil {
			return
		}
	}

	marks, err := m.getGeomarks(token.Cross, toMars)
	if err == nil {
		for _, d := range marks {
			if isTutorial && !hasCreated {
				hasCreated = true
			}
			err := stream.Write(d)
			if err != nil {
				return
			}
		}
	} else {
		logger.ERROR("can't get route of cross %d: %s", token.Cross.ID, err)
	}

	for {
		select {
		case d := <-c:
			mark, ok := d.(Geomark)
			if !ok {
				continue
			}
			if isTutorial && !hasCreated {
				if mark.Id == fmt.Sprintf("%d@exfe", token.UserId) {
					locale := ""
					for _, i := range token.Cross.Exfee.Invitations {
						if i.Identity.UserID == token.UserId {
							locale = i.Identity.Locale
							break
						}
					}
					tutorialMark, err := m.setTutorial(mark.Latitude, mark.Longitude, token.UserId, int64(token.Cross.ID), locale)
					if err != nil {
						logger.ERROR("create tutorial geomark for user %d in cross %d failed: %s", token.UserId, token.Cross.ID, err)
					} else {
						hasCreated = true
						if toMars {
							tutorialMark.ToMars(m.conversion)
						}
						err := stream.Write(tutorialMark)
						if err != nil {
							return
						}
					}
				}
			}
			if toMars {
				mark.ToMars(m.conversion)
			}
			stream.SetWriteDeadline(time.Now().Add(broker.NetworkTimeout))
			err := stream.Write(mark)
			if err != nil {
				return
			}
		case <-time.After(broker.NetworkTimeout):
			err := stream.Ping()
			if err != nil {
				return
			}
		}
	}
}

func (m RouteMap) HandleOptions() {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	m.WriteHeader(http.StatusNoContent)
}

func (m RouteMap) HandleSendNotification() {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth(true)
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	id := m.Vars()["identity_id"]
	identity, ok := model.FromIdentityId(id), false
	for _, inv := range token.Cross.Exfee.Invitations {
		if inv.Identity.Equal(identity) {
			ok = true
			break
		}
	}
	if !ok {
		m.Error(http.StatusForbidden, fmt.Errorf("%s is not attend cross %d", id, token.Cross.ID))
		return
	}

	recipients, err := m.platform.GetRecipientsById(id)
	if err != nil {
		m.Error(http.StatusInternalServerError, err)
		return
	}

	m.routexRepo.Update(token.UserId, int64(token.Cross.ID))
	var fromIdentity model.Identity
	for _, inv := range token.Cross.Exfee.Invitations {
		if inv.Identity.UserID == token.UserId {
			fromIdentity = inv.Identity
		}
	}
	for _, recipient := range recipients {
		switch recipient.Provider {
		case "wechat":
			fallthrough
		case "iOS":
			fallthrough
		case "Android":
			body, err := json.Marshal(notifier.RequestArg{
				To:      recipient,
				CrossId: token.Cross.ID,
				From:    fromIdentity,
			})
			if err != nil {
				logger.ERROR("can't marshal: %s with %+v", err, recipient)
				continue
			}
			url := fmt.Sprintf("http://%s:%d/v3/notifier/routex/request", m.config.ExfeService.Addr, m.config.ExfeService.Port)
			resp, err := broker.HttpResponse(broker.Http("POST", url, "applicatioin/json", body))
			if err != nil {
				logger.ERROR("call %s error: %s with %#v", url, err, string(body))
				continue
			}
			resp.Close()
			return
		}
	}
	m.Error(http.StatusNotAcceptable, fmt.Errorf("can't find provider avaliable"))
	return
}

func (m *RouteMap) auth(checkCross bool) (Token, bool) {
	var token Token

	authData := m.Request().Header.Get("Exfe-Auth-Data")
	if authData == "" {
		authData = `{"token_type":"user_token","user_id":475,"signin_time":1374046388,"last_authenticate":1374046388}`
	}

	if authData != "" {
		if err := json.Unmarshal([]byte(authData), &token); err != nil {
			return token, false
		}
	}

	if !checkCross {
		if token.TokenType == "user_token" {
			return token, true
		}
		return token, false
	}

	crossIdStr := m.Vars()["cross_id"]
	crossId, err := strconv.ParseUint(crossIdStr, 10, 64)
	if err != nil {
		return token, false
	}

	query := make(url.Values)
	switch token.TokenType {
	case "user_token":
		query.Set("user_id", fmt.Sprintf("%d", token.UserId))
	case "cross_access_token":
		if token.CrossId != crossId {
			return token, false
		}
	default:
		t := m.Request().URL.Query().Get("token")
		cross, err := m.platform.GetCrossByInvitationToken(t)
		if err != nil {
			return token, false
		}
		if cross.ID != crossId {
			return token, false
		}
		token.Cross, token.Readonly = cross, true
		return token, true
	}

	token.Cross, err = m.platform.FindCross(int64(crossId), query)
	if err != nil {
		return token, false
	}

	if token.TokenType == "user_token" {
		return token, true
	}

	for _, inv := range token.Cross.Exfee.Invitations {
		if inv.Identity.ID == token.IdentityId {
			token.UserId = inv.Identity.UserID
			return token, true
		}
	}
	return token, false
}
