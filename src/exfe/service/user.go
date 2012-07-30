package exfe_service

import (
	"bytes"
	"email/service"
	"encoding/json"
	"exfe/model"
	"fmt"
	"gobus"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"text/template"
	"twitter/service"
)

type UserArg struct {
	To_identity exfe_model.Identity
	User_name   string
	Token       string
	Action      string

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
	config  *Config
	log     *log.Logger
	email   *gobus.Client
	twitter *gobus.Client
}

func NewUser(config *Config) *User {
	log := log.New(os.Stderr, "exfe.auth", log.LstdFlags)
	email := gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "email")
	twitter := gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "twitter")
	return &User{
		config:  config,
		log:     log,
		email:   email,
		twitter: twitter,
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

	mailarg := &email_service.MailArg{
		To:      []*mail.Address{&mail.Address{to.Name, to.External_id}},
		From:    &mail.Address{"x@exfe.com", "x@exfe.com"},
		Subject: content[0],
		Html:    content[1],
	}

	client.Send("EmailSend", &mailarg, 5)

	return nil
}

func (s *User) Welcome(arg *UserArg, reply *int) error {
	arg.Config = s.config

	err := executeTemplate("user_welcome.html", &arg.To_identity, arg, s.email)
	if err != nil {
		log.Printf("Execute template error: %s", err)
	}
	return nil
}

func (s *User) Verify(arg *UserArg, reply *int) error {
	arg.Config = s.config

	template := fmt.Sprintf("user_%s", strings.ToLower(arg.Action))
	err := executeTemplate(template, &arg.To_identity, arg, s.email)
	if err != nil {
		log.Printf("Execute template error: %s", err)
	}
	return nil
}

func (s *User) TwitterFriends(arg *twitter_service.FriendsArg, reply *int) error {
	var friendReply twitter_service.FriendsReply
	err := s.twitter.Do("Friends", arg, &friendReply, 5)
	if err != nil {
		log.Printf("Get friend error: %s", err)
		return err
	}
	users := friendReply.Ids
	lookupArg := twitter_service.UsersLookupArg{
		ClientToken:  arg.ClientToken,
		ClientSecret: arg.ClientSecret,
		AccessToken:  arg.AccessToken,
		AccessSecret: arg.AccessSecret,
	}
	for len(users) > 0 {
		if len(users) > 100 {
			lookupArg.UserId = users[:100]
			users = users[100:]
		} else {
			lookupArg.UserId = users
			users = nil
		}
		var reply []twitter_service.UserInfo
		err := s.twitter.Do("Lookup", &lookupArg, &reply, 5)
		if err != nil {
			log.Printf("Lookup users error: %s", err)
			continue
		}

		infos := make([]map[string]string, 0, 0)
		for _, i := range reply {
			info := make(map[string]string)
			info["provider"] = "twitter"
			info["external_id"] = fmt.Sprintf("%d", i.Id)
			if i.Name != nil {
				info["name"] = *i.Name
			}
			if i.Description != nil {
				info["bio"] = *i.Description
			}
			if i.Profile_image_url != nil {
				info["avatar_filename"] = *i.Profile_image_url
			}
			if i.Screen_name != nil {
				info["external_username"] = *i.Screen_name
			}
			infos = append(infos, info)
		}

		buf := bytes.NewBuffer(nil)
		e := json.NewEncoder(buf)
		e.Encode(infos)
		resp, err := http.Post(fmt.Sprintf("%s/v2/Friends", s.config.Site_api), "application/json", buf)
		if err != nil {
			log.Printf("Send twitter friend to server error: %s", err)
			continue
		}
		if resp.StatusCode != 200 {
			log.Printf("Send twitter friend to server error: %s", resp.Status)
			continue
		}
	}
	return nil
}
