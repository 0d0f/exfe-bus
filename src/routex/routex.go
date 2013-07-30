package routex

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-broadcast"
	"github.com/googollee/go-rest"
	"logger"
	"model"
	"net/http"
	"net/url"
	"notifier"
	"strconv"
	"sync"
	"time"
)

type RouteMap struct {
	rest.Service `prefix:"/v3/routex" mime:"application/json"`

	SetUser          rest.Processor `path:"/user" method:"POST"`
	UpdateBreadcrums rest.Processor `path:"/breadcrumbs" method:"POST"`
	GetBreadcrums    rest.Processor `path:"/crosses/:cross_id/breadcrumbs" method:"GET"`
	UpdateGeomarks   rest.Processor `path:"/crosses/:cross_id/geomarks" method:"POST"`
	GetGeomarks      rest.Processor `path:"/crosses/:cross_id/geomarks" method:"GET"`
	Notification     rest.Streaming `path:"/crosses/:cross_id" method:"POST"`
	SendRequest      rest.Processor `path:"/crosses/:cross_id/request" method:"POST"`
	Options          rest.Processor `path:"/crosses/:cross_id" method:"OPTIONS"`

	breadcrumbsRepo BreadcrumbsRepo
	geomarksRepo    GeomarksRepo
	conversion      GeoConversionRepo
	platform        *broker.Platform
	config          *model.Config
	crossCast       map[uint64]*broadcast.Broadcast
	castLocker      sync.RWMutex
}

func New(breadcrumbsRepo BreadcrumbsRepo, geomarksRepo GeomarksRepo, conversion GeoConversionRepo, platform *broker.Platform, config *model.Config) *RouteMap {
	return &RouteMap{
		breadcrumbsRepo: breadcrumbsRepo,
		geomarksRepo:    geomarksRepo,
		conversion:      conversion,
		platform:        platform,
		config:          config,
		crossCast:       make(map[uint64]*broadcast.Broadcast),
	}
}

type UserCrossSetup struct {
	CrossId         uint64 `json:"cross_id,omitempty"`
	SaveBreadcrumbs bool   `json:"save_breadcrumbs,omitempty"`
}

func (m RouteMap) HandleSetUser(setup []UserCrossSetup) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	var token Token
	token, ok := m.auth()
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	var err error
	userId := token.UserId
	for _, s := range setup {
		if s.SaveBreadcrumbs {
			err = m.breadcrumbsRepo.EnableCross(userId, s.CrossId)
		} else {
			err = m.breadcrumbsRepo.DisableCross(userId, s.CrossId)
		}
	}
	if err != nil {
		logger.ERROR("setup user cross failed: %s", err)
		m.Error(http.StatusInternalServerError, err)
	}
}

func (m RouteMap) HandleUpdateBreadcrums(breadcrumbs []SimpleLocation) map[string]float64 {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	var token Token
	token, ok := m.auth()
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return nil
	}
	userId := token.UserId

	breadcrumb := breadcrumbs[0]
	if breadcrumb.Accuracy > 70 {
		m.Error(http.StatusBadRequest, fmt.Errorf("accuracy too large: %d", breadcrumb.Accuracy))
		return nil
	}

	if m.Request().URL.Query().Get("coordinate") == "mars" {
		breadcrumb.ToEarth(m.conversion)
	}

	breadcrumb.Timestamp = time.Now().Unix()
	route, err := m.breadcrumbsRepo.Load(userId)
	if err != nil {
		logger.ERROR("can't get route %s: %s", userId, err)
	}
	if len(route.Positions) == 0 {
		logger.INFO("routex", userId, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy)
		if err := m.breadcrumbsRepo.Save(userId, breadcrumb); err != nil {
			logger.ERROR("can't save repo %s: %s with %+v", userId, err, breadcrumb)
			m.Error(http.StatusInternalServerError, err)
			return nil
		}
	} else {
		distance := Distance(breadcrumb.Latitude, breadcrumb.Longitude, route.Positions[0].Latitude, route.Positions[0].Longitude)
		if distance > 30 {
			logger.INFO("routex", userId, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy, "distance", distance)
			if err := m.breadcrumbsRepo.Save(userId, breadcrumb); err != nil {
				logger.ERROR("can't save repo %s: %s with %+v", userId, err, breadcrumb)
				m.Error(http.StatusInternalServerError, err)
				return nil
			}
		} else {
			logger.INFO("routex", userId, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy, "distance", distance, "nosave")
		}
	}
	route.Positions = append([]SimpleLocation{breadcrumb}, route.Positions...)

	earth := breadcrumb
	mars := breadcrumb
	mars.ToMars(m.conversion)
	ret := map[string]float64{
		"earth_to_mars_latitude":  mars.Latitude - earth.Latitude,
		"earth_to_mars_longitude": mars.Longitude - earth.Longitude,
	}

	crosses, err := m.breadcrumbsRepo.Crosses(userId)
	if err != nil {
		logger.ERROR("can't get user %d cross: %s", err)
		return ret
	}
	for _, cross := range crosses {
		m.castLocker.RLock()
		b, ok := m.crossCast[cross]
		m.castLocker.RUnlock()
		if !ok {
			continue
		}
		b.Send(route)
	}

	return ret
}

