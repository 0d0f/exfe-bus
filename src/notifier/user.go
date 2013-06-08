package notifier

import (
	"broker"
	"fmt"
	"formatter"
	"model"
)

type User struct {
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

func (u User) V3Welcome(arg model.UserWelcome) error {
	err := arg.Parse(u.config)
	if err != nil {
		return err
	}

	to := arg.To
	to = to.Tunnel()
	text, err := GenerateContent(u.localTemplate, "user_welcome", to.Provider, to.Language, arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	_, err = u.platform.Send(to, text)
	return err
}

func (u User) V3Verify(arg model.UserVerify) error {
	err := arg.Parse(u.config)
	if err != nil {
		return err
	}

	to := arg.To
	to = to.Tunnel()
	text, err := GenerateContent(u.localTemplate, "user_verify", to.Provider, to.Language, arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	_, err = u.platform.Send(to, text)
	return err
}

func (u User) V3ResetPassword(arg model.UserVerify) error {
	err := arg.Parse(u.config)
	if err != nil {
		return err
	}

	to := arg.To
	to = to.Tunnel()
	text, err := GenerateContent(u.localTemplate, "user_resetpass", to.Provider, to.Language, arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	_, err = u.platform.Send(to, text)
	return err
}
