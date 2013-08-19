package wechat

import (
	"broker"
	"fmt"
	"model"
	"thirdpart"
	"time"
)

type Wechat struct {
	url string
}

func New(config *model.Config) *Wechat {
	u := fmt.Sprintf("%s/v3/bus/sendWechatMessage", config.SiteApi)
	return &Wechat{
		url: u,
	}
}

func (w *Wechat) Provider() string {
	return "wechat"
}

func (w *Wechat) SetPosterCallback(callback thirdpart.Callback) (time.Duration, bool) {
	return 0, true
}

func (w *Wechat) Post(from, to, content string) (string, error) {
	resp, err := broker.Http("POST", w.url, "application/javascript", []byte(content))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return "", nil
}
