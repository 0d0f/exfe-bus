package twitter_job

import (
	"exfe"
	"twitter/service"
	"bytes"
	"fmt"
	"gobus"
	"log"
	"strings"
	"text/template"
)

type TwitterJobArg struct {
	Cross exfe.Cross
	To_invitation exfe.Invitation
}

func (s *TwitterJobArg) isHost() bool {
	return s.Cross.By_identity.Id == s.To_invitation.Identity.Id
}

func (t *TwitterJobArg) Queue() string {
	return "twitter"
}

type Twitter_job struct {
	Config        *Config
	Client        *gobus.Client
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
	Place         string
	SiteUrl       string
	CrossIdBase62 string
	Token         string
}

func (s *TwitterJobArg) CreateData(siteUrl string) *TemplateData {
	t, err := s.Cross.Time.StringInZone(s.To_invitation.Identity.Timezone)
	if err != nil {
		log.Printf("Time parse error: %s", err)
		return nil
	}
	return &TemplateData{
		ToUserName:    s.To_invitation.Identity.External_username,
		IsHost:        s.isHost(),
		Title:         s.Cross.Title,
		Time:          t,
		Place:         s.Cross.Place.String(),
		SiteUrl:       siteUrl,
		CrossIdBase62: s.Cross.Id_base62,
		Token:         s.To_invitation.Token,
	}
}

func LoadTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles(fmt.Sprintf("./template/default/%s", name)))
}

func (s *Twitter_job) Perform(arg *TwitterJobArg) {
	log.Printf("[TwitterJob]Get a job")

	toIdentityId := arg.To_invitation.Identity.Id

	if arg.To_invitation.Identity.External_id == "" {
		// get to_identity info
		s.Client.Send("GetInfo", &twitter_service.UsersShowArg{
			ClientToken:  s.Config.Twitter.Client_token,
			ClientSecret: s.Config.Twitter.Client_secret,
			AccessToken:  s.Config.Twitter.Access_token,
			AccessSecret: s.Config.Twitter.Access_secret,
			ScreenName:   &arg.To_invitation.Identity.External_username,
			IdentityId:   &toIdentityId,
		}, 5)
	}

	// check friendship
	f := &twitter_service.FriendshipsExistsArg{
		ClientToken:  s.Config.Twitter.Client_token,
		ClientSecret: s.Config.Twitter.Client_secret,
		AccessToken:  s.Config.Twitter.Access_token,
		AccessSecret: s.Config.Twitter.Access_secret,
		UserA:        arg.To_invitation.Identity.External_username,
		UserB:        s.Config.Twitter.Screen_name,
	}
	var isFriend bool
	err := s.Client.Do("GetFriendship", f, &isFriend)
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

	if isFriend {
		tweet := ShortTweet(strings.Trim(buf.String(), "\n \t")) + arg.Cross.LinkTo(s.Config.Site_url, &arg.To_invitation)
		s.sendDM(toIdentityId, data.ToUserName, tweet)
	} else {
		tweet := ShortTweet(strings.Trim(buf.String(), "\n \t")) + arg.Cross.Link(s.Config.Site_url)
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
	err := s.Client.Do("SendTweet", tweet, &response)
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
	s.Client.Send("SendDM", dm, 5)
}
