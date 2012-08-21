package exfe_service

import (
	"bytes"
	"email/service"
	"encoding/json"
	"exfe/model"
	"fmt"
	"gobus"
	"io/ioutil"
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

	content := strings.SplitN(buf.String(), "//////////////////////////////////\n\n", 3)

	mailarg := &email_service.MailArg{
		To:      []*mail.Address{&mail.Address{to.Name, to.External_id}},
		From:    &mail.Address{"EXFE ·X·", "x@exfe.com"},
		Subject: content[0],
		Text:    content[1],
		Html:    content[2],
	}

	client.Send("EmailSend", &mailarg, 5)

	return nil
}

func (s *User) Welcome(arg *UserArg, reply *int) error {
	s.log.Printf("welcome to %s", arg.To_identity.ExternalId())
	arg.Config = s.config

	err := executeTemplate("user_welcome.html", &arg.To_identity, arg, s.email)
	if err != nil {
		s.log.Printf("Execute template error: %s", err)
	}
	return nil
}

func (s *User) Verify(arg *UserArg, reply *int) error {
	s.log.Printf("verify to %s", arg.To_identity.ExternalId())
	arg.Config = s.config

	template := fmt.Sprintf("user_%s.html", strings.ToLower(arg.Action))
	err := executeTemplate(template, &arg.To_identity, arg, s.email)
	if err != nil {
		s.log.Printf("Execute template error: %s", err)
	}
	return nil
}

type GetFriendsArg struct {
	ClientToken  string `json:"client_token"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
	AccessSecret string `json:"access_secret"`
	UserID       uint64 `json:"user_id"`
	ExternalID   string `json:"external_id"`
	Provider     string `json:"provider"`
}

func (s *User) GetFriends(arg *GetFriendsArg, reply *int) error {
	s.log.Printf("get friend from %s@%s(%d)", arg.ExternalID, arg.Provider, arg.UserID)

	switch arg.Provider {
	case "twitter":
		s.getTwitterFriends(arg)
	case "facebook":
		s.getFacebookFriends(arg)
	default:
		s.log.Printf("don't know how to get provider %s friend", arg.Provider)
	}
	return nil
}

func (s *User) getTwitterFriends(arg *GetFriendsArg) {
	friendsArg := &twitter_service.FriendsArg{
		ClientToken:  arg.ClientToken,
		ClientSecret: arg.ClientSecret,
		AccessToken:  arg.AccessToken,
		AccessSecret: arg.AccessSecret,
		UserId:       arg.ExternalID,
	}
	var friendReply twitter_service.FriendsReply
	err := s.twitter.Do("Friends", friendsArg, &friendReply, 5)
	if err != nil {
		s.log.Printf("Get friend error: %s", err)
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
		var identities []*exfe_model.Identity
		var reply []twitter_service.UserInfo
		err := s.twitter.Do("Lookup", &lookupArg, &reply, 5)
		if err != nil {
			s.log.Printf("Lookup users error: %s", err)
			continue
		}

		for _, i := range reply {
			identity := &exfe_model.Identity{
				Name:              *i.Name,
				Bio:               *i.Description,
				Avatar_filename:   *i.Profile_image_url,
				External_id:       fmt.Sprintf("%d", i.Id),
				External_username: *i.Screen_name,
				Provider:          "twitter",
			}
			identities = append(identities, identity)
		}

		go s.UpdateIdentities(arg.UserID, identities)
	}
}

type FacebookUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type FacebookPaging struct {
	Next string `json:"next"`
}

type FacebookFriendsReply struct {
	Data   []FacebookUser `json:"data"`
	Paging FacebookPaging `json:"paging"`
}

func (s *User) getFacebookFriends(arg *GetFriendsArg) {
	url := fmt.Sprintf("https://graph.facebook.com/%s/friends?access_token=%s", arg.ExternalID, arg.AccessToken)
	for {
		resp, err := http.Get(url)
		if err != nil {
			s.log.Printf("facebook get friends from %s error: %s", url, err)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.log.Printf("facebook get body from %s fail: %s", url, err)
		}
		if resp.StatusCode != 200 {
			s.log.Printf("facebook get friends from %s fail: (%s) %s", url, resp.Status, string(body))
			return
		}
		var friends FacebookFriendsReply
		err = json.Unmarshal(body, &friends)
		if err != nil {
			s.log.Printf("facebook get friends json error: %s", err)
			return
		}
		if len(friends.Data) == 0 {
			break
		}
		var identities []*exfe_model.Identity
		for _, f := range friends.Data {
			identity := &exfe_model.Identity{
				Name:              f.Name,
				Avatar_filename:   fmt.Sprintf("http://graph.facebook.com/%s/picture", f.Username),
				External_id:       f.Id,
				External_username: f.Username,
				Provider:          "facebook",
			}
			identities = append(identities, identity)
		}
		s.UpdateIdentities(arg.UserID, identities)
		url = friends.Paging.Next
	}
}

type UpdateIdentitiesArg struct {
	UserId     uint64                 `json:"user_id"`
	Identities []*exfe_model.Identity `json:"identities"`
}

func (s *User) UpdateIdentities(userId uint64, identities []*exfe_model.Identity) {
	arg := &UpdateIdentitiesArg{
		UserId:     userId,
		Identities: identities,
	}

	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	err := e.Encode(arg)
	if err != nil {
		s.log.Printf("encoding arg error: %s", err)
		return
	}
	url := fmt.Sprintf("%s/v2/AddFriends", s.config.Site_api)
	s.log.Printf("send to url: %s, post: %s", url, buf.String())
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		s.log.Printf("Send twitter friend to server error: %s", err)
		return
	}
	if resp.StatusCode != 200 {
		s.log.Printf("Send twitter friend to server error: %s", resp.Status)
	}
}
