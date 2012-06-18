package exfe_service

import (
	"bytes"
	"gomail"
	"strings"
	"log"
	"exfe/model"
	"fmt"
	"text/template"
	"gobus"
	"os"
)

type UserArg struct {
	To_identity exfe_model.Identity
	User_name string
	Token string
	Action string

	Config *Config
}

func (a *UserArg) Shorten(s string) string {
	if len(s) < 10 {
		return s
	}
	return fmt.Sprintf("%s…%s", s[0:3], s[len(s)-5:len(s)])
}

func (a *UserArg) NeedVerify() bool {
	return a.Token != ""
}

type User struct {
	config *Config
	log *log.Logger
	client *gobus.Client
}

func NewUser(config *Config) *User {
	log := log.New(os.Stderr, "exfe.auth", log.LstdFlags)
	client := gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "email")
	return &User{
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

func (s *User) Welcome(arg *UserArg, reply *int) error {
	arg.Config = s.config

	err := executeTemplate("auth_welcome.html", &arg.To_identity, arg, s.client)
	if err != nil {
		log.Printf("Execute template error: %s", err)
	}
	return nil
}

func (s *User) Verify(arg *UserArg, reply *int) error {
	arg.Config = s.config

	template := fmt.Sprintf("user_%s", strings.ToLower(arg.Action))
	err := executeTemplate(template, &arg.To_identity, arg, s.client)
	if err != nil {
		log.Printf("Execute template error: %s", err)
	}
	return nil
}