func (m RouteMap) HandleGetBreadcrums() []Geomark {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	toMars := m.Request().URL.Query().Get("coordinate") == "mars"
	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return nil
	}
	var ret []Geomark
	for _, invitation := range token.Cross.Exfee.Invitations {
		userId := invitation.Identity.UserID
		route, err := m.breadcrumbsRepo.Load(userId)
		if err != nil {
			logger.ERROR("can't get user %d breadcrumbs: %s", userId, err)
			continue
		}
		if len(route.Positions) == 0 {
			continue
		}
		if toMars {
			route.ToMars(m.conversion)
		}
		ret = append(ret, route)
	}
	return ret
}

func (m RouteMap) HandleUpdateGeomarks(data []Geomark) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	m.castLocker.RLock()
	broadcast := m.crossCast[token.Cross.ID]
	m.castLocker.RUnlock()
	mars := m.Request().URL.Query().Get("coordinate") == "mars"
	if broadcast != nil || mars {
		for i, d := range data {
			if mars {
				d.ToEarth(m.conversion)
			}
			broadcast.Send(d)
			data[i] = d
		}
	}
	if err := m.geomarksRepo.Save(token.Cross.ID, data); err != nil {
		logger.ERROR("save route for cross %d failed: %s", token.Cross.ID, err)
		m.Error(http.StatusInternalServerError, err)
		return
	}
}

func (m RouteMap) HandleGetGeomarks() []Geomark {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return nil
	}
	data, err := m.geomarksRepo.Load(token.Cross.ID)
	if err != nil {
		logger.ERROR("can't get route of cross %d: %s", token.Cross.ID, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	if data == nil {
		lng, err := strconv.ParseFloat(token.Cross.Place.Lng, 64)
		if err != nil {
			token.Cross.Place = nil
		}
		lat, err := strconv.ParseFloat(token.Cross.Place.Lat, 64)
		if err != nil {
			token.Cross.Place = nil
		}
		if token.Cross.Place != nil {
			createdAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.CreatedAt)
			if err != nil {
				createdAt = time.Now()
			}
			updatedAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.UpdatedAt)
			if err != nil {
				updatedAt = time.Now()
			}
			destinaion := Geomark{
				Id:          "destination",
				Type:        "location",
				CreatedAt:   createdAt.Unix(),
				CreatedBy:   token.Cross.By.Id(),
				UpdatedAt:   updatedAt.Unix(),
				UpdatedBy:   token.Cross.By.Id(),
				Tags:        []string{"destination"},
				Title:       token.Cross.Place.Title,
				Description: token.Cross.Place.Description,
				Longitude:   lng,
				Latitude:    lat,
			}
			data = []Geomark{destinaion}
		}
	}
	if m.Request().URL.Query().Get("coordinate") == "mars" {
		for i := range data {
			data[i].ToMars(m.conversion)
		}
	}
	return data
}

