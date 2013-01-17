package main

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"gobus"
	"model"
	"net/http"
	"streaming"
	"time"
)

type Streaming struct {
	streaming *streaming.Streaming
	config    *model.Config
	gate      *Gate
	log       *logger.SubLogger
}

func NewStreaming(config *model.Config, gate *Gate) (*Streaming, error) {
	return &Streaming{
		streaming: streaming.New(time.Second),
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
		log.Debug("invalid token: %s", token)
		http.Error(w, fmt.Sprintf("token(%s) invalid", token), http.StatusForbidden)
		return
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support streaming", http.StatusInternalServerError)
		return
	}

	conn, buf, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer conn.Close()

	log.Info("connect to %d", userID)
	defer func() { log.Info("disconnect: %s", err) }()

	err = s.streaming.Connect(fmt.Sprintf("%d", userID), conn, buf)
}

func (s *Streaming) Provider() string {
	return "streaming"
}

func (s *Streaming) Send(to *model.Recipient, privateMessage string, publicMessage string, data *model.InfoData) (string, error) {
	err := s.streaming.Feed(fmt.Sprintf("%d", to.UserID), privateMessage)
	if err != nil {
		s.log.Err("send error: %s", err)
	}
	return "1", err
}

func (s *Streaming) SetRoute(r gobus.RouteCreater) error {
	json := new(gobus.JSON)
	return r().Methods("POST").Path("/streaming").HandlerMethod(json, s, "Receive")
}

type receiveArg struct {
	To *model.Recipient `json:"to"`
}

func (s *Streaming) Receive(params map[string]string, data map[string]interface{}) (int, error) {
	content, err := json.Marshal(data)
	if err != nil {
		return -1, err
	}
	var to receiveArg
	err = json.Unmarshal(content, &to)
	if err != nil || to.To == nil {
		return -1, fmt.Errorf("field 'to' invalid")
	}
	_, err = s.Send(to.To, string(content), "", nil)
	return 1, err
}
