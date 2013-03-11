package model

type Token struct {
	Key       string `json:"key"`
	Data      string `json:"data"`
	TouchedAt string `json:"touched_at"`
	ExpireAt  int64  `json:"expire_at"`
}
