package main

import (
	"fmt"
	"github.com/googollee/go-logger"
	"model"
	"net/http"
	"streaming"
)

type Streaming struct {
	streaming *streaming.Streaming
	config    *model.Config
	gate      *Gate
	log       *logger.SubLogger
}

func NewStreaming(config *model.Config, gate *Gate) (*Streaming, error) {
	return &Streaming{
		streaming: streaming.New(),
		config:    config,
		gate:      gate,
		log:       config.Log.SubPrefix("streaming"),
	}, nil
}

func (s *Streaming) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := s.log.SubCode()
	token := r.URL.Query().Get("token")
	userID, err := s.gate.Verify(token)
	if err != nil {
		log.Info("connect to %d", userID)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	log.Info("connect to %d", userID)
	s.streaming.Connect(fmt.Sprintf("%d", userID), w)
}

func (s *Streaming) Provider() string {
	return "streaming"
}

func (s *Streaming) Send(to *model.Recipient, privateMessage string, publicMessage string, data *model.InfoData) (string, error) {
	err := s.streaming.Feed(fmt.Sprintf("%d", to.UserID), privateMessage)
	return "1", err
}
