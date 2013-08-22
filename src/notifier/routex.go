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
	domain        string
}

func NewRoutex(localTemplate *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) *Routex {
	return &Routex{
		localTemplate: localTemplate,
		config:        config,
		platform:      platform,
		domain:        fmt.Sprintf("http://%s:%d", config.ExfeService.Addr, config.ExfeService.Port),
	}
}

type RequestArg struct {
	To      model.Recipient `json:"to"`
	CrossId uint64          `json:"cross_id"`
	From    model.Identity  `json:"from"`
	Cross   model.Cross     `json:"cross"`

	Config *model.Config `json:"-"`
}

func (w Routex) HandleRequest(arg RequestArg) {
	arg.Config = w.config
	var err error
	if arg.Cross, err = w.platform.FindCross(int64(arg.CrossId), nil); err != nil {
		w.Error(http.StatusBadRequest, err)
		return
	}

	go SendAndSave(w.localTemplate, w.platform, &arg.To, arg, "routex_request", w.domain+"/v3/notifier/routex/request", &arg)
	w.WriteHeader(http.StatusAccepted)
}
