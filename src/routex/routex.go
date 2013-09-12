package routex

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-pubsub"
	"github.com/googollee/go-rest"
	"logger"
	"math/rand"
	"model"
	"net/http"
	"net/url"
	"notifier"
	"os"
	"routex/model"
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
	SetGeomark     rest.Processor `path:"/geomarks/crosses/:cross_id/:mark_type/:kind.:mark_id" method:"PUT"`
	DeleteGeomark  rest.Processor `path:"/geomarks/crosses/:cross_id/:mark_type/:kind.:mark_id" method:"DELETE"`

	Stream  rest.Streaming `path:"/crosses/:cross_id" method:"WATCH"`
	Options rest.Processor `path:"/crosses/:cross_id" method:"OPTIONS"`

	SendNotification rest.Processor `path:"/notification/crosses/:cross_id" method:"POST"`

	rand            *rand.Rand
	routexRepo      rmodel.RoutexRepo
	breadcrumbCache rmodel.BreadcrumbCache
	breadcrumbsRepo rmodel.BreadcrumbsRepo
	geomarksRepo    rmodel.GeomarksRepo
	conversion      rmodel.GeoConversionRepo
	platform        *broker.Platform
	config          *model.Config
	tutorialDatas   map[int64][]rmodel.TutorialData
	pubsub          *pubsub.Pubsub
	castLocker      sync.RWMutex
}

func New(routexRepo rmodel.RoutexRepo, breadcrumbCache rmodel.BreadcrumbCache, breadcrumbsRepo rmodel.BreadcrumbsRepo, geomarksRepo rmodel.GeomarksRepo, conversion rmodel.GeoConversionRepo, platform *broker.Platform, config *model.Config) (*RouteMap, error) {
	tutorialDatas := make(map[int64][]rmodel.TutorialData)
	for _, userId := range config.TutorialBotUserIds {
		file := config.Routex.TutorialDataFile[fmt.Sprintf("%d", userId)]
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("can't find tutorial file %s for tutorial bot %d", file, userId)
		}
		var datas []rmodel.TutorialData
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
		pubsub:          pubsub.New(20),
	}
	return ret, nil
}

type UserCrossSetup struct {
	SaveBreadcrumbs bool `json:"save_breadcrumbs,omitempty"`
	AfterInSeconds  int  `json:"after_in_seconds,omitempty"`
}

func (m RouteMap) HandleSetUser(setup UserCrossSetup) {
	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	m.Vars()["user_id"] = fmt.Sprintf("%d", token.UserId)
	m.HandleSetUserInner(setup)
}

func (m RouteMap) HandleSetUserInner(setup UserCrossSetup) {
	userIdStr, crossIdStr := m.Vars()["user_id"], m.Vars()["cross_id"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, fmt.Errorf("invalid user id %s", userIdStr))
		return
	}
	crossId, err := strconv.ParseInt(crossIdStr, 10, 64)
	if err != nil {
		m.Error(http.StatusBadRequest, fmt.Errorf("invalid cross id %s", crossIdStr))
		return
	}
	if setup.AfterInSeconds == 0 {
		setup.AfterInSeconds = 60 * 60
	}
	m.switchWindow(userId, crossId, setup.SaveBreadcrumbs, setup.AfterInSeconds)
}

func (m RouteMap) HandleSearchRoutex(crossIds []int64) []rmodel.Routex {
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
		m.Error(http.StatusBadRequest, fmt.Errorf("invalid cross id %s", crossIdStr))
		return nil
	}
	endAt, err := m.breadcrumbsRepo.GetWindowEnd(userId, crossId)
	if err != nil {
		logger.ERROR("get user %d cross %d routex failed: %s", userId, crossId, err)
		m.Error(http.StatusInternalServerError, err)
		return nil
	}
	if endAt == 0 {
		return nil
	}
	ret := endAt >= time.Now().Unix()
	return &ret
}

