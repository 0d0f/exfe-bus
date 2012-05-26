package exfe_service

import (
	"bytes"
	"gomail"
	"strings"
	"log/syslog"
	"exfe/model"
	"fmt"
	"text/template"
	"gobus"
)

type WelcomeArg struct {
	To_identity exfe_model.Identity
	config *Config
}

type VerifyArg struct {
	To_identity exfe_model.Identity
	config *Config
}

type ResetPasswordArg struct {
	To_identity exfe_model.Identity
	config *Config
}

type ActiveArg struct {
	To_identity exfe_model.Identity
	config *Config
}

type WelcomeActiveArg struct {
	To_identity exfe_model.Identity
	config *Config
}

type Authentication struct {
	config *Config
	log *syslog.Writer
	client *gobus.Client
}

func NewAuthentication(config *Config) *Authentication {
	log, err := syslog.New(syslog.LOG_DEBUG, "exfe.auth")
	if err != nil {
		panic(err)
	}
	client := gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "email")
	return &Authentication{
		config: config,
		log: log,
		client: client,
	}
}

func executeTemplate(name string, to *exfe_model.Identity, data interface{}, client *gobus.Client) error {
	buf := bytes.NewBuffer(nil)
	tmpl, err := template.ParseFiles(fmt.Sprintf("./template/default/%s", name))
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}

	fmt.Println(buf.String())
	content := strings.SplitN(buf.String(), "\n", 2)

	mailarg := gomail.Mail{
		To: []gomail.MailUser{},
		From: gomail.MailUser{"x@exfe.com", "x@exfe.com"},
		Subject: content[0],
		Html: content[1],
	}

	client.Send("EmailSend", &mailarg, 5)

	return nil
}

func (s *Authentication) Welcome(arg *WelcomeArg, reply *int) error {
	arg.config = s.config

	err := executeTemplate("auth_welcome.html", &arg.To_identity, arg, s.client)
	if err != nil {
		s.log.Err(fmt.Sprintf("Execute template error: %s", err))
	}
	return nil
}

func (s *Authentication) Verify(arg *VerifyArg, reply *int) error {
	arg.config = s.config

	err := executeTemplate("auth_verify.html", &arg.To_identity, arg, s.client)
	if err != nil {
		s.log.Err(fmt.Sprintf("Execute template error: %s", err))
	}
	return nil
}

func (s *Authentication) ResetPassword(arg *ResetPasswordArg, reply *int) error {
	arg.config = s.config

	err := executeTemplate("auth_reset_password.html", &arg.To_identity, arg, s.client)
	if err != nil {
		s.log.Err(fmt.Sprintf("Execute template error: %s", err))
	}
	return nil
}

func (s *Authentication) Active(arg *ActiveArg, reply *int) error {
	arg.config = s.config

	err := executeTemplate("auth_active.html", &arg.To_identity, arg, s.client)
	if err != nil {
		s.log.Err(fmt.Sprintf("Execute template error: %s", err))
	}
	return nil
}

func (s *Authentication) WelcomeActive(arg *WelcomeActiveArg, reply *int) error {
	arg.config = s.config

	err := executeTemplate("auth_welcome_active.html", &arg.To_identity, arg, s.client)
	if err != nil {
		s.log.Err(fmt.Sprintf("Execute template error: %s", err))
	}
	return nil
}
