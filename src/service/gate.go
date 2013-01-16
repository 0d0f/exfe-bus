package main

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"gobus"
	"model"
)

type Gate struct {
	client *gobus.Client
	url    string
	log    *logger.SubLogger
}

func NewGate(config *model.Config) (*Gate, error) {
	json := new(gobus.JSON)
	return &Gate{
		client: gobus.NewClient(json),
		url:    fmt.Sprintf("http://%s:%d/tokenmanager/token/%%s", config.ExfeService.Addr, config.ExfeService.Port),
		log:    config.Log.SubPrefix("gate"),
	}, nil
}

type tokenData struct {
	Type   string `json:"token_type"`
	UserID int64  `json:"user_id"`
}

func (g *Gate) Verify(token string) (int64, error) {
	url := fmt.Sprintf(g.url, token)
	var tk model.Token
	err := g.client.Do(url, "GET", nil, &tk)
	if err != nil {
		return 0, err
	}
	var data tokenData
	err = json.Unmarshal([]byte(tk.Data), &data)
	if err != nil {
		return 0, err
	}
	if data.Type != "user_token" {
		return 0, fmt.Errorf("%s is not user token", data.Type)
	}
	return data.UserID, nil
}
