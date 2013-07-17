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
	"notifier"
	"strconv"
	"time"
)

type RouteMap struct {
	rest.Service `prefix:"/v3/routex" mime:"application/json"`

	UpdateBreadcrums rest.Processor `path:"/crosses/:cross_id/breadcrumbs" method:"POST"`
	GetBreadcrums    rest.Processor `path:"/crosses/:cross_id/breadcrumbs" method:"GET"`
	UpdateGeomarks   rest.Processor `path:"/crosses/:cross_id/geomarks" method:"POST"`
	GetGeomarks      rest.Processor `path:"/crosses/:cross_id/geomarks" method:"GET"`
	Notification     rest.Streaming `path:"/crosses/:cross_id" method:"POST"`
	Options          rest.Processor `path:"/crosses/:cross_id" method:"OPTIONS"`
	SendRequest      rest.Processor `path:"/crosses/:cross_id/request" method:"POST"`

	breadcrumbsRepo BreadcrumbsRepo
	geomarksRepo    GeomarksRepo
	platform        *broker.Platform
	config          *model.Config
	broadcasts      map[uint64]*broadcast.Broadcast
}

func New(breadcrumbsRepo BreadcrumbsRepo, geomarksRepo GeomarksRepo, platform *broker.Platform, config *model.Config) *RouteMap {
	return &RouteMap{
		breadcrumbsRepo: breadcrumbsRepo,
		geomarksRepo:    geomarksRepo,
		platform:        platform,
		config:          config,
		broadcasts:      make(map[uint64]*broadcast.Broadcast),
	}
}

func (m RouteMap) HandleUpdateBreadcrums(breadcrumb Location) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok {
		return
	}
	id := token.Identity.Id()

	lat, lng, _, err := breadcrumb.GetGeo()
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return
	}

	breadcrumb.Timestamp = time.Now().Unix()
	breadcrumbs, err := m.breadcrumbsRepo.Load(id, token.Cross.ID)
	if err != nil {
		logger.ERROR("can't get breadcrumbs %s of cross %d: %s", id, token.Cross.ID, err)
	}
	if len(breadcrumbs) == 0 {
		if err := m.breadcrumbsRepo.Save(id, token.Cross.ID, breadcrumb); err != nil {
			logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.Cross.ID, err, breadcrumb)
			m.Error(http.StatusInternalServerError, err)
			return
		}
	} else {
		last := breadcrumbs[0]
		lastLat, lastLng, _, err := last.GetGeo()
		if err != nil {
			if err := m.breadcrumbsRepo.Save(id, token.Cross.ID, breadcrumb); err != nil {
				logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.Cross.ID, err, breadcrumb)
				m.Error(http.StatusInternalServerError, err)
				return
			}
		} else {
			a := math.Cos(lastLat) * math.Cos(lat) * math.Cos(lastLng-lng)
			b := math.Sin(lastLat) * math.Sin(lat)
			alpha := math.Acos(a + b)
			distance := alpha * 6371000
			logger.INFO("routex", "cross", token.Cross.ID, id, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy, "last", last.Longitude, last.Latitude, last.Accuracy, "alpha", alpha, "distance", distance)
			if distance > 30 {
				if err := m.breadcrumbsRepo.Save(id, token.Cross.ID, breadcrumb); err != nil {
					logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.Cross.ID, err, breadcrumb)
					m.Error(http.StatusInternalServerError, err)
					return
				}
			}
		}
	}
	breadcrumbs = append([]Location{breadcrumb}, breadcrumbs...)

	broadcast, ok := m.broadcasts[token.Cross.ID]
	if !ok {
		return
	}
	if breadcrumbs == nil {
		return
	}
	broadcast.Send(map[string]interface{}{
		"type": "/v3/crosses/routex/breadcrumbs",
		"data": map[string]interface{}{
			id: breadcrumbs,
		},
	})
}

func (m RouteMap) HandleGetBreadcrums() map[string][]Location {
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
		breadcrumbs, err := m.breadcrumbsRepo.Load(id, token.Cross.ID)
		if err != nil {
			logger.ERROR("can't get breadcrumbs %s of cross %d: %s", id, token.Cross.ID, err)
			continue
		}
		if breadcrumbs == nil {
			continue
		}
		ret[id] = breadcrumbs
	}
	return ret
}

