package twitter_job

import (
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

type TwitterJobArg struct {
	Title             string
	Description       string
	Begin_at          ExfeTime
	Time_type         string
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
	To_identity_time_zone *string
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
}

type Twitter_job struct {
	Config        *Config
	Getfriendship *gobus.Client
	Getinfo       *gobus.Client
	Sendtweet     *gobus.Client
	Senddm        *gobus.Client
}

func (s *TwitterJobArg) CrossLink(siteUrl string, withToken bool) string {
	link := fmt.Sprintf(" %s/!%s", siteUrl, s.Cross_id_base62)
	if withToken {
		return fmt.Sprintf("%s?token=%s", link, s.Token)
	}
	return link
}

func (s *TwitterJobArg) isHost() bool {
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

func (s *TwitterJobArg) CreateData(siteUrl string) *TemplateData {
	return &TemplateData{
		ToUserName:    s.To_identity.External_username,
		IsHost:        s.isHost(),
		Title:         s.Title,
		Time:          s.Begin_at.Datetime,
		Place1:        s.Place_line1,
		Place2:        s.Place_line2,
		SiteUrl:       siteUrl,
		CrossIdBase62: s.Cross_id_base62,
		Token:         s.Token,
	}
}

func LoadTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles(fmt.Sprintf("./template/default/%s", name)))
}

func (s *Twitter_job) Perform(arg *TwitterJobArg) {
	log.Printf("[TwitterJob]Get a job")

	toIdentityId, _ := strconv.ParseUint(arg.To_identity.Id, 10, 64)

	if arg.To_identity.External_identity == "" {
		// get to_identity info
		s.Getinfo.Send(&twitter_service.UsersShowArg{
			ClientToken:  s.Config.Twitter.Client_token,
			ClientSecret: s.Config.Twitter.Client_secret,
			AccessToken:  s.Config.Twitter.Access_token,
			AccessSecret: s.Config.Twitter.Access_secret,
			ScreenName:   &arg.To_identity.External_username,
			IdentityId:   &toIdentityId,
		})
	}

	// check friendship
	f := &twitter_service.FriendshipsExistsArg{
		ClientToken:  s.Config.Twitter.Client_token,
		ClientSecret: s.Config.Twitter.Client_secret,
		AccessToken:  s.Config.Twitter.Access_token,
		AccessSecret: s.Config.Twitter.Access_secret,
		UserA:        arg.To_identity.External_username,
		UserB:        s.Config.Twitter.Screen_name,
	}
	var isFriend bool
	err := s.Getfriendship.Do(f, &isFriend)
	if err != nil {
		isFriend = false
	}

	data := arg.CreateData(s.Config.Site_url)

	var tmpl *template.Template
	if isFriend {
		tmpl = LoadTemplate("twitter_sender_dm.tmpl")
	} else {
		tmpl = LoadTemplate("twitter_sender_tweet.tmpl")
	}
	buf := bytes.NewBuffer(nil)
	tmpl.Execute(buf, data)

	tweet := ShortTweet(strings.Trim(buf.String(), "\n \t")) + arg.CrossLink(s.Config.Site_url, isFriend)

	if isFriend {
		s.sendDM(toIdentityId, data.ToUserName, tweet)
	} else {
		s.sendTweet(tweet)
	}
}

func (s *Twitter_job) sendTweet(t string) {
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

func (s *Twitter_job) sendDM(identityId uint64, toUserName string, t string) {
	dm := &twitter_service.DirectMessagesNewArg{
		ClientToken:  s.Config.Twitter.Client_token,
		ClientSecret: s.Config.Twitter.Client_secret,
		AccessToken:  s.Config.Twitter.Access_token,
		AccessSecret: s.Config.Twitter.Access_secret,
		Message:      t,
		ToUserName:   &toUserName,
		IdentityId:   &identityId,
	}
	s.Senddm.Send(dm)
}
