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

func (h LiveService) Card_(data here.Data) []string {
	h.Header().Set("access-control-allow-origin", h.config.AccessDomain)
	h.Header().Set("access-control-allow-credentials", "True")

	token := h.Request().URL.Query().Get("token")
	if token == "" {
		token = fmt.Sprintf("%04d", rand.Int31n(10000))
		if h.here.Exist(token) {
			h.Error(http.StatusNotFound, fmt.Errorf("please wait and try again."))
			return nil
		}
		data.Card.Id = fmt.Sprintf("%032d", rand.Int31())
		h.config.Log.Debug("new card: %s", token)
	} else if !h.here.Exist(token) {
		h.Error(http.StatusForbidden, fmt.Errorf("invalid token"))
		return nil
	}
	data.Token = token
	remote := h.Request().RemoteAddr
	remotes := strings.Split(remote, ":")
	data.Traits = append(data.Traits, remotes[0])

	err := h.here.Add(&data)

	if err != nil {
		h.Error(http.StatusBadRequest, err)
		return nil
	}
	h.config.Log.Debug("update card: %s", token)

	return []string{token, data.Card.Id}
}

func (h LiveService) Streaming_() string {
	h.Header().Set("access-control-allow-origin", h.config.AccessDomain)
	h.Header().Set("access-control-allow-credentials", "True")
	token := h.Request().URL.Query().Get("token")
	if !h.here.Exist(token) {
		h.Error(http.StatusForbidden, fmt.Errorf("invalid token"))
		return ""
	}
	return token
}

func NewLive(config *model.Config) (http.Handler, error) {
	service := &LiveService{
		config: config,
		here:   here.New(config.Here.Threshold, config.Here.SignThreshold, time.Duration(config.Here.TimeoutInSecond)*time.Second),
		rand:   rand.New(rand.NewSource(time.Now().Unix())),
	}

	go service.here.Serve()

	go func() {
		c := service.here.UpdateChannel()
		for {
			group := <-c
			var cards []here.Card
			if group.Name != "" {
				for _, d := range group.Data {
					cards = append(cards, d.Card)
				}
			}
			for token := range group.Data {
				service.Streaming.Feed(token, cards)
				if group.Name == "" {
					service.Streaming.Disconnect(token)
				}
			}
		}
	}()

	return rest.New(service)
}