func (m RouteMap) HandleNotification(stream rest.Stream) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	if m.Request().URL.Query().Get("_method") != "WATCH" {
		m.Error(http.StatusBadRequest, m.DetailError(-1, "method not watch"))
		return
	}
	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	m.castLocker.Lock()
	b, ok := m.crossCast[token.Cross.ID]
	if !ok {
		b = broadcast.NewBroadcast(-1)
		m.crossCast[token.Cross.ID] = b
	}
	m.castLocker.Unlock()
	c := make(chan interface{})
	b.Register(c)
	defer func() {
		b.Unregister(c)
		close(c)

		if b.Len() == 0 {
			m.castLocker.Lock()
			delete(m.crossCast, token.Cross.ID)
			m.castLocker.Unlock()
		}
	}()

	toMars := m.Request().URL.Query().Get("coordinate") == "mars"

	for _, invitation := range token.Cross.Exfee.Invitations {
		userId := invitation.Identity.UserID
		route, err := m.breadcrumbsRepo.Load(userId)
		if err != nil {
			logger.ERROR("can't get user %d breadcrumbs: %s", userId, err)
			continue
		}
		if toMars {
			route.ToMars(m.conversion)
		}
		err = stream.Write(route)
		if err != nil {
			return
		}
	}

	marks, err := m.geomarksRepo.Load(token.Cross.ID)
	if err == nil {
		if marks == nil {
			lng, err := strconv.ParseFloat(token.Cross.Place.Lng, 64)
			if err != nil {
				token.Cross.Place = nil
			}
			lat, err := strconv.ParseFloat(token.Cross.Place.Lat, 64)
			if err != nil {
				token.Cross.Place = nil
			}
			if token.Cross.Place != nil {
				createdAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.CreatedAt)
				if err != nil {
					createdAt = time.Now()
				}
				updatedAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.UpdatedAt)
				if err != nil {
					updatedAt = time.Now()
				}
				destinaion := Geomark{
					Id:          "destination",
					Type:        "location",
					CreatedAt:   createdAt.Unix(),
					CreatedBy:   token.Cross.By.Id(),
					UpdatedAt:   updatedAt.Unix(),
					UpdatedBy:   token.Cross.By.Id(),
					Tags:        []string{"destination"},
					Title:       token.Cross.Place.Title,
					Description: token.Cross.Place.Description,
					Longitude:   lng,
					Latitude:    lat,
				}
				marks = append(marks, destinaion)
			}
		}
		for _, d := range marks {
			d.ToMars(m.conversion)
			err := stream.Write(d)
			if err != nil {
				return
			}
		}
	} else if err != nil {
		logger.ERROR("can't get route of cross %d: %s", token.Cross.ID, err)
	}

	for {
		select {
		case d := <-c:
			mark, ok := d.(Geomark)
			if !ok {
				continue
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

func (m RouteMap) HandleSendRequest(id string) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

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

	var fromIdentity model.Identity
	for _, inv := range token.Cross.Exfee.Invitations {
		if inv.Identity.UserID == token.UserId {
			fromIdentity = inv.Identity
		}
	}
	for _, recipient := range recipients {
		switch recipient.Provider {
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
	return
}

func (m *RouteMap) auth() (Token, bool) {
	var token Token

	authData := m.Request().Header.Get("Exfe-Auth-Data")
	// if authData == "" {
	// 	// token: 345ac9296016c858a752a7e5fea35b7682fa69f922c4cefa30cfc22741da3109
	// 	authData = `{"token_type":"cross_access_token","cross_id":100758,"identity_id":907,"user_id":652,"created_time":1374636534,"updated_time":1374636534}`
	// }
	// logger.DEBUG("auth data: %s", authData)

	if err := json.Unmarshal([]byte(authData), &token); err != nil {
		return token, false
	}

	crossIdStr := m.Vars()["cross_id"]
	if crossIdStr == "" {
		switch token.TokenType {
		case "user_token":
			return token, true
		case "cross_access_token":
			var err error
			token.Cross, err = m.platform.FindCross(int64(token.CrossId), nil)
			if err != nil {
				return token, false
			}
			identity, err := m.platform.GetIdentityById(token.IdentityId)
			if err != nil {
				return token, false
			}
			token.UserId = identity.UserID
			return token, true
		}
		return token, false
	}
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
