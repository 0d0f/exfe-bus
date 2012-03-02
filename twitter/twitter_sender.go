package main

import (
	"fmt"
	"gosque"
	"config"
	"strconv"
	"strings"
	"gobus"
)

type ExfeTime struct {
	Time string
	Data string
	Datetime string
	Time_type string
}

type TwitterSender struct {
	Title string
	Description string
	Begin_at ExfeTime
	Place_line1 string
	Place_line2 string
	Cross_id int64
	Cross_id_base62 string
	Invitation_id string
	Token string
	Identity_id string
	Host_identity_id int64
	Provider string
	External_identity string
	Name string
	Avatar_file_name string
	Host_identity struct {
		Name string
		Avatar_file_name string
	}
	Rsvp_status int64
	By_identity struct {
		Id string
		External_identity string
		Name string
		Bio string
		Avatar_file_name string
		External_username string
		Provider string
	}
	To_identity struct {
		Id string
		External_identity string
		Name string
		Bio string
		Avatar_file_name string
		External_username string
		Provider string
	}
	Invitations []struct {
		Invitation_id string
		State int64
		By_identity_id string
		Token string
		Updated_at string
		Identity_id string
		Provider string
		External_identity string
		Name string
		Bio string
		Avatar_file_name string
		External_username string
		Identities []struct {
			Identity_id string
			Status string
			Provider string
			External_identity string
			Name string
			Bio string
			Avatar_file_name string
			External_username string
		}
		User_id int64
	}
	config *config.Configure
	getfriendship *gobus.Client
	getinfo *gobus.Client
	sendtweet *gobus.Client
	senddm *gobus.Client
}

func (s *TwitterSender) CrossLink(withToken bool) string {
	link := fmt.Sprintf(" %s/!%s", s.config.String("site_url"), s.Cross_id_base62)
	if withToken {
		return fmt.Sprintf("%s?token=%s", link, s.Token)
	}
	return link
}

func (s *TwitterSender) Time() string {
	if s.Begin_at.Datetime == "" {
		return ""
	}
	return fmt.Sprintf(" %s", s.Begin_at.Datetime)
}

func (s *TwitterSender) Place() string {
	if s.Place_line1 == "" {
		return ""
	}
	if s.Place_line2 == "" {
		return fmt.Sprintf(" at %s", s.Place_line1)
	}
	return fmt.Sprintf(" at %s, %s", s.Place_line1, s.Place_line2)
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

type Friendship struct {
	ClientToken string
	ClientSecret string
	AccessToken string
	AccessSecret string
	UserA string
	UserB string
}

func (s *TwitterSender) Do() {
	if (s.External_identity != "") ||
			(strings.ToLower(s.External_identity) == strings.ToLower(fmt.Sprintf("@%s@twitter", s.To_identity.External_username))) {
		// update user info
	}

	// check friendship
	response, _ := s.getfriendship.Do(&Friendship{
		ClientToken: s.config.String("twitter.client_token"),
		ClientSecret: s.config.String("twitter.client_secret"),
		AccessToken: s.config.String("twitter.access_token"),
		AccessSecret: s.config.String("twitter.access_secret"),
		UserA: s.To_identity.External_username,
		UserB: s.config.String("twitter.screen_name"),
	})
	isFriend := response.(*TwitterResponse).Result == "true"

	// build tweet
	var tweet string
	if s.isHost() {
		tweet = "You're successfully gathering this X"
	} else {
		tweet = "Invitation"
	}

	tweet = fmt.Sprintf("%s: %s.%s%s", tweet, s.Title, s.Time(), s.Place())

	if !isFriend {
		tweet = fmt.Sprintf("@%s %s", s.To_identity.External_username, tweet)
	}

	tweet = ShortTweet(tweet) + s.CrossLink(isFriend)

	if isFriend {
		fmt.Println("in dm")
		s.sendDM(s.To_identity.External_username, tweet)
	} else {
		fmt.Println("in tweet")
		s.sendTweet(tweet)
	}
}

type Tweet struct {
	ClientToken string
	ClientSecret string
	AccessToken string
	AccessSecret string
	Tweet string
}

func (s *TwitterSender) sendTweet(t string) {
	tweet := &Tweet{
		ClientToken: s.config.String("twitter.client_token"),
		ClientSecret: s.config.String("twitter.client_secret"),
		AccessToken: s.config.String("twitter.access_token"),
		AccessSecret: s.config.String("twitter.access_secret"),
		Tweet: t,
	}
	s.sendtweet.Do(tweet)
}

type DM struct {
	ClientToken string
	ClientSecret string
	AccessToken string
	AccessSecret string
	Message string
	ToUserName string
	ToUserId string
}

func (s *TwitterSender) sendDM(to_user string, t string) {
	dm := &DM{
		ClientToken: s.config.String("twitter.client_token"),
		ClientSecret: s.config.String("twitter.client_secret"),
		AccessToken: s.config.String("twitter.access_token"),
		AccessSecret: s.config.String("twitter.access_secret"),
		Message: t,
		ToUserName: to_user,
	}
	s.senddm.Do(dm)
}

func TwitterSenderGenerator() interface{} {
	return &TwitterSender{}
}

type TwitterResponse struct {
	Error string
	Result string
}

func TwitterResponseGenerator() interface{} {
	return &TwitterResponse{}
}

func main() {
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
		"gobus:queue:twitter:tweet",
		TwitterResponseGenerator)

	senddm := gobus.CreateClient(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"gobus:queue:twitter:directmessage",
		TwitterResponseGenerator)

	getinfo := gobus.CreateClient(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"gobus:queue:twitter:usershow",
		TwitterResponseGenerator)

	getfriendship := gobus.CreateClient(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		"gobus:queue:twitter:friendship",
		TwitterResponseGenerator)

	recv := client.IncomingJob(TwitterSenderGenerator, 5e9)
	for {
		select {
		case job := <-recv:
			twitterSender := job.(*TwitterSender)
			twitterSender.config = config
			twitterSender.sendtweet = sendtweet
			twitterSender.senddm = senddm
			twitterSender.getinfo = getinfo
			twitterSender.getfriendship = getfriendship
			twitterSender.Do()
		}
	}
}
