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
	text, err := GenerateContent(u.localTemplate, "v3_user_welcome", to.Provider, to.Language, arg)
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
	text, err := GenerateContent(u.localTemplate, "v3_user_verify", to.Provider, to.Language, arg)
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
	text, err := GenerateContent(u.localTemplate, "v3_user_resetpass", to.Provider, to.Language, arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	_, err = u.platform.Send(to, text)
	return err
}

func (u User) Welcome(arg model.UserWelcome) error {
	err := arg.Parse(u.config)
	if err != nil {
		return err
	}

	content, err := GetContent(u.localTemplate, "user_welcome", arg.To, arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	return u.send(content, arg.To)
}

func (u User) Verify(arg model.UserVerify) error {
	err := arg.Parse(u.config)
	if err != nil {
		return err
	}

	content, err := GetContent(u.localTemplate, "user_verify", arg.To, arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	return u.send(content, arg.To)
}

func (u User) ResetPassword(arg model.UserVerify) error {
	err := arg.Parse(u.config)
	if err != nil {
		return err
	}

	content, err := GetContent(u.localTemplate, "user_resetpass", arg.To, arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	return u.send(content, arg.To)
}

func (u User) send(content string, to model.Recipient) error {
	_, err := u.platform.Send(to, content)

	if err != nil {
		return fmt.Errorf("send error: %s", err)
	}
	return nil
}
