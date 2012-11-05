package notifier

import (
	"bytes"
	"fmt"
	"formatter"
	"gobus"
	"model"
	"service/args"
	"thirdpart"
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
	private, public, err := u.getWelcomeContent(arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}

	url := fmt.Sprintf("http://%s:%d", u.config.ExfeService.Addr, u.config.ExfeService.Port)
	client, err := gobus.NewClient(fmt.Sprintf("%s/%s", url, "Thirdpart"))
	if err != nil {
		return fmt.Errorf("can't create gobus client: %s", err)
	}

	a := args.SendArg{
		To:             &arg.To,
		PrivateMessage: private,
		PublicMessage:  public,
		Info: &thirdpart.InfoData{
			CrossID: 0,
			Type:    thirdpart.CrossInvitation,
		},
	}
	var ids string
	err = client.Do("Send", &a, &ids)
	if err != nil {
		return fmt.Errorf("send error: %s", err)
	}
	return nil
}

func (u User) Confirm(arg ConfirmArg) error {
	private, public, err := u.getConfirmContent(arg)
	if err != nil {
		return fmt.Errorf("can't get content: %s", err)
	}

	url := fmt.Sprintf("http://%s:%d", u.config.ExfeService.Addr, u.config.ExfeService.Port)
	client, err := gobus.NewClient(fmt.Sprintf("%s/%s", url, "Thirdpart"))
	if err != nil {
		return fmt.Errorf("can't create gobus client: %s", err)
	}

	a := args.SendArg{
		To:             &arg.To,
		PrivateMessage: private,
		PublicMessage:  public,
		Info: &thirdpart.InfoData{
			CrossID: 0,
			Type:    thirdpart.CrossInvitation,
		},
	}
	var ids string
	err = client.Do("Send", &a, &ids)
	if err != nil {
		return fmt.Errorf("send error: %s", err)
	}
	return nil
}

func (u User) getWelcomeContent(arg WelcomeArg) (string, string, error) {
	arg.Config = u.config
	messageType, err := thirdpart.MessageTypeFromProvider(arg.To.Provider)
	if err != nil {
		return "", "", err
	}

	templateName := fmt.Sprintf("user_welcome.%s", messageType)
	private := bytes.NewBuffer(nil)
	err = u.localTemplate.Execute(private, arg.To.Language, templateName, arg)
	if err != nil {
		return "", "", fmt.Errorf("private template(%s) failed: %s", templateName, err)
	}

	return private.String(), "", nil
}

func (u User) getConfirmContent(arg ConfirmArg) (string, string, error) {
	arg.Config = u.config
	messageType, err := thirdpart.MessageTypeFromProvider(arg.To.Provider)
	if err != nil {
		return "", "", err
	}

	templateName := fmt.Sprintf("user_confirm.%s", messageType)
	private := bytes.NewBuffer(nil)
	err = u.localTemplate.Execute(private, arg.To.Language, templateName, arg)
	if err != nil {
		return "", "", fmt.Errorf("private template(%s) failed: %s", templateName, err)
	}

	return private.String(), "", nil
}
