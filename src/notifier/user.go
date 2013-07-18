package notifier

import (
	"broker"
	"fmt"
	"formatter"
	"github.com/googollee/go-rest"
	"model"
	"net/http"
)

type User struct {
	rest.Service `prefix:"/v3/notifier/user"`

	Welcome rest.Processor `path:"/welcome" method:"POST"`
	Verify  rest.Processor `path:"/verify" method:"POST"`
	Reset   rest.Processor `path:"/rest" method:"POST"`

	localTemplate *formatter.LocalTemplate
	config        *model.Config
	platform      *broker.Platform
	domain        string
}

func NewUser(localTemplate *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) *User {
	return &User{
		localTemplate: localTemplate,
		config:        config,
		platform:      platform,
		domain:        fmt.Sprintf("http://%s:%d", config.ExfeService.Addr, config.ExfeService.Port),
	}
}

func (u User) HandleWelcome(arg model.UserWelcome) {
	err := arg.Parse(u.config)
	if err != nil {
		u.Error(http.StatusBadRequest, err)
		return
	}

	to := arg.To.PopRecipient()
	err = SendAndSave(u.localTemplate, u.platform, to, arg, "user_welcome", u.domain+"/v3/notifier/user/welcome")
	if err != nil {
		u.Error(http.StatusInternalServerError, err)
		return
	}
}

func (u User) HandleVerify(arg model.UserVerify) {
	err := arg.Parse(u.config)
	if err != nil {
		u.Error(http.StatusBadRequest, err)
		return
	}

	to := arg.To.PopRecipient()
	err = SendAndSave(u.localTemplate, u.platform, to, arg, "user_verify", u.domain+"/v3/notifier/user/verify")
	if err != nil {
		u.Error(http.StatusInternalServerError, err)
		return
	}
}

func (u User) HandleReset(arg model.UserVerify) {
	err := arg.Parse(u.config)
	if err != nil {
		u.Error(http.StatusBadRequest, err)
		return
	}

	to := arg.To.PopRecipient()
	err = SendAndSave(u.localTemplate, u.platform, to, arg, "user_resetpass", u.domain+"/v3/notifier/user/reset")
	if err != nil {
		u.Error(http.StatusInternalServerError, err)
		return
	}
}
