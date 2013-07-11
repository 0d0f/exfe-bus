package notifier

import (
	"broker"
	"formatter"
	"model"
)

type Wechat struct {
	localTemplate *formatter.LocalTemplate
	config        *model.Config
	platform      *broker.Platform
}

func NewWechat(localTemplate *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) *Wechat {
	return &Wechat{
		localTemplate: localTemplate,
		config:        config,
		platform:      platform,
	}
}

func (w *Wechat) RoutexNotice(to model.Recipient) error {
	text, err := GenerateContent(w.localTemplate, "wechat_routex", to.Provider, to.Language, to)
	if err != nil {
		return err
	}
	_, err = w.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}
