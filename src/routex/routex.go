package routex

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-broadcast"
	"github.com/googollee/go-rest"
	"logger"
	"math"
	"model"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type RouteMap struct {
	rest.Service `prefix:"/v3/routex" mime:"application/json"`

	UpdateLocation rest.Processor `path:"/cross/:cross_id/location" method:"POST"`
	GetLocation    rest.Processor `path:"/cross/:cross_id/location" method:"GET"`
	UpdateRoute    rest.Processor `path:"/cross/:cross_id/route" method:"POST"`
	GetRoute       rest.Processor `path:"/cross/:cross_id/route" method:"GET"`
	Notification   rest.Streaming `path:"/cross/:cross_id" method:"POST"`

	locationRepo LocationRepo
	routeRepo    RouteRepo
	platform     *broker.Platform
	config       *model.Config
	broadcasts   map[uint64]*broadcast.Broadcast
}

func New(locationRepo LocationRepo, routeRepo RouteRepo, platform *broker.Platform, config *model.Config) *RouteMap {
	return &RouteMap{
		locationRepo: locationRepo,
		routeRepo:    routeRepo,
		platform:     platform,
		config:       config,
		broadcasts:   make(map[uint64]*broadcast.Broadcast),
	}
}

func (m RouteMap) HandleUpdateLocation(location Location) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok {
		return
	}
	id := token.Identity.Id()

	lat, lng, _, err := location.GetGeo()
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return
	}

	location.Timestamp = time.Now().Unix()
	locations, err := m.locationRepo.Load(id, token.Cross.ID)
	if err != nil {
		logger.ERROR("can't get locations %s of cross %d: %s", id, token.Cross.ID, err)
	}
	if len(locations) == 0 {
		if err := m.locationRepo.Save(id, token.Cross.ID, location); err != nil {
			logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.Cross.ID, err, location)
			m.Error(http.StatusInternalServerError, err)
			return
		}
	} else {
		last := locations[0]
		lastLat, lastLng, _, err := last.GetGeo()
		if err != nil {
			if err := m.locationRepo.Save(id, token.Cross.ID, location); err != nil {
				logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.Cross.ID, err, location)
				m.Error(http.StatusInternalServerError, err)
				return
			}
		} else {
			a := math.Cos(lastLat) * math.Cos(lat) * math.Cos(lastLng-lng)
			b := math.Sin(lastLat) * math.Sin(lat)
			alpha := math.Acos(a + b)
			distance := alpha * 6371000
			if distance > 30 {
				if err := m.locationRepo.Save(id, token.Cross.ID, location); err != nil {
					logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.Cross.ID, err, location)
					m.Error(http.StatusInternalServerError, err)
					return
				}
			}
		}
	}
	locations = append([]Location{location}, locations...)

	broadcast, ok := m.broadcasts[token.Cross.ID]
	if !ok {
		return
	}
	if locations == nil {
		return
	}
	broadcast.Send(map[string]interface{}{
		"name": m.UpdateLocation.Path("cross_id", fmt.Sprintf("%d", token.Cross.ID)),
		"data": map[string]interface{}{
			id: locations,
		},
	})
}

func (m RouteMap) HandleGetLocation() map[string][]Location {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok {
		return nil
	}
	ret := make(map[string][]Location)
	for _, invitation := range token.Cross.Exfee.Invitations {
		id := invitation.Identity.Id()
		locations, err := m.locationRepo.Load(id, token.Cross.ID)
		if err != nil {
			logger.ERROR("can't get locations %s of cross %d: %s", id, token.Cross.ID, err)
			continue
		}
		if locations == nil {
			continue
		}
		ret[id] = locations
	}
	return ret
}

func (m RouteMap) HandleUpdateRoute(data []map[string]interface{}) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok {
		return
	}
	if err := m.routeRepo.Save(token.Cross.ID, data); err != nil {
		logger.ERROR("save route for cross %d failed: %s", token.Cross.ID, err)
		m.Error(http.StatusInternalServerError, err)
		return
	}
	broadcast, ok := m.broadcasts[token.Cross.ID]
	if !ok {
		return
	}
	broadcast.Send(map[string]interface{}{
		"name": m.UpdateRoute.Path("cross_id", fmt.Sprintf("%d", token.Cross.ID)),
		"data": data,
	})
}

func (m RouteMap) HandleGetRoute() []map[string]interface{} {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok {
		return nil
	}
	data, err := m.routeRepo.Load(token.Cross.ID)
	if err != nil {
		logger.ERROR("can't get route of cross %d: %s", token.Cross.ID, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	if data == nil {
		return make([]map[string]interface{}, 0)
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
		return
	}
	b, ok := m.broadcasts[token.Cross.ID]
	if !ok {
		b = broadcast.NewBroadcast(-1)
		m.broadcasts[token.Cross.ID] = b
	}
	c := make(chan interface{})
	b.Register(c)
	defer func() {
		b.Unregister(c)
		close(c)
	}()

	ret := make(map[string][]Location)
	for _, invitation := range token.Cross.Exfee.Invitations {
		id := invitation.Identity.Id()
		locations, err := m.locationRepo.Load(id, token.Cross.ID)
		if err != nil {
			logger.ERROR("can't get locations %s of cross %d: %s", id, token.Cross.ID, err)
			continue
		}
		if locations == nil {
			continue
		}
		ret[id] = locations
	}
	err := stream.Write(map[string]interface{}{
		"name": m.UpdateLocation.Path("cross_id", fmt.Sprintf("%d", token.Cross.ID)),
		"data": ret,
	})
	if err != nil {
		return
	}

	data, err := m.routeRepo.Load(token.Cross.ID)
	if err == nil && data != nil {
		err := stream.Write(map[string]interface{}{
			"name": m.UpdateRoute.Path("cross_id", fmt.Sprintf("%d", token.Cross.ID)),
			"data": data,
		})
		if err != nil {
			return
		}
	} else if err != nil {
		logger.ERROR("can't get route of cross %d: %s", token.Cross.ID, err)
	}

	for {
		select {
		case d := <-c:
			stream.SetWriteDeadline(time.Now().Add(broker.NetworkTimeout))
			err := stream.Write(d)
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

func (m *RouteMap) auth() (Token, bool) {
	var token Token

	crossIdStr := m.Vars()["cross_id"]
	crossId, err := strconv.ParseUint(crossIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusNotFound, m.DetailError(-1, "invalid cross id"))
		return token, false
	}

	authData := m.Request().Header.Get("Exfe-Auth-Data")
	if err := json.Unmarshal([]byte(authData), &token); err != nil {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return token, false
	}

	if token.TokenType != "cross_access_token" && token.TokenType != "user_token" {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return token, false
	}

	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", token.UserId))
	token.Cross, err = m.platform.FindCross(int64(crossId), query)
	if err != nil {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return token, false
	}

	for _, inv := range token.Cross.Exfee.Invitations {
		if inv.Identity.UserID == token.UserId {
			token.Identity = inv.Identity
			return token, true
		}
	}

	m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
	return token, false
}
