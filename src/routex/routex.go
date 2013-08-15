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
	"os"
	"strconv"
	"sync"
	"time"
)

type RouteMap struct {
	rest.Service `prefix:"/v3/routex" mime:"application/json"`

	SearchRoutex rest.Processor `path:"/_inner/search/crosses" method:"POST"`
	SetUserInner rest.Processor `path:"/_inner/users/:user_id/crosses" method:"POST"`
	GetRoutex    rest.Processor `path:"/_inner/users/:user_id/crosses/:cross_id" method:"GET"`
	SetUser      rest.Processor `path:"/users/crosses" method:"POST"`

	UpdateBreadcrums      rest.Processor `path:"/breadcrumbs" method:"POST"`
	UpdateBreadcrumsInner rest.Processor `path:"/_inner/breadcrumbs/users/:user_id" method:"POST"`
	GetBreadcrums         rest.Processor `path:"/breadcrumbs/crosses/:cross_id" method:"GET"`
	GetUserBreadcrums     rest.Processor `path:"/breadcrumbs/crosses/:cross_id/users/:user_id" method:"GET"`

	GetGeomarks   rest.Processor `path:"/geomarks/crosses/:cross_id" method:"GET"`
	SetGeomark    rest.Processor `path:"/geomarks/crosses/:cross_id/:mark_type/:mark_id" method:"PUT"`
	DeleteGeomark rest.Processor `path:"/geomarks/crosses/:cross_id/:mark_type/:mark_id" method:"DELETE"`

	Stream  rest.Streaming `path:"/crosses/:cross_id" method:"WATCH"`
	Options rest.Processor `path:"/crosses/:cross_id" method:"OPTIONS"`

	SendNotification rest.Processor `path:"/crosses/notification/:cross_id/:identity_id" method:"POST"`

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

func (m *RouteMap) getTutorialData(currentTime time.Time, userId int64, number int) []SimpleLocation {
	data, ok := m.tutorialDatas[userId]
	if !ok {
		return nil
	}
	currentTime = currentTime.UTC()
	now := currentTime.Unix()
	todayTime, _ := time.Parse("2006-01-02 15:04:05", currentTime.Format("2006-01-02 00:00:00"))
	today := todayTime.Unix()

	oneDaySeconds := int64(24 * time.Hour / time.Second)
	totalPoint := len(data)
	fmt.Println("current second:", now-today)
	currentPoint := int((now - today) * int64(totalPoint) / oneDaySeconds)
	fmt.Println("current point:", currentPoint, "total:", totalPoint)

	var ret []SimpleLocation
	for ; number > 0; number-- {
		ret = append(ret, SimpleLocation{
			Timestamp: today + data[currentPoint].Offset,
			Accuracy:  data[currentPoint].Accuracy,
			Latitude:  data[currentPoint].Latitude,
			Longitude: data[currentPoint].Longitude,
		})
		if currentPoint > 0 {
			currentPoint--
		} else {
			currentPoint = totalPoint - 1
			today -= oneDaySeconds
		}
	}
	return ret
}

type UserCrossSetup struct {
	CrossId         int64 `json:"cross_id,omitempty"`
	SaveBreadcrumbs bool  `json:"save_breadcrumbs,omitempty"`
	AfterInSeconds  int   `json:"after_in_seconds,omitempty"`
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

	m.Vars()["user_id"] = fmt.Sprintf("%d", token.UserId)
	m.HandleSetUserInner(setup)
}

func (m RouteMap) HandleSetUserInner(setup []UserCrossSetup) {
	userIdStr := m.Vars()["user_id"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return
	}
	go func() {
		for _, s := range setup {
			if s.SaveBreadcrumbs {
				if s.AfterInSeconds == 0 {
					s.AfterInSeconds = 7200
				}
				if err := m.routexRepo.EnableCross(userId, s.CrossId, s.AfterInSeconds); err != nil {
					logger.ERROR("set user %d enable cross %d routex repo failed: %s", userId, s.CrossId, err)
				}
				if err := m.breadcrumbsRepo.EnableCross(userId, s.CrossId, s.AfterInSeconds); err != nil {
					logger.ERROR("set user %d enable cross %d breadcrumbs repo failed: %s", userId, s.CrossId, err)
				}
			} else {
				if err := m.routexRepo.DisableCross(userId, s.CrossId); err != nil {
					logger.ERROR("set user %d disable cross %d routex repo failed: %s", userId, s.CrossId, err)
				}
				if err := m.breadcrumbsRepo.DisableCross(userId, s.CrossId); err != nil {
					logger.ERROR("set user %d disable cross %d breadcrumbs repo failed: %s", userId, s.CrossId, err)
				}
			}
		}
	}()

	for _, s := range setup {
		if s.SaveBreadcrumbs {
			if s.AfterInSeconds == 0 {
				s.AfterInSeconds = 7200
			}
			if err := m.breadcrumbCache.EnableCross(userId, s.CrossId, s.AfterInSeconds); err != nil {
				logger.ERROR("set user %d enable cross %d breadcrumb cache failed: %s", userId, s.CrossId, err)
			}
		} else {
			if err := m.breadcrumbCache.DisableCross(userId, s.CrossId); err != nil {
				logger.ERROR("set user %d disable cross %d breadcrumb cache failed: %s", userId, s.CrossId, err)
			}
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

type BreadcrumbOffset struct {
	Latitude  float64 `json:"earth_to_mars_latitude"`
	Longitude float64 `json:"earth_to_mars_longitude"`
}

func (o BreadcrumbOffset) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"earth_to_mars_latitude":%.4f,"earth_to_mars_longitude":%.4f}`, o.Latitude, o.Longitude)), nil
}

func (m RouteMap) HandleUpdateBreadcrums(breadcrumbs []SimpleLocation) BreadcrumbOffset {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	var token Token
	var ret BreadcrumbOffset
	token, ok := m.auth()
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return ret
	}
	m.Vars()["user_id"] = fmt.Sprintf("%d", token.UserId)

	return m.HandleUpdateBreadcrumsInner(breadcrumbs)
}

func (m RouteMap) HandleUpdateBreadcrumsInner(breadcrumbs []SimpleLocation) BreadcrumbOffset {
	var ret BreadcrumbOffset

	userIdStr, breadcrumb := m.Vars()["user_id"], breadcrumbs[0]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, err)
		return ret
	}
	if breadcrumb.Accuracy > 70 {
		m.Error(http.StatusBadRequest, fmt.Errorf("accuracy too large: %d", breadcrumb.Accuracy))
		return ret
	}

	if m.Request().URL.Query().Get("coordinate") == "mars" {
		breadcrumb.ToEarth(m.conversion)
	}

	breadcrumb.Timestamp = time.Now().Unix()
	last, err := m.breadcrumbCache.Load(userId)
	distance := Distance(breadcrumb.Latitude, breadcrumb.Longitude, last.Latitude, last.Longitude)
	var crossIds []int64
	if err != nil || distance > 30 {
		logger.INFO("routex", "user", userId, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy)
		if err := m.breadcrumbCache.Save(userId, breadcrumb); err != nil {
			logger.ERROR("can't save cache %d: %s with %+v", userId, err, breadcrumb)
			m.Error(http.StatusInternalServerError, err)
			return ret
		}
		if crossIds, err = m.breadcrumbCache.SaveCross(userId, breadcrumb); err != nil {
			logger.ERROR("can't save cache %d: %s with %+v", userId, err, breadcrumb)
			m.Error(http.StatusInternalServerError, err)
			return ret
		}
		go func() {
			if err := m.breadcrumbsRepo.Save(userId, breadcrumb); err != nil {
				logger.ERROR("can't save user %d breadcrumb: %s with %+v", userId, err, breadcrumb)
			}
		}()
	} else {
		logger.INFO("routex", "user", userId, "breadcrumb", breadcrumb.Longitude, breadcrumb.Latitude, breadcrumb.Accuracy, "distance", distance, "nosave")
		if crossIds, err = m.breadcrumbCache.SaveCross(userId, breadcrumb); err != nil {
			logger.ERROR("can't save cache %d: %s with %+v", userId, err, breadcrumb)
			m.Error(http.StatusInternalServerError, err)
			return ret
		}
	}

	earth := breadcrumb
	mars := breadcrumb
	mars.ToMars(m.conversion)
	ret = BreadcrumbOffset{
		Latitude:  mars.Latitude - earth.Latitude,
		Longitude: mars.Longitude - earth.Longitude,
	}

	go func() {
		route := Geomark{
			Id:   fmt.Sprintf("%d@exfe", userId),
			Type: "route",
			Tags: []string{"breadcrumbs"},
		}
		for _, cross := range crossIds {
			m.castLocker.RLock()
			b, ok := m.crossCast[cross]
			m.castLocker.RUnlock()
			if !ok {
				continue
			}
			l, exist, err := m.breadcrumbCache.LoadCross(userId, cross)
			if err != nil {
				logger.ERROR("get user %d last breadcrumb from cross %d failed: %s", userId, cross, err)
				continue
			}
			if !exist {
				continue
			}
			route.Positions = []SimpleLocation{l}
			b.Send(route)
		}
	}()

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
	ret := make([]Geomark, 0)
	for _, invitation := range token.Cross.Exfee.Invitations {
		userId := invitation.Identity.UserID
		route := Geomark{
			Id:   fmt.Sprintf("%d@exfe", userId),
			Type: "route",
		}
		var err error
		if route.Positions, err = m.breadcrumbsRepo.Load(userId, int64(token.Cross.ID)); err != nil {
			logger.ERROR("can't get user %d breadcrumbs of cross %d: %s", userId, token.Cross.ID, err)
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

func (m RouteMap) HandleGetUserBreadcrums() Geomark {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	toMars, userIdStr := m.Request().URL.Query().Get("coordinate") == "mars", m.Vars()["user_id"]
	token, ok := m.auth()
	var ret Geomark
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return ret
	}
	var userId int64
	for _, invitation := range token.Cross.Exfee.Invitations {
		if fmt.Sprintf("%d", invitation.Identity.UserID) == userIdStr {
			userId = invitation.Identity.UserID
			break
		}
	}
	if userId == 0 {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return ret
	}
	var err error
	if ret.Positions, err = m.breadcrumbsRepo.Load(userId, int64(token.Cross.ID)); err != nil {
		logger.ERROR("can't get user %d breadcrumbs of cross %d: %s", userId, token.Cross.ID, err)
		return ret
	}
	ret.Id, ret.Type = fmt.Sprintf("%d", userId), "route"
	if toMars {
		ret.ToMars(m.conversion)
	}
	return ret
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
	data, err := m.geomarksRepo.Get(int64(token.Cross.ID))
	if err != nil {
		logger.ERROR("can't get route of cross %d: %s", token.Cross.ID, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	if data == nil {
		var lat, lng float64
		if token.Cross.Place != nil {
			if lng, err = strconv.ParseFloat(token.Cross.Place.Lng, 64); err != nil {
				token.Cross.Place = nil
			} else if lat, err = strconv.ParseFloat(token.Cross.Place.Lat, 64); err != nil {
				token.Cross.Place = nil
			}
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
				Icon:        "http://panda.0d0f.com/static/img/map_pin_blue@2x.png",
				Title:       token.Cross.Place.Title,
				Description: token.Cross.Place.Description,
				Longitude:   lng,
				Latitude:    lat,
			}
			go func() {
				m.geomarksRepo.Set(int64(token.Cross.ID), destinaion)
			}()
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

func (m RouteMap) HandleSetGeomark(mark Geomark) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	mark.Type, mark.Id, mark.UpdatedAt, mark.Action = m.Vars()["mark_type"], m.Vars()["mark_id"], time.Now().Unix(), ""
	if m.Request().URL.Query().Get("coordinate") == "mars" {
		mark.ToEarth(m.conversion)
	}
	if err := m.geomarksRepo.Set(int64(token.Cross.ID), mark); err != nil {
		m.Error(http.StatusInternalServerError, err)
		return
	}

	go func() {
		mark.Action = "update"
		m.castLocker.RLock()
		broadcast := m.crossCast[int64(token.Cross.ID)]
		m.castLocker.RUnlock()
		if broadcast != nil {
			broadcast.Send(mark)
		}
	}()

	return
}

func (m RouteMap) HandleDeleteGeomark() {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
	if !ok || token.Readonly {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	var mark Geomark
	mark.Type, mark.Id = m.Vars()["mark_type"], m.Vars()["mark_id"]
	if err := m.geomarksRepo.Delete(int64(token.Cross.ID), mark.Type, mark.Id); err != nil {
		m.Error(http.StatusInternalServerError, err)
		return
	}

	go func() {
		mark.Action = "delete"
		m.castLocker.RLock()
		broadcast := m.crossCast[int64(token.Cross.ID)]
		m.castLocker.RUnlock()
		if broadcast != nil {
			broadcast.Send(mark)
		}
	}()

	return
}

func (m RouteMap) HandleStream(stream rest.Stream) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
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

	for _, invitation := range token.Cross.Exfee.Invitations {
		userId := invitation.Identity.UserID
		route := Geomark{
			Id:   fmt.Sprintf("%d@exfe", userId),
			Type: "route",
			Tags: []string{"breadcrumbs"},
		}
		l, exist, err := m.breadcrumbCache.LoadCross(userId, int64(token.Cross.ID))
		if err != nil {
			logger.ERROR("can't get user %d breadcrumbs of cross %d: %s", userId, token.Cross.ID, err)
			continue
		}
		if !exist {
			continue
		}
		route.Positions = []SimpleLocation{l}
		if len(route.Positions) == 0 {
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

	marks, err := m.geomarksRepo.Get(int64(token.Cross.ID))
	if err == nil {
		if len(marks) == 0 {
			var lat, lng float64
			if token.Cross.Place != nil {
				if lng, err = strconv.ParseFloat(token.Cross.Place.Lng, 64); err != nil {
					token.Cross.Place = nil
				} else if lat, err = strconv.ParseFloat(token.Cross.Place.Lat, 64); err != nil {
					token.Cross.Place = nil
				}
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
					Icon:        "http://panda.0d0f.com/static/img/map_pin_blue@2x.png",
					Title:       token.Cross.Place.Title,
					Description: token.Cross.Place.Description,
					Longitude:   lng,
					Latitude:    lat,
				}
				go func() {
					m.geomarksRepo.Set(int64(token.Cross.ID), destinaion)
				}()
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

func (m RouteMap) HandleSendNotification() {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	token, ok := m.auth()
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
	// 	authData = `{"token_type":"user_token","user_id":475,"signin_time":1374046388,"last_authenticate":1374046388}`
	// }

	if authData != "" {
		if err := json.Unmarshal([]byte(authData), &token); err != nil {
			return token, false
		}
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
