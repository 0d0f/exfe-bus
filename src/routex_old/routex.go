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
	"time"
)

type RouteMap struct {
	rest.Service `prefix:"/v3/routex_old" mime:"application/json"`

	UpdateBreadcrums rest.Processor `path:"/crosses/:cross_id/breadcrumbs" method:"POST"`
	GetBreadcrums    rest.Processor `path:"/crosses/:cross_id/breadcrumbs" method:"GET"`
	UpdateGeomarks   rest.Processor `path:"/crosses/:cross_id/geomarks" method:"POST"`
	GetGeomarks      rest.Processor `path:"/crosses/:cross_id/geomarks" method:"GET"`
	Notification     rest.Streaming `path:"/crosses/:cross_id" method:"WATCH"`
	Options          rest.Processor `path:"/crosses/:cross_id" method:"OPTIONS"`
	SendRequest      rest.Processor `path:"/crosses/:cross_id/request" method:"POST"`

	breadcrumbsRepo BreadcrumbsRepo
	geomarksRepo    GeomarksRepo
	conversion      GeoConversionRepo
	platform        *broker.Platform
	config          *model.Config
	broadcasts      map[uint64]*broadcast.Broadcast
}

func New(breadcrumbsRepo BreadcrumbsRepo, geomarksRepo GeomarksRepo, conversion GeoConversionRepo, platform *broker.Platform, config *model.Config) *RouteMap {
	return &RouteMap{
		breadcrumbsRepo: breadcrumbsRepo,
		geomarksRepo:    geomarksRepo,
		conversion:      conversion,
		platform:        platform,
		config:          config,
		broadcasts:      make(map[uint64]*broadcast.Broadcast),
	}
}

func (m RouteMap) HandleUpdateBreadcrums(breadcrumb Location) map[string]string {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return nil
	}
	id := token.Identity.Id()

	_, _, acc, err := breadcrumb.GetGeo()
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return nil
	}
	if acc > 70 {
		m.Error(http.StatusBadRequest, fmt.Errorf("accuracy too large: %d", acc))
		return nil
	}

	earth := breadcrumb
	mars := breadcrumb
	if m.Request().URL.Query().Get("coordinate") == "mars" {
		breadcrumb.ToEarth(m.conversion)
		earth = breadcrumb
	} else {
		mars.ToMars(m.conversion)
	}

	breadcrumb.Timestamp = time.Now().Unix()
	breadcrumbs, err := m.breadcrumbsRepo.Load(id, token.Cross.ID)
	if err != nil {
		logger.ERROR("can't get breadcrumbs %s of cross %d: %s", id, token.Cross.ID, err)
	}
	if len(breadcrumbs) == 0 {
		logger.INFO("routex", "cross", token.Cross.ID, id, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy)
		if err := m.breadcrumbsRepo.Save(id, token.Cross.ID, breadcrumb); err != nil {
			logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.Cross.ID, err, breadcrumb)
			m.Error(http.StatusInternalServerError, err)
			return nil
		}
	} else {
		distance, err := breadcrumb.Distance(breadcrumbs[0])
		if err != nil {
			logger.INFO("routex", "cross", token.Cross.ID, id, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy, "err", err)
			if err := m.breadcrumbsRepo.Save(id, token.Cross.ID, breadcrumb); err != nil {
				logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.Cross.ID, err, breadcrumb)
				m.Error(http.StatusInternalServerError, err)
				return nil
			}
		} else if distance > 30 {
			logger.INFO("routex", "cross", token.Cross.ID, id, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy, "distance", distance)
			if err := m.breadcrumbsRepo.Save(id, token.Cross.ID, breadcrumb); err != nil {
				logger.ERROR("can't save repo %s of cross %d: %s with %+v", id, token.Cross.ID, err, breadcrumb)
				m.Error(http.StatusInternalServerError, err)
				return nil
			}
		} else {
			logger.INFO("routex", "cross", token.Cross.ID, id, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy, "distance", distance, "nosave")
		}
	}
	breadcrumbs = append([]Location{breadcrumb}, breadcrumbs...)
	earthLat, earthLng, _, _ := earth.GetGeo()
	marsLat, marsLng, _, _ := mars.GetGeo()
	ret := map[string]string{
		"earth_to_mars_latitude":  fmt.Sprintf("%.4f", marsLat-earthLat),
		"earth_to_mars_longitude": fmt.Sprintf("%.4f", marsLng-earthLng),
	}

	broadcast, ok := m.broadcasts[token.Cross.ID]
	if !ok {
		return ret
	}
	if broadcast == nil {
		delete(m.broadcasts, token.Cross.ID)
		return ret
	}
	d := make([]Location, len(breadcrumbs))
	for i := range breadcrumbs {
		d[i] = breadcrumbs[i]
	}
	broadcast.Send(map[string]interface{}{
		"type": "/v3/crosses/routex/breadcrumbs",
		"data": map[string][]Location{
			id: d,
		},
	})
	return ret
}

