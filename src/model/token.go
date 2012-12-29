package model

type Token struct {
	Key       string `json:"key"`
	Data      string `json:"data"`
	IsExpired bool   `json:"is_expired"`
}