func (m RouteMap) HandleUpdateGeomarks(data []map[string]interface{}) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok {
		return
	}
	if err := m.geomarksRepo.Save(token.Cross.ID, data); err != nil {
		logger.ERROR("save route for cross %d failed: %s", token.Cross.ID, err)
		m.Error(http.StatusInternalServerError, err)
		return
	}
	broadcast, ok := m.broadcasts[token.Cross.ID]
	if !ok {
		return
	}
	broadcast.Send(map[string]interface{}{
		"type": "/v3/crosses/routex/geomarks",
		"data": data,
	})
}

func (m RouteMap) HandleGetGeomarks() []map[string]interface{} {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok {
		return nil
	}
	data, err := m.geomarksRepo.Load(token.Cross.ID)
	if err != nil {
		logger.ERROR("can't get route of cross %d: %s", token.Cross.ID, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	if data == nil {
		ret := make([]map[string]interface{}, 0)
		if token.Cross.Place.Lng != "" && token.Cross.Place.Lat != "" {
			createdAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.CreatedAt)
			if err != nil {
				createdAt = time.Now()
			}
			updatedAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.UpdatedAt)
			if err != nil {
				updatedAt = time.Now()
			}
			destinaion := map[string]interface{}{
				"id":          "destination",
				"type":        "location",
				"created_at":  createdAt.Unix(),
				"created_by":  token.Cross.By.Id(),
				"updated_at":  updatedAt.Unix(),
				"updated_by":  token.Cross.By.Id(),
				"tags":        []string{"destination"},
				"title":       token.Cross.Place.Title,
				"description": token.Cross.Place.Description,
				"timestamp":   time.Now().Unix(),
				"longitude":   token.Cross.Place.Lng,
				"latitude":    token.Cross.Place.Lat,
			}
			ret = append(ret, destinaion)
		}
		return ret
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
		breadcrumbs, err := m.breadcrumbsRepo.Load(id, token.Cross.ID)
		if err != nil {
			logger.ERROR("can't get breadcrumbs %s of cross %d: %s", id, token.Cross.ID, err)
			continue
		}
		if breadcrumbs == nil {
			continue
		}
		ret[id] = breadcrumbs
	}
	err := stream.Write(map[string]interface{}{
		"type": "/v3/crosses/routex/breadcrumbs",
		"data": ret,
	})
	if err != nil {
		return
	}

	data, err := m.geomarksRepo.Load(token.Cross.ID)
	if err == nil {
		if data == nil {
			data = make([]map[string]interface{}, 0)
			if token.Cross.Place.Lng != "" && token.Cross.Place.Lat != "" {
				createdAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.CreatedAt)
				if err != nil {
					createdAt = time.Now()
				}
				updatedAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.UpdatedAt)
				if err != nil {
					updatedAt = time.Now()
				}
				destinaion := map[string]interface{}{
					"id":          "destination",
					"type":        "location",
					"created_at":  createdAt.Unix(),
					"created_by":  token.Cross.By.Id(),
					"updated_at":  updatedAt.Unix(),
					"updated_by":  token.Cross.By.Id(),
					"tags":        []string{"destination"},
					"title":       token.Cross.Place.Title,
					"description": token.Cross.Place.Description,
					"timestamp":   time.Now().Unix(),
					"longitude":   token.Cross.Place.Lng,
					"latitude":    token.Cross.Place.Lat,
				}
				data = append(data, destinaion)
			}
		}
		err := stream.Write(map[string]interface{}{
			"type": "/v3/crosses/routex/geomarks",
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
	if !ok {
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
	for _, recipient := range recipients {
		switch recipient.Provider {
		case "iOS":
			fallthrough
		case "Android":
			body, err := json.Marshal(notifier.RequestArg{
				To:      recipient,
				CrossId: token.Cross.ID,
				From:    token.Identity,
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
	if token.TokenType == "user_token" {
		query.Set("user_id", fmt.Sprintf("%d", token.UserId))
	}
	token.Cross, err = m.platform.FindCross(int64(crossId), query)
	if err != nil {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return token, false
	}

	if token.TokenType == "user_token" {
		for _, inv := range token.Cross.Exfee.Invitations {
			if inv.Identity.UserID == token.UserId {
				token.Identity = inv.Identity
				return token, true
			}
		}
	} else {
		for _, inv := range token.Cross.Exfee.Invitations {
			if inv.Identity.ID == token.IdentityId {
				token.Identity = inv.Identity
				return token, true
			}
		}
	}

	m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
	return token, false
}
