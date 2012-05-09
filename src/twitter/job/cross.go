package twitter_job

import (
	"exfe"
	"gobus"
	"log"
	"text/template"
	"fmt"
	"twitter/service"
	"bytes"
	"strings"
)

func ShortTweet(tweet string) string {
	const linkLength = 25
	if (len(tweet) + linkLength) > 140 {
		return tweet[0:(137-linkLength)] + "..."
	}
	return tweet
}

func LoadTemplate(name string) *template.Template {
	return template.Must(template.ParseFiles(fmt.Sprintf("./template/default/%s", name)))
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

type UpdateExfeeArg struct {
	Cross exfe.Cross
	Event map[string][]exfe.Invitation
	By_identity exfe.Identity
	To_invitation exfe.Invitation
}

func (a *UpdateExfeeArg) isHost() bool {
	return a.Cross.By_identity.Id == a.To_invitation.Identity.Id
}

func (a *UpdateExfeeArg) CreateData(siteUrl string) *TemplateData {
	t, err := a.Cross.Time.StringInZone(a.To_invitation.Identity.Timezone)
	if err != nil {
		log.Printf("Time parse error: %s", err)
		return nil
	}
	return &TemplateData{
		ToUserName:    a.To_invitation.Identity.External_username,
		IsHost:        a.isHost(),
		Title:         a.Cross.Title,
		Time:          t,
		Place:         a.Cross.Place.String(),
		SiteUrl:       siteUrl,
		CrossIdBase62: a.Cross.Id_base62,
		Token:         a.To_invitation.Token,
	}
}

type CrossJob struct {
	Config *Config
	Client *gobus.Client
}

func (j *CrossJob) Update_exfee(args []*UpdateExfeeArg) error {
	log.Printf("[INFO][Update Exfee]Get a job")

	for _, arg := range args {
		toIdentityId := arg.To_invitation.Identity.Id

		if arg.To_invitation.Identity.External_id == "" {
			// get to_identity info
			j.Client.Send("GetInfo", &twitter_service.UsersShowArg{
				ClientToken:  j.Config.Twitter.Client_token,
				ClientSecret: j.Config.Twitter.Client_secret,
				AccessToken:  j.Config.Twitter.Access_token,
				AccessSecret: j.Config.Twitter.Access_secret,
				ScreenName:   &arg.To_invitation.Identity.External_username,
				IdentityId:   &toIdentityId,
			}, 5)
		}

		// check friendship
		f := &twitter_service.FriendshipsExistsArg{
			ClientToken:  j.Config.Twitter.Client_token,
			ClientSecret: j.Config.Twitter.Client_secret,
			AccessToken:  j.Config.Twitter.Access_token,
			AccessSecret: j.Config.Twitter.Access_secret,
			UserA:        arg.To_invitation.Identity.External_username,
			UserB:        j.Config.Twitter.Screen_name,
		}
		var isFriend bool
		err := j.Client.Do("GetFriendship", f, &isFriend, 3)
		if err != nil {
			log.Printf("[ERROR][Update Exfee]Can't require friendship: %s(%d)", arg.To_invitation.Identity.External_username, arg.To_invitation.Identity.External_id)
			isFriend = false
		}

		data := arg.CreateData(j.Config.Site_url)

		var tmpl *template.Template
		if isFriend {
			tmpl = LoadTemplate("twitter_sender_dm.tmpl")
		} else {
			tmpl = LoadTemplate("twitter_sender_tweet.tmpl")
		}
		buf := bytes.NewBuffer(nil)
		tmpl.Execute(buf, data)

		if isFriend {
			tweet := ShortTweet(strings.Trim(buf.String(), "\n \t")) + " " + arg.Cross.LinkTo(j.Config.Site_url, &arg.To_invitation)
			j.sendDM(toIdentityId, data.ToUserName, tweet)
		} else {
			tweet := ShortTweet(strings.Trim(buf.String(), "\n \t")) + " " + arg.Cross.Link(j.Config.Site_url)
			j.sendTweet(tweet)
		}
	}
}

func (j *CrossJob) sendTweet(t string) {
	tweet := &twitter_service.StatusesUpdateArg{
		ClientToken:  j.Config.Twitter.Client_token,
		ClientSecret: j.Config.Twitter.Client_secret,
		AccessToken:  j.Config.Twitter.Access_token,
		AccessSecret: j.Config.Twitter.Access_secret,
		Tweet:        t,
	}
	var response twitter_service.StatusesUpdateReply
	j.Client.Send("SendTweet", tweet, 5)
}

func (j *CrossJob) sendDM(identityId uint64, toUserName string, t string) {
	dm := &twitter_service.DirectMessagesNewArg{
		ClientToken:  j.Config.Twitter.Client_token,
		ClientSecret: j.Config.Twitter.Client_secret,
		AccessToken:  j.Config.Twitter.Access_token,
		AccessSecret: j.Config.Twitter.Access_secret,
		Message:      t,
		ToUserName:   &toUserName,
		IdentityId:   &identityId,
	}
	j.Client.Send("SendDM", dm, 5)
}
