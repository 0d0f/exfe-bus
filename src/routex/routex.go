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
	"sync"
	"time"
)

type RouteMap struct {
	rest.Service `prefix:"/v3/routex" mime:"application/json"`

	updateIdentity rest.SimpleNode `route:"/_inner/update_identity" method:"POST"`
	updateExfee    rest.SimpleNode `route:"/_inner/update_exfee" method:"POST"`
	searchRoutex   rest.SimpleNode `route:"/_inner/search/crosses" method:"POST"`
	getRoutex      rest.SimpleNode `route:"/_inner/users/:user_id/crosses/:cross_id" method:"GET"`
	setUser        rest.SimpleNode `route:"/users/crosses/:cross_id" method:"POST"`

	updateBreadcrums       rest.SimpleNode `route:"/breadcrumbs" method:"POST"`
	updateBreadcrumsInner  rest.SimpleNode `route:"/_inner/breadcrumbs/users/:user_id" method:"POST"`
	getBreadcrums          rest.SimpleNode `route:"/breadcrumbs/crosses/:cross_id" method:"GET"`
	getUserBreadcrums      rest.SimpleNode `route:"/breadcrumbs/crosses/:cross_id/users/:user_id" method:"GET"`
	getUserBreadcrumsInner rest.SimpleNode `route:"/_inner/breadcrumbs/users/:user_id" method:"GET"`

	searchGeomarks rest.SimpleNode `route:"/_inner/geomarks/crosses/:cross_id" method:"GET"`
	getGeomarks    rest.SimpleNode `route:"/geomarks/crosses/:cross_id" method:"GET"`
	setGeomark     rest.SimpleNode `route:"/geomarks/crosses/:cross_id/:mark_type/:kind.:mark_id" method:"PUT"`
	deleteGeomark  rest.SimpleNode `route:"/geomarks/crosses/:cross_id/:mark_type/:kind.:mark_id" method:"DELETE"`

	stream  rest.Streaming  `route:"/crosses/:cross_id" method:"WATCH"`
	options rest.SimpleNode `route:"/crosses/:cross_id" method:"OPTIONS"`

	sendNotification rest.SimpleNode `route:"/notification/crosses/:cross_id" method:"POST"`

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
	quit            chan int
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
		geomarksRepo:    geomarksRepo, conversion: conversion,
		platform:      platform,
		tutorialDatas: tutorialDatas,
		config:        config,
		pubsub:        pubsub.New(20),
		quit:          make(chan int),
	}
	go ret.tutorialGenerator()
	return ret, nil
}

func (m RouteMap) UpdateIdentity(ctx rest.Context, identity model.Identity) {
	id := rmodel.Identity{
		Identity: identity,
		Type:     "identity",
		Action:   "update",
	}
	m.pubsub.Publish(m.identityName(identity), id)
}

func (m RouteMap) UpdateExfee(ctx rest.Context, invitations model.Invitation) {
	var crossId int64
	var action string
	ctx.Bind("cross_id", &crossId)
	ctx.Bind("action", &action)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	if action != "join" && action != "remove" {
		ctx.Return(http.StatusBadRequest, "invalid action: %s", action)
		return
	}
	id := rmodel.Invitation{
		Identity:      invitations.Identity,
		Notifications: invitations.Notifications,
		Type:          "invitation",
		Action:        action,
	}
	m.pubsub.Publish(m.publicName(crossId), id)
}

type UserCrossSetup struct {
	SaveBreadcrumbs bool `json:"save_breadcrumbs,omitempty"`
	AfterInSeconds  int  `json:"after_in_seconds,omitempty"`
}

func (m RouteMap) SetUser(ctx rest.Context, setup UserCrossSetup) {
	token, ok := m.auth(ctx)
	if !ok {
		ctx.Return(http.StatusUnauthorized, "invalid token")
		return
	}

	var crossId int64
	ctx.Bind("cross_id", &crossId)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	if setup.AfterInSeconds == 0 {
		setup.AfterInSeconds = 60 * 60
	}
	m.switchWindow(crossId, token.Identity, setup.SaveBreadcrumbs, setup.AfterInSeconds)
}

