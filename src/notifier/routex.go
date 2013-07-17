package notifier

import (
	"broker"
	"fmt"
	"formatter"
	"github.com/googollee/go-rest"
	"model"
	"net/http"
)

type Routex struct {
	rest.Service `prefix:"/v3/notifier/routex"`

	Request rest.Processor `path:"/request" method:"POST"`

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
	CrossId uint64          `json:"cross_id"`
	From    model.Identity  `json:"from"`

	Config *model.Config `json:"-"`
}

func (w Routex) HandleRequest(arg RequestArg) {
	arg.Config = w.config
	text, err := GenerateContent(w.localTemplate, "routex_request", arg.To.Provider, arg.To.Language, arg)
	if err != nil {
		w.Error(http.StatusInternalServerError, err)
		return
	}
	id, ontime, defaultOk, err := w.platform.Send(arg.To, text)
	if err != nil {
		w.Error(http.StatusInternalServerError, err)
		return
	}
	arg.To.Fallbacks = arg.To.Fallbacks[1:]
	WaitResponse(id, ontime, defaultOk, arg.To, fmt.Sprintf("http://%s:%d/v3/notifier/routex/request", w.config.ExfeService.Addr, w.config.ExfeService.Port), arg)
}