func (m RouteMap) HandleGetBreadcrums() map[string][]Location {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	toMars := m.Request().URL.Query().Get("coordinate") == "mars"
	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
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
		if toMars {
			for i := range breadcrumbs {
				breadcrumbs[i].ToMars(m.conversion)
			}
		}
		ret[id] = breadcrumbs
	}
	return ret
}

func (m RouteMap) HandleUpdateGeomarks(data []Location) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	if m.Request().URL.Query().Get("coordinate") == "mars" {
		for i := range data {
			data[i].ToEarth(m.conversion)
		}
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

func (m RouteMap) HandleGetGeomarks() []Location {
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
		if token.Cross.Place != nil && token.Cross.Place.Lng != "" && token.Cross.Place.Lat != "" {
			createdAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.CreatedAt)
			if err != nil {
				createdAt = time.Now()
			}
			updatedAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.UpdatedAt)
			if err != nil {
				updatedAt = time.Now()
			}
			destinaion := Location{
				Id:          "destination",
				Type:        "location",
				CreatedAt:   createdAt.Unix(),
				CreatedBy:   token.Cross.By.Id(),
				UpdatedAt:   updatedAt.Unix(),
				UpdatedBy:   token.Cross.By.Id(),
				Tags:        []string{"destination"},
				Title:       token.Cross.Place.Title,
				Description: token.Cross.Place.Description,
				Timestamp:   time.Now().Unix(),
				Longitude:   token.Cross.Place.Lng,
				Latitude:    token.Cross.Place.Lat,
			}
			data = []Location{destinaion}
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

	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
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

	toMars := m.Request().URL.Query().Get("coordinate") == "mars"

	for _, invitation := range token.Cross.Exfee.Invitations {
		ret := make(map[string][]Location)
		id := invitation.Identity.Id()
		breadcrumbs, err := m.breadcrumbsRepo.Load(id, token.Cross.ID)
		if err != nil {
			logger.ERROR("can't get breadcrumbs %s of cross %d: %s", id, token.Cross.ID, err)
			continue
		}
		if breadcrumbs == nil {
			continue
		}
		if toMars {
			for i := range breadcrumbs {
				breadcrumbs[i].ToMars(m.conversion)
			}
		}
		ret[id] = breadcrumbs
		err = stream.Write(map[string]interface{}{
			"type": "/v3/crosses/routex/breadcrumbs",
			"data": ret,
		})
		if err != nil {
			return
		}
	}

	data, err := m.geomarksRepo.Load(token.Cross.ID)
	if err == nil {
		if data == nil {
			data = make([]Location, 0)
			if token.Cross.Place != nil && token.Cross.Place.Lng != "" && token.Cross.Place.Lat != "" {
				createdAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.CreatedAt)
				if err != nil {
					createdAt = time.Now()
				}
				updatedAt, err := time.Parse("2006-01-02 15:04:05 -0700", token.Cross.UpdatedAt)
				if err != nil {
					updatedAt = time.Now()
				}
				destinaion := Location{
					Id:          "destination",
					Type:        "location",
					CreatedAt:   createdAt.Unix(),
					CreatedBy:   token.Cross.By.Id(),
					UpdatedAt:   updatedAt.Unix(),
					UpdatedBy:   token.Cross.By.Id(),
					Tags:        []string{"destination"},
					Title:       token.Cross.Place.Title,
					Description: token.Cross.Place.Description,
					Timestamp:   time.Now().Unix(),
					Longitude:   token.Cross.Place.Lng,
					Latitude:    token.Cross.Place.Lat,
				}
				data = append(data, destinaion)
			}
		}
		for i := range data {
			data[i].ToMars(m.conversion)
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
			if toMars {
				if data, ok := d.(map[string]interface{}); ok {
					sendData := data["data"]
					if breadcrumbs, ok := sendData.(map[string][]Location); ok {
						for k, v := range breadcrumbs {
							for i := range v {
								d := v[i]
								d.ToMars(m.conversion)
								v[i] = d
							}
							breadcrumbs[k] = v
						}
						sendData = breadcrumbs
					} else if marks, ok := sendData.([]Location); ok {
						for i := range marks {
							d := marks[i]
							d.ToMars(m.conversion)
							marks[i] = d
						}
						sendData = marks
					} else {
						continue
					}
					data["data"] = sendData
					d = data
				} else {
					continue
				}
			}
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
		return token, false
	}

	authData := m.Request().Header.Get("Exfe-Auth-Data")
	// if authData == "" {
	// 	// token: 345ac9296016c858a752a7e5fea35b7682fa69f922c4cefa30cfc22741da3109
	// 	authData = `{"token_type":"cross_access_token","cross_id":100758,"identity_id":907,"user_id":652,"created_time":1374636534,"updated_time":1374636534}`
	// }
	// logger.DEBUG("auth data: %s", authData)

	if err := json.Unmarshal([]byte(authData), &token); err != nil {
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

	return token, false
}
