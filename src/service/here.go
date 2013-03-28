package main

import (
	"fmt"
	"github.com/googollee/go-rest"
	"here"
	"model"
	"net/http"
	"strings"
	"time"
)

type HereService struct {
	rest.Service `prefix:"/v3/here"`

	Users     rest.Processor `path:"/users" method:"POST"`
	Streaming rest.Streaming `path:"/streaming" method:"GET" end:"" timeout:"60"`

	config *model.Config
	here   *here.Here
}

func (h HereService) Users_(user here.User) {
	remote := h.Request().RemoteAddr
	remotes := strings.Split(remote, ":")
	user.Traits = append(user.Traits, remotes[0])
	h.here.Add(user)
}

func (h HereService) Streaming_() string {
	token := h.Request().URL.Query().Get("token")
	if token == "" {
		h.Error(http.StatusBadRequest, fmt.Errorf("invalid token"))
	}
	return token
}

func NewHere(config *model.Config) (http.Handler, error) {
	service := new(HereService)
	service.config = config
	service.here = here.New(config.Here.Threshold, config.Here.SignThreshold, time.Duration(config.Here.TimeoutInSecond)*time.Second)

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
		}
	}()

	return rest.New(service)
}
