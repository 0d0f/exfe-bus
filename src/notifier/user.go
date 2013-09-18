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

	welcome rest.SimpleNode `route:"/welcome" method:"POST"`
	verify  rest.SimpleNode `route:"/verify" method:"POST"`
	reset   rest.SimpleNode `route:"/reset" method:"POST"`

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

func (u User) Welcome(ctx rest.Context, arg model.UserWelcome) {
	err := arg.Parse(u.config)
	if err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}

	go SendAndSave(u.localTemplate, u.platform, &arg.To, arg, "user_welcome", u.domain+"/v3/notifier/user/welcome", &arg)
	ctx.Return(http.StatusAccepted)
}

func (u User) Verify(ctx rest.Context, arg model.UserVerify) {
	err := arg.Parse(u.config)
	if err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}

	go SendAndSave(u.localTemplate, u.platform, &arg.To, arg, "user_verify", u.domain+"/v3/notifier/user/verify", &arg)
	ctx.Return(http.StatusAccepted)
}

func (u User) Reset(ctx rest.Context, arg model.UserVerify) {
	err := arg.Parse(u.config)
	if err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}

	go SendAndSave(u.localTemplate, u.platform, &arg.To, arg, "user_resetpass", u.domain+"/v3/notifier/user/reset", &arg)
	ctx.Return(http.StatusAccepted)
}
