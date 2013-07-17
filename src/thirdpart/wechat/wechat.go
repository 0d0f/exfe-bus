package wechat

import (
	"broker"
	"fmt"
	"model"
	"thirdpart"
	"time"
)

type Wechat struct {
	bots map[string]string
}

func New(config *model.Config) *Wechat {
	bots := make(map[string]string)
	for k, v := range config.Wechat {
		bots[k] = fmt.Sprintf("http://%s:%d/send?to=%%s", v.Addr, v.Port)
	}
	return &Wechat{
		bots: bots,
	}
}

func (w *Wechat) Provider() string {
	return "wechat"
}

func (w *Wechat) SetPosterCallback(callback thirdpart.Callback) (time.Duration, bool) {
	return 0, true
}

func (w *Wechat) Post(from, to, content string) (string, error) {
	bot, ok := w.bots[from]
	if !ok {
		return "", fmt.Errorf("can't find bot: %s", from)
	}
	resp, err := broker.Http("POST", fmt.Sprintf(bot, to), "text/plain", []byte(content))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return "", nil
}
