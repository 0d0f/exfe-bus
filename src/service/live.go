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

	Card      rest.Processor `path:"/card" method:"POST"`
	Streaming rest.Streaming `path:"/streaming" method:"GET" end:"" timeout:"60"`

	config *model.Config
	here   *here.Here
	rand   *rand.Rand
	tokens map[string]bool
}

func (h LiveService) Card_(user here.User) string {
	token := h.Request().URL.Query().Get("token")
	if token == "" {
		token = fmt.Sprintf("%04d", rand.Int31n(10000))
		if h.tokens[token] {
			h.Error(http.StatusNotFound, fmt.Errorf("please wait and try again."))
			return ""
		}
	}
	if !h.tokens[token] {
		h.Error(http.StatusForbidden, fmt.Errorf("invalid token"))
		return ""
	}
	user.Id = token
	remote := h.Request().RemoteAddr
	remotes := strings.Split(remote, ":")
	user.Traits = append(user.Traits, remotes[0])
	h.here.Add(user)

	return token
}

func (h LiveService) Streaming_() string {
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
			id := <-c
			group := service.here.UserInGroup(id)
			users := make(map[string]*here.User)
			if group != nil {
				if _, ok := group.Users[id]; ok {
					users = group.Users
				}
			}
			service.Streaming.Feed(id, users)
			if len(users) == 0 {
				service.Streaming.Disconnect(id)
			}
		}
	}()

	return rest.New(service)
}