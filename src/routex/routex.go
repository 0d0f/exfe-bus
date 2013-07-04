package routex

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-broadcast"
	"github.com/googollee/go-rest"
	"io/ioutil"
	"logger"
	"model"
	"net/http"
	"time"
)

type RouteMap struct {
	rest.Service `prefix:"/v3/routex" mime:"application/json"`

	UpdateLocation rest.Processor `path:"/cross/:cross_id/location" method:"POST"`
	GetLocation    rest.Processor `path:"/cross/:cross_id/location" method:"GET"`
	UpdateRoute    rest.Processor `path:"cross/:cross_id/route" method:"POST"`
	GetRoute       rest.Processor `path:/cross/:cross_id/route" method:"GET"`
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
	token, ok := m.auth()
	if !ok {
		return
	}
	identity, err := m.platform.GetIdentityById(token.IdentityId)
	if err != nil {
		logger.ERROR("can't find identity %d: %s", token.IdentityId, err)
		m.Error(http.StatusInternalServerError, err)
		return
	}
	id := identity.Id()
	if err := m.locationRepo.Save(id, token.CrossId, location); err != nil {
		logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.CrossId, err, location)
		m.Error(http.StatusInternalServerError, err)
		return
	}
	broadcast, ok := m.broadcasts[token.CrossId]
	if !ok {
		return
	}
	locations, err := m.locationRepo.Load(id, token.CrossId)
	if err != nil {
		logger.ERROR("can't get locations %s of cross %d: %s", id, token.CrossId, err)
		return
	}
	broadcast.Send(map[string]interface{}{
		"name": fmt.Sprintf("/cross/%d/location", token.CrossId),
		"data": map[string]interface{}{
			id: locations,
		},
	})
}

func (m RouteMap) HandleGetLocation() map[string][]Location {
	token, ok := m.auth()
	if !ok {
		return nil
	}
	cross, err := m.platform.FindCross(int64(token.CrossId), nil)
	if err != nil {
		logger.ERROR("can't find cross %d: %s", token.CrossId, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	ret := make(map[string][]Location)
	for _, invitation := range cross.Exfee.Invitations {
		id := invitation.Identity.Id()
		locations, err := m.locationRepo.Load(id, token.CrossId)
		if err != nil {
			logger.ERROR("can't get locations %s of cross %d: %s", id, token.CrossId, err)
			continue
		}
		ret[id] = locations
	}
	return ret
}

func (m RouteMap) HandleUpdateRoute() {
	token, ok := m.auth()
	if !ok {
		return
	}
	b, err := ioutil.ReadAll(m.Request().Body)
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return
	}
	content := string(b)
	if err := m.routeRepo.Save(token.CrossId, content); err != nil {
		logger.ERROR("save route for cross %d failed: %s", token.CrossId, err)
		m.Error(http.StatusInternalServerError, err)
		return
	}
	broadcast, ok := m.broadcasts[token.CrossId]
	if !ok {
		return
	}
	broadcast.Send(map[string]interface{}{
		"name": fmt.Sprintf("/cross/%d/route", token.CrossId),
		"data": content,
	})
}

func (m RouteMap) HandleGetRoute() string {
	token, ok := m.auth()
	if !ok {
		return ""
	}
	content, err := m.routeRepo.Load(token.CrossId)
	if err != nil {
		logger.ERROR("can't get content of cross %d: %s", token.CrossId, err)
		m.Error(http.StatusInternalServerError, err)
		return ""
	}
	return content
}

func (m RouteMap) HandleNotification(stream *rest.Stream) {
	if m.Request().URL.Query().Get("_method") != "WATCH" {
		m.Error(http.StatusBadRequest, m.DetailError(-1, "method not watch"))
		return
	}
	token, ok := m.auth()
	if !ok {
		return
	}
	b, ok := m.broadcasts[token.CrossId]
	if !ok {
		b = broadcast.NewBroadcast()
		m.broadcasts[token.CrossId] = b
	}
	c := make(chan interface{})
	b.Register(c)
	for {
		d := <-c
		stream.SetWriteDeadline(time.Now().Add(broker.NetworkTimeout))
		err := stream.Write(d)
		if err != nil {
			return
		}
	}
}

func (m *RouteMap) auth() (CrossToken, bool) {
	authData := m.Request().Header.Get("Exfe-Auth-Data")
	crossIdStr := m.Vars()["cross_id"]
	var token CrossToken
	if err := json.Unmarshal([]byte(authData), &token); err != nil || crossIdStr != fmt.Sprintf("%d", token.CrossId) {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "auth failed"))
		return token, false
	}
	return token, true
}
