package main

import (
	"broker"
	"fmt"
	"github.com/googollee/go-broadcast"
	"github.com/googollee/go-rest"
	"here"
	"logger"
	"math/rand"
	"model"
	"net/http"
	"strings"
	"time"
)

type LiveService struct {
	rest.Service `prefix:"/v3/live"`

	card      rest.SimpleNode `route:"/cards" method:"POST"`
	streaming rest.Streaming  `route:"/streaming" method:"POST" end:""`

	platform  *broker.Platform
	config    *model.Config
	here      *here.Here
	rand      *rand.Rand
	tokens    map[string]bool
	broadcast map[string]*broadcast.Broadcast
}

func NewLive(config *model.Config, platform *broker.Platform) (*LiveService, error) {
	service := &LiveService{
		config:    config,
		here:      here.New(config.Here.Threshold, config.Here.SignThreshold, time.Duration(config.Here.TimeoutInSecond)*time.Second),
		rand:      rand.New(rand.NewSource(time.Now().Unix())),
		platform:  platform,
		broadcast: make(map[string]*broadcast.Broadcast),
	}

	go service.here.Serve()

	go func() {
		c := service.here.UpdateChannel()
		for {
			group := <-c
			cards := make([]here.Card, 0)
			if group.Name != "" {
				for _, d := range group.Data {
					cards = append(cards, d.Card)
				}
			}
			for token := range group.Data {
				if b, ok := service.broadcast[token]; ok {
					b.Send(cards)
				}
			}
		}
	}()

	return service, nil
}

func (h LiveService) Card(ctx rest.Context, data here.Data) {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", h.config.AccessDomain)
	ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Response().Header().Set("Cache-Control", "no-cache")

	var token string
	ctx.Bind("token", &token)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
	if token == "" {
		token = fmt.Sprintf("%04d", rand.Int31n(10000))
		if h.here.Exist(token) != nil {
			ctx.Return(http.StatusNotFound, "please wait and try again.")
			return
		}
		data.Card.Id = fmt.Sprintf("%032d", rand.Int31())
	} else if h.here.Exist(token) == nil {
		ctx.Return(http.StatusForbidden, "invalid token")
		return
	}
	data.Token = token
	remote := ctx.Request().RemoteAddr
	remotes := strings.Split(remote, ":")
	data.Traits = append(data.Traits, remotes[0])

	if data.Card.Avatar == "" {
		ids, err := h.platform.GetIdentity(data.Card.Identities)
		if err == nil {
			for _, id := range ids {
				if strings.Index(id.Avatar, "/v2/avatar/default?name=") < 0 {
					data.Card.Avatar = id.Avatar
					break
				}
			}
			logger.DEBUG("token %s can't find avatar", data.Token)
		} else {
			logger.DEBUG("get avatar failed: %s", err)
		}
	}

	err := h.here.Add(&data)
	logger.INFO("live", "add", "token", data.Token, "card", data.Card.Id, "name", data.Card.Name, "long", data.Longitude, "lat", data.Latitude, "acc", data.Accuracy, "traits", data.Traits)

	if err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}

	ctx.Render([]string{token, data.Card.Id})
}

func (h LiveService) Streaming(ctx rest.StreamContext) {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", h.config.AccessDomain)
	ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Response().Header().Set("Cache-Control", "no-cache")

	var token string
	ctx.Bind("token", &token)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
	group := h.here.Exist(token)
	if group == nil {
		ctx.Return(http.StatusForbidden, "invalid token")
		return
	}
	c := make(chan interface{}, 10)
	b, ok := h.broadcast[token]
	if !ok {
		b = broadcast.NewBroadcast(5)
		h.broadcast[token] = b
	}
	b.Register(c)
	defer func() {
		b.Unregister(c)
		close(c)
	}()

	cards := make([]here.Card, 0)
	for _, d := range group.Data {
		cards = append(cards, d.Card)
	}
	ctx.Render(cards)
	for {
		select {
		case d := <-c:
			cards, ok := d.([]here.Card)
			if !ok {
				continue
			}
			err := ctx.Render(cards)
			if err != nil || len(cards) == 0 {
				logger.INFO("live", "clear", "token", token)
				return
			}
		case <-time.After(broker.NetworkTimeout):
			err := ctx.Ping()
			if err != nil {
				return
			}
		}
	}
}
