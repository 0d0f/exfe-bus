package twitter_job

import (
	"net/http"
	"twitter/service"
	"bytes"
	"fmt"
	"gobus"
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

	Config        *Config
	Getfriendship *gobus.Client
	Getinfo       *gobus.Client
	Sendtweet     *gobus.Client
	Senddm        *gobus.Client
}

func (s *TwitterSender) updateUserInfo(id uint64, i *twitter_service.TwitterUserInfo) {
	url := fmt.Sprintf("%s/identity/update", s.Config.Site_url)
	_, err := http.PostForm(url, i.MakeUrlValues(id))
	if err != nil {
		log.Printf("[Error]Update identity info fail: %s", err)
	} else {
		log.Printf("[Info]Update identity info success")
	}
}

func (s *TwitterSender) CrossLink(withToken bool) string {
	link := fmt.Sprintf(" %s/!%s", s.Config.Site_url, s.Cross_id_base62)
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
		SiteUrl:       sender.Config.Site_url,
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
		var reply twitter_service.TwitterUserInfo
		err := s.Getinfo.Do(&twitter_service.UsersShowArg{
			ClientToken:  s.Config.Twitter.Client_token,
			ClientSecret: s.Config.Twitter.Client_secret,
			AccessToken:  s.Config.Twitter.Access_token,
			AccessSecret: s.Config.Twitter.Access_secret,
			ScreenName:   s.To_identity.External_username,
		}, &reply)
		if err == nil {
			id, _ := strconv.ParseUint(s.Identity_id, 10, 64)
			go s.updateUserInfo(id, &reply)
		}
	}

	// check friendship
	f := &twitter_service.FriendshipsExistsArg{
		ClientToken:  s.Config.Twitter.Client_token,
		ClientSecret: s.Config.Twitter.Client_secret,
		AccessToken:  s.Config.Twitter.Access_token,
		AccessSecret: s.Config.Twitter.Access_secret,
		UserA:        s.To_identity.External_username,
		UserB:        s.Config.Twitter.Screen_name,
	}
	var isFriend bool
	err := s.Getfriendship.Do(f, &isFriend)
	if err != nil {
		log.Printf("Twitter check friendship(%s/%s) fail: %s", f.UserA, f.UserB, err)
		return
	}

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
	tweet := &twitter_service.StatusesUpdateArg{
		ClientToken:  s.Config.Twitter.Client_token,
		ClientSecret: s.Config.Twitter.Client_secret,
		AccessToken:  s.Config.Twitter.Access_token,
		AccessSecret: s.Config.Twitter.Access_secret,
		Tweet:        t,
	}
	var response twitter_service.StatusesUpdateReply
	err := s.Sendtweet.Do(tweet, &response)
	if err != nil {
		log.Printf("Can't send tweet: %s", err)
		return
	}
}

func (s *TwitterSender) sendDM(to_user string, t string) {
	dm := &twitter_service.DirectMessagesNewArg{
		ClientToken:  s.Config.Twitter.Client_token,
		ClientSecret: s.Config.Twitter.Client_secret,
		AccessToken:  s.Config.Twitter.Access_token,
		AccessSecret: s.Config.Twitter.Access_secret,
		Message:      t,
		ToUserName:   to_user,
	}
	var response twitter_service.DirectMessagesNewReply
	err := s.Senddm.Do(dm, &response)
	if err != nil {
		log.Printf("Can't send tweet: %s", err)
		return
	}

	i, _ := strconv.ParseUint(s.To_identity.Id, 10, 64)
	go s.updateUserInfo(i, &response.Recipient)
}

func TwitterSenderGenerator() interface{} {
	return &TwitterSender{}
}
