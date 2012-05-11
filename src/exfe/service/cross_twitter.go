package exfe_service

import (
	"github.com/simonz05/godis"
	"exfe/model"
	"gobus"
	"log/syslog"
	"twitter/service"
	"fmt"
	"time"
	"bytes"
	"text/template"
	"strings"
)

const crossTwitterQueueName = "exfe:queue:cross:twitter"

type CrossTwitter struct {
	queue *gobus.TailDelayQueue
	config *Config
	log *syslog.Writer
	client *gobus.Client
}

func NewCrossTwitter(config *Config) *CrossTwitter {
	arg := []OneIdentityUpdateArg{}
	redis := godis.New(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)
	log, err := syslog.New(syslog.LOG_DEBUG, "exfe.cross.twitter")
	if err != nil {
		panic(err)
	}
	queue, err := gobus.NewTailDelayQueue(crossTwitterQueueName, config.Cross.Twitter_delay, arg, redis)
	if err != nil {
		panic(err)
	}
	return &CrossTwitter{
		queue: queue,
		config: config,
		log: log,
		client: gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "twitter"),
	}
}

func (s *CrossTwitter) Serve() {
	for {
		t, err := s.queue.NextWakeup()
		if err != nil {
			s.log.Crit(fmt.Sprintf("next wakeup error: %s", err))
			break
		}
		time.Sleep(t)
		args, err := s.queue.Pop()
		if err != nil {
			s.log.Err(fmt.Sprintf("pop from delay queue failed: %s", err))
			continue
		}
		if args != nil {
			s.handle(args.([]OneIdentityUpdateArg))
		}
	}
}

func (s *CrossTwitter) shortTweet(tweet string) string {
	const linkLength = 25
	if (len(tweet) + linkLength) > 140 {
		return tweet[0:(137-linkLength)] + "..."
	}
	return tweet
}

func (s *CrossTwitter) getIdentityInfo(id *exfe_model.Identity) {
	if id.External_id != "" {
		return
	}
	// get to_identity info
	s.client.Send("GetInfo", &twitter_service.UsersShowArg{
		ClientToken:  s.config.Twitter.Client_token,
		ClientSecret: s.config.Twitter.Client_secret,
		AccessToken:  s.config.Twitter.Access_token,
		AccessSecret: s.config.Twitter.Access_secret,
		ScreenName:   &id.External_username,
		IdentityId:   &id.Id,
	}, 5)
}

func (s *CrossTwitter) checkFriend(to *exfe_model.Identity) (isFriend bool) {
	f := &twitter_service.FriendshipsExistsArg{
		ClientToken:  s.config.Twitter.Client_token,
		ClientSecret: s.config.Twitter.Client_secret,
		AccessToken:  s.config.Twitter.Access_token,
		AccessSecret: s.config.Twitter.Access_secret,
		UserA:        to.External_username,
		UserB:        s.config.Twitter.Screen_name,
	}
	err := s.client.Do("GetFriendship", f, &isFriend, 10)
	if err != nil {
		s.log.Err(fmt.Sprintf("Can't require identity %d friendship: %s", to.Id, err))
		isFriend = false
	}
	return
}

func (s *CrossTwitter) sendTweet(t string) {
	tweet := &twitter_service.StatusesUpdateArg{
		ClientToken:  s.config.Twitter.Client_token,
		ClientSecret: s.config.Twitter.Client_secret,
		AccessToken:  s.config.Twitter.Access_token,
		AccessSecret: s.config.Twitter.Access_secret,
		Tweet:        t,
	}
	s.client.Send("SendTweet", tweet, 5)
}

func (s *CrossTwitter) sendDM(identityId uint64, toUserName string, t string) {
	dm := &twitter_service.DirectMessagesNewArg{
		ClientToken:  s.config.Twitter.Client_token,
		ClientSecret: s.config.Twitter.Client_secret,
		AccessToken:  s.config.Twitter.Access_token,
		AccessSecret: s.config.Twitter.Access_secret,
		Message:      t,
		ToUserName:   &toUserName,
		IdentityId:   &identityId,
	}
	s.client.Send("SendDM", dm, 5)
}

func (s *CrossTwitter) handle(args []OneIdentityUpdateArg) {
	old_cross := args[0].Old_cross
	cross := &args[len(args)-1].Cross
	to_identity := &args[0].To_identity

	s.sendNewInvitation(to_identity, old_cross, cross)
	s.sendNoticeInvitation(to_identity, old_cross, cross)
	s.sendCrossChange(to_identity, old_cross, cross)
	s.sendExfeeChange(to_identity, old_cross, cross)
}

type NewInvitationData struct {
	ToUserName    string
	IsHost        bool
	Title         string
	Time          string
	Place         string
	SiteUrl       string
	CrossIdBase62 string
	Token         string
}

func (s *CrossTwitter) createInvitationData(siteUrl string, to *exfe_model.Identity, cross *exfe_model.Cross) *NewInvitationData {
	t, err := cross.Time.StringInZone(to.Timezone)
	if err != nil {
		s.log.Err(fmt.Sprintf("Time parse error: %s", err))
		return nil
	}
	isHost := cross.By_identity.Connected_user_id == to.Connected_user_id
	var token *string
	for _, invitation := range cross.Exfee.Invitations {
		if invitation.Identity.Connected_user_id == to.Connected_user_id {
			token = &invitation.Token
			break
		}
	}
	if token == nil {
		s.log.Err(fmt.Sprintf("Can't find identity %d in cross %d", to.Id, cross.Id))
	}
	return &NewInvitationData{
		ToUserName:    to.External_username,
		IsHost:        isHost,
		Title:         cross.Title,
		Time:          t,
		Place:         cross.Place.String(),
		SiteUrl:       siteUrl,
		CrossIdBase62: cross.Id_base62,
		Token:         *token,
	}
}

func (s *CrossTwitter) sendNewInvitation(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old != nil {
		return
	}

	data := s.createInvitationData(s.config.Site_url, to, current)
	if data == nil {
		s.log.Err(fmt.Sprintf("Can't send cross %d invitation to identity %d", current.Id, to.Id))
		return
	}

	s.getIdentityInfo(to)
	isFriend := s.checkFriend(to)

	buf := bytes.NewBuffer(nil)
	if isFriend {
		tmpl := template.Must(template.New("NewInvitation").Parse(
			"{{ if .IsHost }}You're successfully gathering this X{{ else }}Invitation{{ end }}: {{ .Title }}.{{ if .Time }} {{ .Time }}{{ end }}{{ if .Place }} at {{ .Place }}{{ end }}"))
		tmpl.Execute(buf, data)
		msg := s.shortTweet(strings.Trim(buf.String(), "\n \t")) + " " + current.LinkTo(s.config.Site_url, data.Token)
		s.sendDM(to.Id, data.ToUserName, msg)
	} else {
		tmpl := template.Must(template.New("NewInvitation").Parse(
			"@{{ .ToUserName }} {{ if .IsHost }}Invited{{ else }}Invitation{{ end }}:"))
		tmpl.Execute(buf, data)
		tweet := s.shortTweet(strings.Trim(buf.String(), "\n \t")) + " " + current.Link(s.config.Site_url)
		s.sendTweet(tweet)
	}
}

func (s *CrossTwitter) sendNoticeInvitation(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
}

func (s *CrossTwitter) sendCrossChange(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
}

func (s *CrossTwitter) sendExfeeChange(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
}
