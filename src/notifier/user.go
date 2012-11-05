package notifier

import (
	"fmt"
	"formatter"
	"gobus"
	"model"
	"service/args"
)

type WelcomeArg struct {
	ArgBase
	NeedVerify bool `json:"need_verify"`
}

type ConfirmArg struct {
	ArgBase
	By *model.Identity `json:"by"`
}

func (a ConfirmArg) NeedShowBy() bool {
	return !a.To.SameUser(a.By)
}

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

func (u User) Welcome(arg WelcomeArg) error {
	err := arg.Parse(u.config)
	if err != nil {
		return err
	}

	content, err := GetContent(u.localTemplate, "user_welcome", arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	return u.send(content, arg.ArgBase)
}

func (u User) Confirm(arg ConfirmArg) error {
	err := arg.Parse(u.config)
	if err != nil {
		return err
	}

	content, err := GetContent(u.localTemplate, "user_confirm", arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	return u.send(content, arg.ArgBase)
}

func (u User) ResetPassword(arg ArgBase) error {
	err := arg.Parse(u.config)
	if err != nil {
		return err
	}

	content, err := GetContent(u.localTemplate, "user_resetpass", arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}
	return u.send(content, arg)
}

func (u User) send(content string, arg ArgBase) error {
	url := fmt.Sprintf("http://%s:%d", u.config.ExfeService.Addr, u.config.ExfeService.Port)
	client, err := gobus.NewClient(fmt.Sprintf("%s/%s", url, "Thirdpart"))
	if err != nil {
		return fmt.Errorf("can't create gobus client: %s", err)
	}

	a := args.SendArg{
		To:             &arg.To,
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
