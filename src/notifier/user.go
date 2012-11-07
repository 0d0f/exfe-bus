package notifier

import (
	"fmt"
	"formatter"
	"gobus"
	"model"
)

type User struct {
	localTemplate *formatter.LocalTemplate
	config        *model.Config
}

func NewUser(localTemplate *formatter.LocalTemplate, config *model.Config) *User {
	return &User{
		localTemplate: localTemplate,
		config:        config,
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

func (u User) ResetPassword(arg model.ThirdpartTo) error {
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
	url := fmt.Sprintf("http://%s:%d", u.config.ExfeService.Addr, u.config.ExfeService.Port)
	client, err := gobus.NewClient(fmt.Sprintf("%s/%s", url, "Thirdpart"))
	if err != nil {
		return fmt.Errorf("can't create gobus client: %s", err)
	}

	a := model.ThirdpartSend{
		To:             to,
		PrivateMessage: content,
		PublicMessage:  "",
	}
	var ids string
	err = client.Do("Send", &a, &ids)
	if err != nil {
		return fmt.Errorf("send error: %s", err)
	}
	return nil
}
