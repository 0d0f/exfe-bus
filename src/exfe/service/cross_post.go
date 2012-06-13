package exfe_service

import (
	"gobus"
	"fmt"
	"twitter/service"
	"apn/service"
	"c2dm/service"
	"net/http"
	"net/url"
	"log"
	"os"
	"io/ioutil"
)

type CrossPost struct {
	twitter *gobus.Client
	ios *gobus.Client
	android *gobus.Client
	log *log.Logger
	config *Config
}

func NewCrossPost(config *Config) *CrossPost {
	return &CrossPost{
		twitter: gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "twitter"),
		ios: gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "iOSAPN"),
		android: gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "Android"),
		log: log.New(os.Stderr, "exfe.cross.post", log.LstdFlags),
		config: config,
	}
}

func (s *CrossPost) SendPost(arg *OneIdentityUpdateArg) {
	by := arg.By_identity.Name
	switch arg.To_identity.Provider {
	case "twitter":
		if arg.By_identity.Provider == "twitter" {
			by = fmt.Sprintf("@%s", arg.By_identity.External_username)
		}
		s.SendTwitter(arg, fmt.Sprintf("%s %s", by, arg.Post.Content))
	case "iOSAPN":
		s.SendApn(arg, fmt.Sprintf("%s %s", by, arg.Post.Content))
	case "Android":
		s.SendAndroid(arg, fmt.Sprintf("%s %s", by, arg.Post.Content))
	}
}

func (s *CrossPost) SendTwitter(arg *OneIdentityUpdateArg, msg string) {
	params := make(url.Values)
	params.Add("data", fmt.Sprintf("%s", arg.Cross.Id))
	resp, err := http.PostForm(fmt.Sprintf("%s/iom/%s"), params)
	if err != nil {
		s.log.Printf("access iom server error: %s", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Printf("read response from iom server error: %s", err)
		return
	}
	if resp.StatusCode != 200 {
		s.log.Printf("iom server return error(%s): %s", resp.Status, string(body))
		return
	}
	urls := []string{fmt.Sprintf("#%s", string(body))}

	isFriend := false
	f := &twitter_service.FriendshipsExistsArg{
		ClientToken:  s.config.Twitter.Client_token,
		ClientSecret: s.config.Twitter.Client_secret,
		AccessToken:  s.config.Twitter.Access_token,
		AccessSecret: s.config.Twitter.Access_secret,
		UserA:        arg.To_identity.External_username,
		UserB:        s.config.Twitter.Screen_name,
	}
	err = s.twitter.Do("GetFriendship", f, &isFriend, 10)
	if err != nil {
		s.log.Printf("Can't require identity %d friendship: %s", arg.To_identity.Id, err)
	}

	if isFriend {
		dm := &twitter_service.DirectMessagesNewArg{
			ClientToken:  s.config.Twitter.Client_token,
			ClientSecret: s.config.Twitter.Client_secret,
			AccessToken:  s.config.Twitter.Access_token,
			AccessSecret: s.config.Twitter.Access_secret,
			Message:      msg,
			Urls:         urls,
			ToUserName:   &arg.To_identity.External_username,
			IdentityId:   &arg.To_identity.Id,
		}
		s.twitter.Send("SendDM", dm, 5)
	} else {
		tweet := &twitter_service.StatusesUpdateArg{
			ClientToken:  s.config.Twitter.Client_token,
			ClientSecret: s.config.Twitter.Client_secret,
			AccessToken:  s.config.Twitter.Access_token,
			AccessSecret: s.config.Twitter.Access_secret,
			Tweet:        fmt.Sprintf("@%s %s", arg.To_identity.External_username, msg),
			Urls:         urls,
		}
		s.twitter.Send("SendTweet", tweet, 5)
	}
}

func (s *CrossPost) SendApn(arg *OneIdentityUpdateArg, msg string) {
	a := apn_service.ApnSendArg{
		DeviceToken: arg.To_identity.External_id,
		Alert: msg,
		Badge: 0,
		Sound: "default",
		Cid: arg.Cross.Id,
		T: "c",
	}
	s.ios.Send("ApnSend", &a, 5)
}

func (s *CrossPost) SendAndroid(arg *OneIdentityUpdateArg, msg string) {
	a := c2dm_service.C2DMSendArg{
		DeviceID: arg.To_identity.External_id,
		Message: msg,
		Cid: arg.Cross.Id,
		T: "c",
		Badge: 0,
		Sound: "default",
	}
	s.android.Send("C2DMSend", &a, 5)
}