func (m RouteMap) HandleStream(stream rest.Stream) {
	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	now := time.Now()
	endAt, err := m.breadcrumbsRepo.GetWindowEnd(token.UserId, int64(token.Cross.ID))
	if err != nil || endAt <= now.Unix() {
		forceOpen, ok := m.Request().URL.Query()["force_window_open"]
		if !ok {
			m.Error(http.StatusForbidden, fmt.Errorf("not in window"))
			return
		}
		after := 15 * 60
		if endAt == 0 {
			after = 60 * 60
		}
		if len(forceOpen) > 0 {
			if i, err := strconv.ParseInt(forceOpen[0], 10, 64); err == nil {
				after = int(i)
			}
		}
		endAt = now.Unix() + int64(after)
		m.switchWindow(token.UserId, int64(token.Cross.ID), true, after)
	}

	c := make(chan interface{}, 10)
	m.pubsub.Subscribe(m.publicName(int64(token.Cross.ID)), c)
	logger.DEBUG("streaming connected by user %d, cross %d", token.UserId, token.Cross.ID)
	defer func() {
		logger.DEBUG("streaming disconnect by user %d, cross %d", token.UserId, token.Cross.ID)
		m.pubsub.Unsubscribe(m.publicName(int64(token.Cross.ID)), c)
		close(c)
	}()

	willEnd := endAt - now.Unix()
	err = stream.Write(map[string]interface{}{
		"type":   "command",
		"action": "close_after",
		"args":   []interface{}{willEnd},
	})
	if err != nil {
		return
	}

	toMars := m.Request().URL.Query().Get("coordinate") == "mars"
	isTutorial := false
	if token.Cross.By.UserID == m.config.Routex.TutorialCreator {
		isTutorial = true
	}
	hasCreated := false

	m.WriteHeader(http.StatusOK)
	quit := make(chan int)
	defer func() { close(quit) }()

	for _, invitation := range token.Cross.Exfee.Invitations {
		userId := invitation.Identity.UserID
		route := rmodel.Geomark{
			Id:   m.breadcrumbsId(userId),
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
						positions := m.getTutorialData(time.Now(), userId, 1)
						if positions == nil {
							continue
						}
						route := rmodel.Geomark{
							Id:        m.breadcrumbsId(userId),
							Action:    "save_to_history",
							Type:      "route",
							Tags:      []string{"breadcrumbs"},
							Positions: positions,
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
			route.Positions = []rmodel.SimpleLocation{l}
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

	stream.SetWriteDeadline(time.Now().Add(broker.NetworkTimeout))
	err = stream.Write(map[string]string{
		"type":   "command",
		"action": "init_end",
	})
	if err != nil {
		return
	}
	fmt.Println("is tutorial:", isTutorial, "has created:", hasCreated)

	lastCheck := now.Unix()
	for {
		select {
		case d := <-c:
			if mark, ok := d.(rmodel.Geomark); ok {
				if isTutorial && !hasCreated {
					if mark.Id == m.breadcrumbsId(token.UserId) {
						locale, by := "", ""
						for _, i := range token.Cross.Exfee.Invitations {
							if i.Identity.UserID == token.UserId {
								locale, by = i.Identity.Locale, i.Identity.Id()
								break
							}
						}
						tutorialMark, err := m.setTutorial(mark.Positions[0].GPS[0], mark.Positions[0].GPS[1], token.UserId, int64(token.Cross.ID), locale, by)
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
				d = mark
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

		case <-time.After(time.Duration(endAt-time.Now().Unix()) * time.Second):
			newEndAt, err := m.breadcrumbsRepo.GetWindowEnd(token.UserId, int64(token.Cross.ID))
			if err != nil || newEndAt == 0 || newEndAt <= time.Now().Unix() {
				return
			}
			endAt = newEndAt
			err = stream.Write(map[string]interface{}{
				"type":   "command",
				"action": "close_after",
				"args":   []interface{}{endAt - time.Now().Unix()},
			})
			if err != nil {
				return
			}
		}
		if time.Now().Unix()-lastCheck > 60 {
			lastCheck = time.Now().Unix()
			newEndAt, err := m.breadcrumbsRepo.GetWindowEnd(token.UserId, int64(token.Cross.ID))
			if err != nil {
				logger.ERROR("can't set user %d cross %d: %s", token.UserId, token.Cross.ID, err)
				continue
			}
			endAt = newEndAt
			err = stream.Write(map[string]interface{}{
				"type":   "command",
				"action": "close_after",
				"args":   []interface{}{endAt - time.Now().Unix()},
			})
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
	token, ok := m.auth()
	if !ok {
		m.Error(http.StatusUnauthorized, m.DetailError(-1, "invalid token"))
		return
	}

	to := model.FromIdentityId(m.Request().URL.Query().Get("id"))
	var toInvitation *model.Invitation
	for _, inv := range token.Cross.Exfee.Invitations {
		if inv.Identity.Equal(to) {
			toInvitation = &inv
			break
		}
	}
	if toInvitation == nil {
		m.Error(http.StatusForbidden, fmt.Errorf("%s is not attend cross %d", to.Id(), token.Cross.ID))
		return
	}
	to = toInvitation.Identity

	recipients, err := m.platform.GetRecipientsById(to.Id())
	if err != nil {
		m.Error(http.StatusInternalServerError, err)
		return
	}

	m.update(token.UserId, int64(token.Cross.ID))

	arg := notifier.RequestArg{
		CrossId: token.Cross.ID,
		From:    token.Identity,
	}
	pushed := false
	for _, recipient := range recipients {
		switch recipient.Provider {
		case "iOS":
			fallthrough
		case "Android":
			arg.To = recipient
			m.sendRequest(arg)
			pushed = true
		}
	}
	if to.Provider == "wechat" {
		if ok, err := m.platform.CheckWechatFollowing(to.ExternalUsername); (err != nil || !ok) && !pushed {
			m.Error(http.StatusNotAcceptable, fmt.Errorf("can't find provider avaliable"))
		}
	}

	go func() {
		arg.To = to.ToRecipient()
		m.sendRequest(arg)
		for _, id := range toInvitation.Notifications {
			to := model.FromIdentityId(id)
			arg.To.ExternalUsername, arg.To.Provider = to.ExternalUsername, to.Provider
			m.sendRequest(arg)
		}
	}()
}

func (m *RouteMap) sendRequest(arg notifier.RequestArg) {
	body, err := json.Marshal(arg)
	if err != nil {
		logger.ERROR("can't marshal: %s with %+v", err, arg)
		return
	}
	url := fmt.Sprintf("http://%s:%d/v3/notifier/routex/request", m.config.ExfeService.Addr, m.config.ExfeService.Port)
	resp, err := broker.HttpResponse(broker.Http("POST", url, "applicatioin/json", body))
	if err != nil {
		logger.ERROR("post %s error: %s with %#v", url, err, string(body))
		return
	}
	resp.Close()
}

func (m RouteMap) switchWindow(userId, crossId int64, save bool, afterInSeconds int) {
	m.update(userId, crossId)
	if save {
		if err := m.breadcrumbsRepo.EnableCross(userId, crossId, afterInSeconds); err != nil {
			logger.ERROR("set user %d enable cross %d breadcrumbs repo failed: %s", userId, crossId, err)
		}
		if err := m.breadcrumbCache.EnableCross(userId, crossId, afterInSeconds); err != nil {
			logger.ERROR("set user %d enable cross %d breadcrumb cache failed: %s", userId, crossId, err)
		}
	} else {
		if err := m.breadcrumbsRepo.DisableCross(userId, crossId); err != nil {
			logger.ERROR("set user %d disable cross %d breadcrumbs repo failed: %s", userId, crossId, err)
		}
		if err := m.breadcrumbCache.DisableCross(userId, crossId); err != nil {
			logger.ERROR("set user %d disable cross %d breadcrumb cache failed: %s", userId, crossId, err)
		}
	}
}

func (m RouteMap) update(userId, crossId int64) {
	if err := m.routexRepo.Update(crossId); err != nil {
		logger.ERROR("update routex user %d cross %d error: %s", err)
	}
	m.platform.BotCrossUpdate("cross_id", fmt.Sprintf("%d", crossId), nil, model.Identity{})
}

func (m *RouteMap) auth() (rmodel.Token, bool) {
	m.Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	m.Header().Set("Access-Control-Allow-Credentials", "true")
	m.Header().Set("Cache-Control", "no-cache")

	var token rmodel.Token

	authData := m.Request().Header.Get("Exfe-Auth-Data")
	// if authData == "" {
	// 	authData = `{"token_type":"user_token","user_id":475,"signin_time":1374046388,"last_authenticate":1374046388}`
	// }

	if authData != "" {
		if err := json.Unmarshal([]byte(authData), &token); err != nil {
			return token, false
		}
	}

	crossIdStr, ok := m.Vars()["cross_id"]
	if !ok {
		if token.TokenType == "user_token" {
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
		return token, false
	}

	if token.Cross, err = m.platform.FindCross(int64(crossId), query); err != nil {
		return token, false
	}

	for _, inv := range token.Cross.Exfee.Invitations {
		switch token.TokenType {
		case "cross_access_token":
			if inv.Identity.ID == token.IdentityId {
				token.UserId = inv.Identity.UserID
				token.Identity = inv.Identity
				return token, true
			}
		case "user_token":
			if inv.Identity.UserID == token.UserId {
				token.Identity = inv.Identity
				return token, true
			}
		}
	}
	return token, false
}

func (m RouteMap) publicName(crossId int64) string {
	return fmt.Sprintf("routex:cross_%d", crossId)
}
