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
	s.sendDiscard(to_identity, old_cross, cross)
	s.sendCrossChange(to_identity, old_cross, cross)
	s.sendExfeeChange(to_identity, old_cross, cross)
}

func (s *CrossTwitter) findToken(to *exfe_model.Identity, cross *exfe_model.Cross) *string {
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
	return token
}

func (s *CrossTwitter) diffExfee(left, right *exfe_model.Exfee) (leftOnly map[uint64]*exfe_model.Identity, rightOnly map[uint64]*exfe_model.Identity) {
	leftOnly = make(map[uint64]*exfe_model.Identity)
	rightOnly = make(map[uint64]*exfe_model.Identity)

	for _, i := range left.Invitations {
		leftOnly[i.Identity.Id] = &i.Identity
	}
	for _, i := range right.Invitations {
		rightOnly[i.Identity.Id] = &i.Identity
	}
	same := make([]uint64, 0, 0)
	for k, _ := range leftOnly {
		if _, ok := rightOnly[k]; ok {
			same = append(same, k)
		}
	}
	fmt.Println("leftOnly:", leftOnly)
	fmt.Println("rightOnly:", rightOnly)
	fmt.Println("same:", same)
	for _, id := range same {
		delete(leftOnly, id)
		delete(rightOnly, id)
	}
	return
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
	return &NewInvitationData{
		ToUserName:    to.External_username,
		IsHost:        isHost,
		Title:         cross.Title,
		Time:          t,
		Place:         cross.Place.String(),
		SiteUrl:       siteUrl,
		CrossIdBase62: cross.Id_base62,
		Token:         *s.findToken(to, cross),
	}
}

func (s *CrossTwitter) sendNewInvitation(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old != nil {
		_, right := s.diffExfee(&old.Exfee, &current.Exfee)
		if _, ok := right[to.Id]; !ok {
			return
		}
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

func (s *CrossTwitter) sendDiscard(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old == nil {
		return
	}
	left, _ := s.diffExfee(&old.Exfee, &current.Exfee)
	fmt.Println(left)
	if _, ok := left[to.Id]; !ok {
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
	tmpl := template.Must(template.New("Discard").Parse(
		"You're discarded from this X"))
	tmpl.Execute(buf, data)
	if isFriend {
		msg := s.shortTweet(strings.Trim(buf.String(), "\n \t")) + " " + current.LinkTo(s.config.Site_url, data.Token)
		s.sendDM(to.Id, data.ToUserName, msg)
	} else {
		tweet := fmt.Sprintf("@%s %s %s", data.ToUserName, s.shortTweet(strings.Trim(buf.String(), "\n \t")), current.Link(s.config.Site_url))
		s.sendTweet(tweet)
	}
}

const messageMaxLen = 140 - 29 /* len("Update http://t.co/fbqqsjky:\n") */ - 5 /* reserved */
const titleMaxLen = 20
const newTitleMaxLen = 13

func sameTitleMessage(time, title, place1, place2 string) string {
	meta := fmt.Sprintf("%s \n%s \n%s", time, place1, place2)

	if len(meta) + len(title) + 2 > messageMaxLen {
		title = strings.Trim(title[0:titleMaxLen], " \n") + "…"
	}
	if len(meta) + len(title) + 2 > messageMaxLen {
		metaLen := messageMaxLen - len(title) - 5
		meta = strings.Trim(meta[0:metaLen], " \n") + "…"
	}
	return fmt.Sprintf("%s \n%s", meta, title)
}

func diffTitleMessage(time, new_title, place1, place2, old_title string) string {
	meta := fmt.Sprintf("%s \n%s \n%s", time, place1, place2)
	title := fmt.Sprintf("\"%s\"\nchanged from \"%s\"", new_title, old_title)

	if len(meta) + len(title) + 2 > messageMaxLen {
		new_title = strings.Trim(new_title[0:newTitleMaxLen], " \n") + "…"
		title = fmt.Sprintf("\"%s\"\nchanged from \"%s\"", new_title, old_title)
	}
	if len(meta) + len(title) + 2 > messageMaxLen {
		old_title = strings.Trim(old_title[0:titleMaxLen], " \n") + "…"
		title = fmt.Sprintf("\"%s\"\nchanged from \"%s\"", new_title, old_title)
	}
	if len(meta) + len(title) + 2 > messageMaxLen {
		metaLen := messageMaxLen - len(title) - 5
		meta = strings.Trim(meta[0:metaLen], " \n") + "…"
	}
	return fmt.Sprintf("%s \n%s", meta, title)
}

func (s *CrossTwitter) sendCrossChange(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old == nil {
		return
	}

	newTime, err := current.Time.StringInZone(to.Timezone)
	if err != nil {
		s.log.Err(fmt.Sprintf("can't convert cross %d time to zone %s", current.Id, to.Timezone))
		return
	}
	newPlaceTitle := current.Place.Title
	newPlaceDesc := current.Place.Description
	isChanged := false

	if old.Title != current.Title {
		isChanged = true
	}
	if old.Place.Title != newPlaceTitle {
		isChanged = true
	}
	if old.Place.Description != newPlaceDesc {
		isChanged = true
	}
	if o, _ := old.Time.StringInZone(to.Timezone); o != newTime {
		isChanged = true
	}
	if !isChanged {
		return
	}

	var message string
	if old.Title != current.Title {
		message = diffTitleMessage(newTime, current.Title, newPlaceTitle, newPlaceDesc, old.Title)
	} else {
		message = sameTitleMessage(newTime, current.Title, newPlaceTitle, newPlaceDesc)
	}

	s.getIdentityInfo(to)
	isFriend := s.checkFriend(to)

	if isFriend {
		msg := fmt.Sprintf("Update %s:\n%s", current.LinkTo(s.config.Site_url, *s.findToken(to, current)), message)
		s.sendDM(to.Id, to.External_username, msg)
	} else {
		tweet := fmt.Sprintf("@%s Update %s:\n%s", to.External_username, current.Link(s.config.Site_url), message)
		s.sendTweet(tweet)
	}
}

func (s *CrossTwitter) sendExfeeChange(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old == nil {
		return
	}
	left, right := s.diffExfee(&old.Exfee, &current.Exfee)
	if len(left) == 0 && len(right) == 0 {
		return
	}
	if _, ok := left[to.Id]; ok {
		return
	}
	if _, ok := right[to.Id]; ok {
		return
	}

	s.getIdentityInfo(to)
	isFriend := s.checkFriend(to)

	var message string
	switch len(right) {
	default:
		message += fmt.Sprintf("Confirmed: %s /w %d others", right[0].External_username, len(right))
	case 2:
		message += fmt.Sprintf("Confirmed: %s /w 1 other", right[0].External_username)
	case 1:
		message += fmt.Sprintf("Confirmed: %s", right[0].External_username)
	case 0:
	}
	if len(message) > 0 {
		message += "\n"
	}
	switch len(left) {
	default:
		message += fmt.Sprintf("Decline: %s /w %d others", left[0].External_username, len(left))
	case 2:
		message += fmt.Sprintf("Decline: %s /w 1 other", left[0].External_username)
	case 1:
		message += fmt.Sprintf("Decline: %s", left[0].External_username)
	case 0:
	}

	if isFriend {
		msg := fmt.Sprintf("Update %s:\n%s", current.LinkTo(s.config.Site_url, *s.findToken(to, current)), message)
		s.sendDM(to.Id, to.External_username, msg)
	} else {
		tweet := fmt.Sprintf("@%s Update %s:\n%s",to.External_username, current.Link(s.config.Site_url), message)
		s.sendTweet(tweet)
	}
}
