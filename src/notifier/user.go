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

	go SendAndSave(u.localTemplate, u.platform, &arg.To, arg, "user_welcome", u.domain+"/v3/notifier/user/welcome", &arg)
	u.WriteHeader(http.StatusAccepted)
}

func (u User) HandleVerify(arg model.UserVerify) {
	err := arg.Parse(u.config)
	if err != nil {
		u.Error(http.StatusBadRequest, err)
		return
	}

	go SendAndSave(u.localTemplate, u.platform, &arg.To, arg, "user_verify", u.domain+"/v3/notifier/user/verify", &arg)
	u.WriteHeader(http.StatusAccepted)
}

func (u User) HandleReset(arg model.UserVerify) {
	err := arg.Parse(u.config)
	if err != nil {
		u.Error(http.StatusBadRequest, err)
		return
	}

	go SendAndSave(u.localTemplate, u.platform, &arg.To, arg, "user_resetpass", u.domain+"/v3/notifier/user/reset", &arg)
	u.WriteHeader(http.StatusAccepted)
}