func (m RouteMap) SearchRoutex(ctx rest.Context, crossIds []int64) {
	ret, err := m.routexRepo.Search(crossIds)
	if err != nil {
		logger.ERROR("search for route failed: %s with %+v", err, crossIds)
		ctx.Return(http.StatusInternalServerError, err)
		return
	}
	ctx.Render(ret)
}

type RoutexInfo struct {
	InWindow *bool            `json:"in_window"`
	Objects  []rmodel.Geomark `json:"objects"`
}

func (m RouteMap) GetRoutex(ctx rest.Context) {
	var userId, crossId int64
	ctx.Bind("cross_id", &crossId)
	ctx.Bind("user_id", &userId)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	endAt, err := m.breadcrumbsRepo.GetWindowEnd(userId, crossId)
	if err != nil {
		logger.ERROR("get user %d cross %d routex failed: %s", userId, crossId, err)
		ctx.Return(http.StatusInternalServerError, err)
		return
	}
	ret := RoutexInfo{}
	if endAt != 0 {
		ret.InWindow = new(bool)
		*ret.InWindow = endAt >= time.Now().Unix()
	}
	query := make(url.Values)
	query.Set("user_id", fmt.Sprintf("%d", userId))
	cross, err := m.platform.FindCross(crossId, query)
	if err == nil {
		ret.Objects = m.getObjects(cross, true)
	} else {
		logger.ERROR("get user %d cross %d failed: %s", userId, crossId, err)
		ctx.Return(http.StatusInternalServerError, err)
		return
	}
	ctx.Render(ret)
}

