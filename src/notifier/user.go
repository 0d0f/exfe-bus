package notifier

import (
	"broker"
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
}

func NewUser(localTemplate *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) *User {
	return &User{
		localTemplate: localTemplate,
		config:        config,
		platform:      platform,
	}
}

func (u User) HandleWelcome(arg model.UserWelcome) {
	err := arg.Parse(u.config)
	if err != nil {
		u.Error(http.StatusBadRequest, err)
		return
	}

	to := arg.To
	text, err := GenerateContent(u.localTemplate, "user_welcome", to.Provider, to.Language, arg)
	if err != nil {
		u.Error(http.StatusInternalServerError, err)
		return
	}
	_, _, _, err = u.platform.Send(to, text)
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

	to := arg.To
	text, err := GenerateContent(u.localTemplate, "user_verify", to.Provider, to.Language, arg)
	if err != nil {
		u.Error(http.StatusInternalServerError, err)
		return
	}
	_, _, _, err = u.platform.Send(to, text)
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

	to := arg.To
	text, err := GenerateContent(u.localTemplate, "user_resetpass", to.Provider, to.Language, arg)
	if err != nil {
		u.Error(http.StatusInternalServerError, err)
		return
	}
	_, _, _, err = u.platform.Send(to, text)
	if err != nil {
		u.Error(http.StatusInternalServerError, err)
		return
	}
}
