package main

import (
	"./pkg/twitter"
	"bytes"
	"config"
	"fmt"
	"gobus"
	"gosque"
	"log"
	"strconv"
	"strings"
	"text/template"
)

type ExfeTime struct {
	Time      string
	Data      string
	Datetime  string
	Time_type string
}

type TwitterSender struct {
	Title             string
	Description       string
	Begin_at          ExfeTime
	Place_line1       string
	Place_line2       string
	Cross_id          int64
	Cross_id_base62   string
	Invitation_id     string
	Token             string
	Identity_id       string
	Host_identity_id  int64
	Provider          string
	External_identity string
	Name              string
	Avatar_file_name  string
	Host_identity     struct {
		Name             string
		Avatar_file_name string
	}
	Rsvp_status int64
	By_identity struct {
		Id                string
		External_identity string
		Name              string
		Bio               string
		Avatar_file_name  string
		External_username string
		Provider          string
	}
	To_identity struct {
		Id                string
		External_identity string
		Name              string
		Bio               string
		Avatar_file_name  string
		External_username string
		Provider          string
	}
	Invitations []struct {
		Invitation_id     string
		State             int64
		By_identity_id    string
		Token             string
		Updated_at        string
		Identity_id       string
		Provider          string
		External_identity string
		Name              string
		Bio               string
		Avatar_file_name  string
		External_username string
		Identities        []struct {
			Identity_id       string
			Status            string
			Provider          string
			External_identity string
			Name              string
			Bio               string
			Avatar_file_name  string
			External_username string
		}
		User_id int64
	}
	config        *config.Configure
	getfriendship *gobus.Client
	getinfo       *gobus.Client
	sendtweet     *gobus.Client
	senddm        *gobus.Client
}

func (s *TwitterSender) CrossLink(withToken bool) string {
	link := fmt.Sprintf(" %s/!%s", s.config.String("site_url"), s.Cross_id_base62)
	if withToken {
		return fmt.Sprintf("%s?token=%s", link, s.Token)
	}
	return link
}

func (s *TwitterSender) isHost() bool {
	identity_id, _ := strconv.ParseInt(s.Identity_id, 10, 0)
	return identity_id == s.Host_identity_id
}

func ShortTweet(tweet string) string {
	const linkLength = 25
	if (len(tweet) + linkLength) > 140 {
		return tweet[0:(137-linkLength)] + "..."
	}
	return tweet
}

type TemplateData struct {
	ToUserName    string
	IsHost        bool
	Title         string
	Time          string
	Place1        string
	Place2        string
	SiteUrl       string
	CrossIdBase62 string
	Token         string
}

func CreateData(sender *TwitterSender) *TemplateData {
	return &TemplateData{
		ToUserName:    sender.To_identity.External_username,
		IsHost:        sender.isHost(),
		Title:         sender.Title,
		Time:          sender.Begin_at.Datetime,
		Place1:        sender.Place_line1,
		Place2:        sender.Place_line2,
		SiteUrl:       sender.config.String("site_url"),
		CrossIdBase62: sender.Cross_id_base62,
		Token:         sender.Token,
	}
}

func LoadTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles(fmt.Sprintf("./template/default/%s", name)))
}

func (s *TwitterSender) Do() {
	log.Printf("Get a job")

	if (s.External_identity != "") ||
		(strings.ToLower(s.External_identity) == strings.ToLower(fmt.Sprintf("@%s@twitter", s.To_identity.External_username))) {
		// update user info
		s.getinfo.Send(&twitter.UserInfo{
			ClientToken:  s.config.String("twitter.client_token"),
			ClientSecret: s.config.String("twitter.client_secret"),
			AccessToken:  s.config.String("twitter.access_token"),
			AccessSecret: s.config.String("twitter.access_secret"),
			ScreenName:   s.To_identity.External_username,
		})
	}

	// check friendship
	f := twitter.Friendship{
		ClientToken:  s.config.String("twitter.client_token"),
		ClientSecret: s.config.String("twitter.client_secret"),
		AccessToken:  s.config.String("twitter.access_token"),
		AccessSecret: s.config.String("twitter.access_secret"),
		UserA:        s.To_identity.External_username,
		UserB:        s.config.String("twitter.screen_name"),
	}
	var response string
	err := s.getfriendship.Do(f, &response)
	if err != nil {
		log.Printf("Twitter check friendship(%s/%s) fail: %s", f.UserA, f.UserB, err)
		return
	}
	isFriend := response == "true"

	data := CreateData(s)

	var tmpl *template.Template
	if isFriend {
		tmpl = LoadTemplate("twitter_sender_dm.tmpl")
	} else {
		tmpl = LoadTemplate("twitter_sender_tweet.tmpl")
	}
	buf := bytes.NewBuffer(nil)
	tmpl.Execute(buf, data)

	tweet := ShortTweet(strings.Trim(buf.String(), "\n \t")) + s.CrossLink(isFriend)

	if isFriend {
		s.sendDM(data.ToUserName, tweet)
	} else {
		s.sendTweet(tweet)
	}
}

func (s *TwitterSender) sendTweet(t string) {
	tweet := twitter.Tweet{
		ClientToken:  s.config.String("twitter.client_token"),
		ClientSecret: s.config.String("twitter.client_secret"),
		AccessToken:  s.config.String("twitter.access_token"),
		AccessSecret: s.config.String("twitter.access_secret"),
		Tweet:        t,
	}
	var response string
	err := s.sendtweet.Do(tweet, &response)
	if err != nil {
		log.Printf("Can't send tweet: %s", err)
	}
}

func (s *TwitterSender) sendDM(to_user string, t string) {
	dm := twitter.DirectMessage{
		ClientToken:  s.config.String("twitter.client_token"),
		ClientSecret: s.config.String("twitter.client_secret"),
		AccessToken:  s.config.String("twitter.access_token"),
		AccessSecret: s.config.String("twitter.access_secret"),
		Message:      t,
		ToUserName:   to_user,
	}
	var response string
	err := s.senddm.Do(dm, &response)
	if err != nil {
		log.Printf("Can't send tweet: %s", err)
	}
}

func TwitterSenderGenerator() interface{} {
	return &TwitterSender{}
}

func main() {
	log.SetPrefix("[TwitterSender]")
	log.Printf("Service start")
	config := config.LoadFile("twitter_sender.yaml")

	client := gosque.CreateQueue(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"resque:queue:twitter")

	sendtweet := gobus.CreateClient(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"twitter:tweet")

	senddm := gobus.CreateClient(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"twitter:directmessage")

	getinfo := gobus.CreateClient(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"twitter:userinfo")

	getfriendship := gobus.CreateClient(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"twitter:friendship")

	recv := client.IncomingJob("twitter_job", TwitterSenderGenerator, 5e9)
	for {
		select {
		case job := <-recv:
			twitterSender := job.(*TwitterSender)
			twitterSender.config = config
			twitterSender.sendtweet = sendtweet
			twitterSender.senddm = senddm
			twitterSender.getinfo = getinfo
			twitterSender.getfriendship = getfriendship
			go func() {
				twitterSender.Do()
			}()
		}
	}
	log.Printf("Service stop")
}
