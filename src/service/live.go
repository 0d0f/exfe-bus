package main

import (
	"fmt"
	"github.com/googollee/go-rest"
	"here"
	"math/rand"
	"model"
	"net/http"
	"strings"
	"time"
)

type LiveService struct {
	rest.Service `prefix:"/v3/live"`

	Card      rest.Processor `path:"/cards" method:"POST"`
	Streaming rest.Streaming `path:"/streaming" method:"GET" end:"" timeout:"60"`

	config *model.Config
	here   *here.Here
	rand   *rand.Rand
	tokens map[string]bool
}

func (h LiveService) Card_(data here.Data) string {
	h.Header().Set("access-control-allow-origin", h.config.AccessDomain)
	h.Header().Set("access-control-allow-credentials", "True")
	token := h.Request().URL.Query().Get("token")
	if token == "" {
		token = fmt.Sprintf("%04d", rand.Int31n(10000))
		if h.tokens[token] {
			h.Error(http.StatusNotFound, fmt.Errorf("please wait and try again."))
			return ""
		}
		h.tokens[token] = true
	} else if !h.tokens[token] {
		h.Error(http.StatusForbidden, fmt.Errorf("invalid token"))
		return ""
	}
	data.Token = token
	data.Card.IsMe = false
	remote := h.Request().RemoteAddr
	remotes := strings.Split(remote, ":")
	data.Traits = append(data.Traits, remotes[0])
	h.here.Add(data)

	return token
}

func (h LiveService) Streaming_() string {
	h.Header().Set("access-control-allow-origin", h.config.AccessDomain)
	h.Header().Set("access-control-allow-credentials", "True")
	token := h.Request().URL.Query().Get("token")
	if !h.tokens[token] {
		h.Error(http.StatusForbidden, fmt.Errorf("invalid token"))
	}
	return token
}

func NewLive(config *model.Config) (http.Handler, error) {
	service := &LiveService{
		config: config,
		here:   here.New(config.Here.Threshold, config.Here.SignThreshold, time.Duration(config.Here.TimeoutInSecond)*time.Second),
		rand:   rand.New(rand.NewSource(time.Now().Unix())),
		tokens: make(map[string]bool),
	}

	go func() {
		c := service.here.UpdateChannel()
		for {
			token := <-c
			group := service.here.UserInGroup(token)
			cards := make([]here.Card, 0)
			if group != nil {
				if _, ok := group.Data[token]; ok {
					for k, d := range group.Data {
						card := d.Card
						if k == token {
							card.IsMe = true
						}
						cards = append(cards, card)
					}
				}
			}
			service.Streaming.Feed(token, cards)
			if len(cards) == 0 {
				service.Streaming.Disconnect(token)
				delete(service.tokens, token)
			}
		}
	}()

	return rest.New(service)
}
