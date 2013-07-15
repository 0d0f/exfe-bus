package notifier

import (
	"broker"
	"formatter"
	"model"
)

type Routex struct {
	localTemplate *formatter.LocalTemplate
	config        *model.Config
	platform      *broker.Platform
}

func NewRoutex(localTemplate *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) *Routex {
	return &Routex{
		localTemplate: localTemplate,
		config:        config,
		platform:      platform,
	}
}

type RequestArg struct {
	To      model.Recipient `json:"to"`
	CrossId int64           `json:"cross_id"`

	Config *model.Config `json:"-"`
}

func (w *Routex) Request(arg RequestArg) error {
	arg.Config = w.config
	text, err := GenerateContent(w.localTemplate, "routex_request", arg.To.Provider, arg.To.Language, arg)
	if err != nil {
		return err
	}
	_, err = w.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}
