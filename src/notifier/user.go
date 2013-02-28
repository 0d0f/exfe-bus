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