func (m RouteMap) Stream(ctx rest.StreamContext) {
	token, ok := m.auth(ctx)
	if !ok {
		ctx.Return(http.StatusUnauthorized, "invalid token")
		return
	}
	var forceOpen bool
	var coordinate string
	ctx.Bind("force_window_open", &forceOpen)
	ctx.Bind("coordinate", &coordinate)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}

	now := time.Now()
	endAt, err := m.breadcrumbsRepo.GetWindowEnd(token.UserId, int64(token.Cross.ID))
	if err != nil || endAt <= now.Unix() {
		if !forceOpen {
			ctx.Return(http.StatusForbidden, "not in window")
			return
		}
		after := 15 * 60
		if endAt == 0 {
			after = 60 * 60
		}
		var openAfter int
		ctx.BindReset()
		ctx.Bind("force_window_open", &openAfter)
		if ctx.BindError() == nil {
			after = openAfter
		}
		endAt = now.Unix() + int64(after)
		m.switchWindow(int64(token.Cross.ID), token.Identity, true, after)
	}

	c := make(chan interface{}, 10)
	m.pubsub.Subscribe(m.publicName(int64(token.Cross.ID)), c)
	if token.Cross.By.UserID == m.config.Routex.TutorialCreator {
		m.pubsub.Subscribe(m.tutorialName(), c)
	}
	for _, inv := range token.Cross.Exfee.Invitations {
		m.pubsub.Subscribe(m.identityName(inv.Identity), c)
	}
	logger.DEBUG("streaming connected by user %d, cross %d", token.UserId, token.Cross.ID)
	defer func() {
		logger.DEBUG("streaming disconnect by user %d, cross %d", token.UserId, token.Cross.ID)
		m.pubsub.UnsubscribeAll(c)
		close(c)
	}()

	willEnd := endAt - now.Unix()
	err = ctx.Render(map[string]interface{}{
		"type":   "command",
		"action": "close_after",
		"args":   []interface{}{willEnd},
	})
	if err != nil {
		return
	}

	toMars := coordinate == "mars"
	isTutorial := false
	if token.Cross.By.UserID == m.config.Routex.TutorialCreator {
		isTutorial = true
	}
	hasCreated := false

	ctx.Return(http.StatusOK)
	quit := make(chan int)
	defer func() { close(quit) }()

	for _, mark := range m.getObjects(token.Cross, toMars) {
		if isTutorial && !hasCreated && !mark.IsBreadcrumbs() {
			hasCreated = true
		}
		if err := ctx.Render(mark); err != nil {
			return
		}
	}

	ctx.SetWriteDeadline(time.Now().Add(broker.NetworkTimeout))
	if err := ctx.Render(map[string]string{"type": "command", "action": "init_end"}); err != nil {
		return
	}

	lastCheck := now.Unix()
	for ctx.Ping() == nil {
		select {
		case d := <-c:
			switch data := d.(type) {
			case rmodel.Geomark:
				if isTutorial && !hasCreated {
					if data.Id == m.breadcrumbsId(token.UserId) {
						locale, by := "", ""
						for _, i := range token.Cross.Exfee.Invitations {
							if i.Identity.UserID == token.UserId {
								locale, by = i.Identity.Locale, i.Identity.Id()
								break
							}
						}
						tutorialMark, err := m.setTutorial(data.Positions[0].GPS[0], data.Positions[0].GPS[1], token.UserId, int64(token.Cross.ID), locale, by)
						if err != nil {
							logger.ERROR("create tutorial geomark for user %d in cross %d failed: %s", token.UserId, token.Cross.ID, err)
						} else {
							hasCreated = true
							if toMars {
								tutorialMark.ToMars(m.conversion)
							}
							err := ctx.Render(tutorialMark)
							if err != nil {
								return
							}
						}
					}
				}
				if toMars {
					data.ToMars(m.conversion)
				}
				d = data
			case rmodel.Identity:
				switch data.Action {
				case "join":
					if token.Cross.Exfee.Join(data.Identity) {
						m.pubsub.Subscribe(m.identityName(data.Identity), c)
					}
				case "remove":
					if token.Cross.Exfee.Remove(data.Identity) {
						m.pubsub.Unsubscribe(m.identityName(data.Identity), c)
					}
				}
			}
			ctx.SetWriteDeadline(time.Now().Add(broker.NetworkTimeout))
			err := ctx.Render(d)
			if err != nil {
				return
			}
		case <-time.After(broker.NetworkTimeout):
		case <-time.After(time.Duration(endAt-time.Now().Unix()) * time.Second):
			newEndAt, err := m.breadcrumbsRepo.GetWindowEnd(token.UserId, int64(token.Cross.ID))
			if err != nil || newEndAt == 0 || newEndAt <= time.Now().Unix() {
				return
			}
			endAt = newEndAt
			err = ctx.Render(map[string]interface{}{
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
			err = ctx.Render(map[string]interface{}{
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

func (m RouteMap) Options(ctx rest.Context) {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Response().Header().Set("Cache-Control", "no-cache")

	ctx.Return(http.StatusNoContent)
}

func (m RouteMap) SendNotification(ctx rest.Context) {
	token, ok := m.auth(ctx)
	if !ok {
		ctx.Return(http.StatusUnauthorized, "invalid token")
		return
	}

	var id string
	ctx.Bind("id", &id)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	to := model.FromIdentityId(id)
	var toInvitation *model.Invitation
	for _, inv := range token.Cross.Exfee.Invitations {
		if inv.Identity.Equal(to) {
			toInvitation = &inv
			break
		}
	}
	if toInvitation == nil {
		ctx.Return(http.StatusForbidden, "%s is not attend cross %d", to.Id(), token.Cross.ID)
		return
	}
	to = toInvitation.Identity

	recipients, err := m.platform.GetRecipientsById(to.Id())
	if err != nil {
		ctx.Return(http.StatusInternalServerError, err)
		return
	}

	m.update(int64(token.Cross.ID), token.Identity)

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
			ctx.Return(http.StatusNotAcceptable, "can't find provider avaliable")
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

func (m *RouteMap) getObjects(cross model.Cross, toMars bool) []rmodel.Geomark {
	isTutorial := false
	if cross.By.UserID == m.config.Routex.TutorialCreator {
		isTutorial = true
	}

	var ret []rmodel.Geomark
	breadcrumbs, err := m.breadcrumbCache.LoadAllCross(int64(cross.ID))
	now := time.Now()
	if isTutorial {
		for _, id := range m.config.TutorialBotUserIds {
			l := m.getTutorialData(now, id, 1)
			if len(l) > 0 {
				breadcrumbs[id] = l[0]
			}
		}
	}

	users := make(map[int64]bool)
	for _, inv := range cross.Exfee.Invitations {
		users[inv.Identity.UserID] = true
	}
	if err == nil {
		for userId, l := range breadcrumbs {
			if !users[userId] {
				if err := m.breadcrumbCache.RemoveCross(userId, int64(cross.ID)); err != nil {
					logger.ERROR("remove user %d cross %d breadcrumb error: %s", userId, cross.ID, err)
				}
				continue
			}
			mark := m.breadcrumbsToGeomark(userId, 1, []rmodel.SimpleLocation{l})
			if toMars {
				mark.ToMars(m.conversion)
			}
			ret = append(ret, mark)
		}
	} else {
		logger.ERROR("can't get current breadcrumb of cross %d: %s", cross.ID, err)
	}

	marks, err := m.getGeomarks_(cross, toMars)
	if err == nil {
		ret = append(ret, marks...)
	} else {
		logger.ERROR("can't get route of cross %d: %s", cross.ID, err)
	}

	return ret
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

func (m RouteMap) switchWindow(crossId int64, identity model.Identity, save bool, afterInSeconds int) {
	m.update(crossId, identity)
	if save {
		if err := m.breadcrumbsRepo.EnableCross(identity.UserID, crossId, afterInSeconds); err != nil {
			logger.ERROR("set user %d enable cross %d breadcrumbs repo failed: %s", identity.UserID, crossId, err)
		}
		if err := m.breadcrumbCache.EnableCross(identity.UserID, crossId, afterInSeconds); err != nil {
			logger.ERROR("set user %d enable cross %d breadcrumb cache failed: %s", identity.UserID, crossId, err)
		}
	} else {
		if err := m.breadcrumbsRepo.DisableCross(identity.UserID, crossId); err != nil {
			logger.ERROR("set user %d disable cross %d breadcrumbs repo failed: %s", identity.UserID, crossId, err)
		}
		if err := m.breadcrumbCache.DisableCross(identity.UserID, crossId); err != nil {
			logger.ERROR("set user %d disable cross %d breadcrumb cache failed: %s", identity.UserID, crossId, err)
		}
	}
}

func (m RouteMap) update(crossId int64, by model.Identity) {
	if err := m.routexRepo.Update(crossId); err != nil {
		logger.ERROR("update routex user %d cross %d error: %s", err)
	}
	cross := make(map[string]interface{})
	cross["widgets"] = []map[string]string{
		map[string]string{"type": "routex"},
	}
	m.platform.BotCrossUpdate("cross_id", fmt.Sprintf("%d", crossId), cross, by)
}

func (m *RouteMap) auth(ctx rest.Context) (rmodel.Token, bool) {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", m.config.AccessDomain)
	ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Response().Header().Set("Cache-Control", "no-cache")

	defer ctx.BindReset()

	var token rmodel.Token

	authData := ctx.Request().Header.Get("Exfe-Auth-Data")
	// if authData == "" {
	// 	authData = `{"token_type":"user_token","user_id":475,"signin_time":1374046388,"last_authenticate":1374046388}`
	// }

	if authData != "" {
		if err := json.Unmarshal([]byte(authData), &token); err != nil {
			return token, false
		}
	}

	var crossIdFlag bool
	ctx.Bind("cross_id", &crossIdFlag)
	if ctx.BindError() != nil || !crossIdFlag {
		if token.TokenType == "user_token" {
			return token, true
		}
		return token, false
	}
	var crossId int64
	ctx.Bind("cross_id", &crossId)
	if err := ctx.BindError(); err != nil {
		return token, false
	}

	query := make(url.Values)
	switch token.TokenType {
	case "user_token":
		query.Set("user_id", fmt.Sprintf("%d", token.UserId))
	case "cross_access_token":
		if int64(token.CrossId) != crossId {
			return token, false
		}
	default:
		return token, false
	}

	var err error
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

func (m RouteMap) tutorialName() string {
	return "routex:tutorial:data"
}

func (m RouteMap) identityName(identity model.Identity) string {
	return fmt.Sprintf("routex:identity:%s", identity.Id())
}

func (m RouteMap) tutorialGenerator() {
	for {
		select {
		case <-m.quit:
			return
		case <-time.After(time.Second * 10):
			now := time.Now()
			for userId := range m.tutorialDatas {
				positions := m.getTutorialData(now, userId, 1)
				if len(positions) == 0 {
					continue
				}
				mark := m.breadcrumbsToGeomark(userId, 1, positions)
				m.pubsub.Publish(m.tutorialName(), mark)
			}
		}
	}
}
