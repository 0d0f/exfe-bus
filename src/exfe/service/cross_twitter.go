package exfe_service

import (
	"exfe/model"
	"twitter/service"
	"fmt"
	"bytes"
	"text/template"
	"strings"
)

type CrossTwitter struct {
	CrossProviderBase
}

func NewCrossTwitter(config *Config) (ret *CrossTwitter) {
	ret = &CrossTwitter {
		CrossProviderBase: NewCrossProviderBase("twitter", config),
	}
	ret.handler = ret
	return
}

func (s *CrossTwitter) Handle(to_identity *exfe_model.Identity, old_cross, cross *exfe_model.Cross) {
	s.sendNewCross(to_identity, old_cross, cross)
	s.sendCrossChange(to_identity, old_cross, cross)
	s.sendExfeeChange(to_identity, old_cross, cross)
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

func (s *CrossTwitter) sendNewCross(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old != nil {
		return
	}

	s.sendInvitation(to, current)
}

func (s *CrossTwitter) sendInvitation(to *exfe_model.Identity, cross *exfe_model.Cross) {
	data := newInvitationData(s.log, s.config.Site_url, to, cross)
	if data == nil {
		s.log.Err(fmt.Sprintf("Can't send cross %d invitation to identity %d", cross.Id, to.Id))
		return
	}

	s.getIdentityInfo(to)
	isFriend := s.checkFriend(to)

	buf := bytes.NewBuffer(nil)
	if isFriend {
		tmpl := template.Must(template.New("NewInvitation").Parse(
			"{{ if .IsHost }}You're successfully gathering this X{{ else }}Invitation{{ end }}: {{ .Title }}.{{ if .Time }} {{ .Time }}{{ end }}{{ if .Place }} at {{ .Place }}{{ end }}"))
		tmpl.Execute(buf, data)
		msg := s.shortTweet(strings.Trim(buf.String(), "\n \t")) + " " + cross.LinkTo(s.config.Site_url, data.Token)
		s.sendDM(to.Id, data.ToUserName, msg)
	} else {
		tmpl := template.Must(template.New("NewInvitation").Parse(
			"@{{ .ToUserName }} {{ if .IsHost }}Invited{{ else }}Invitation{{ end }}:"))
		tmpl.Execute(buf, data)
		tweet := s.shortTweet(strings.Trim(buf.String(), "\n \t")) + " " + cross.Link(s.config.Site_url)
		s.sendTweet(tweet)
	}
}

func (s *CrossTwitter) sendQuit(to *exfe_model.Identity, cross *exfe_model.Cross) {
	s.getIdentityInfo(to)
	isFriend := s.checkFriend(to)

	msg := fmt.Sprintf("You quit the Cross %s", cross.Link(s.config.Site_url))
	if isFriend {
		s.sendDM(to.Id, to.External_username, msg)
	} else {
		tweet := fmt.Sprintf("@%s %s", to.External_username, msg, "\n \t")
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
	accepted, declined, newlyInvited, removed := diffExfee(s.log, &old.Exfee, &current.Exfee)

	if len(accepted) > 0 {
		s.sendAccepted(to, accepted, current)
	}
	if len(declined) > 0 {
		s.sendDeclined(to, declined, current)
	}
	if len(newlyInvited) > 0 {
		if _, ok := newlyInvited[to.Connected_user_id]; ok {
			s.sendInvitation(to, current)
		} else {
			s.sendNewlyInvited(to, newlyInvited, current)
		}
	}
	if len(removed) > 0 {
		if _, ok := removed[to.Connected_user_id]; ok {
			s.sendQuit(to, current)
		} else {
			s.sendRemoved(to, removed, current)
		}
	}
}

func (s *CrossTwitter) sendAccepted(to *exfe_model.Identity, identities map[uint64]*exfe_model.Identity, cross *exfe_model.Cross) {
	totalAccepted := 0
	for _, i := range cross.Exfee.Invitations {
		if i.Rsvp_status == "ACCEPTED" {
			totalAccepted++
		}
	}
	msg := fmt.Sprintf("%d Accepted:", totalAccepted)
	for _, i := range identities {
		msg = fmt.Sprintf("%s %s,", msg, i.Name)
	}
	otherCount := totalAccepted - len(identities)
	switch otherCount {
	case 0:
		msg = msg[0:len(msg) - 1]
	case 1:
		msg = fmt.Sprintf("%s and 1 other", msg)
	default:
		msg = fmt.Sprintf("%s and %d others", msg, totalAccepted - len(identities))
	}

	s.getIdentityInfo(to)
	isFriend := s.checkFriend(to)

	if isFriend {
		msg = fmt.Sprintf("Cross %s(%s) %s", cross.Title, cross.Link(s.config.Site_url), msg)
		s.sendDM(to.Id, to.External_username, msg)
	} else {
		tweet := fmt.Sprintf("@s Cross %s(%s) %s", to.External_username, cross.Title, cross.Link(s.config.Site_url), msg)
		s.sendTweet(tweet)
	}
}

func (s *CrossTwitter) sendDeclined(to *exfe_model.Identity, identities map[uint64]*exfe_model.Identity, cross *exfe_model.Cross) {
	msg := "Declined:"
	for _, i := range identities {
		msg = fmt.Sprintf("%s %s,", msg, i.Name)
	}
	msg = msg[0:len(msg) - 1]

	s.getIdentityInfo(to)
	isFriend := s.checkFriend(to)

	if isFriend {
		msg = fmt.Sprintf("Cross %s %s %s", cross.Title, cross.Link(s.config.Site_url), msg)
		s.sendDM(to.Id, to.External_username, msg)
	} else {
		tweet := fmt.Sprintf("@%s Cross %s %s %s", to.External_username, cross.Title, cross.Link(s.config.Site_url), msg)
		s.sendTweet(tweet)
	}
}

func (s *CrossTwitter) sendNewlyInvited(to *exfe_model.Identity, invitations map[uint64]*exfe_model.Invitation, cross *exfe_model.Cross) {
	msg := "Newly invited:"
	for _, i := range invitations {
		msg = fmt.Sprintf("%s %s,", msg, i.Identity.Name)
	}
	msg = msg[0:len(msg) - 1]

	s.getIdentityInfo(to)
	isFriend := s.checkFriend(to)

	if isFriend {
		msg = fmt.Sprintf("Cross %s(%s) %s", cross.Title, cross.Link(s.config.Site_url), msg)
		s.sendDM(to.Id, to.External_username, msg)
	} else {
		tweet := fmt.Sprintf("@s Cross %s(%s) %s", to.External_username, cross.Title, cross.Link(s.config.Site_url), msg)
		s.sendTweet(tweet)
	}
}

func (s *CrossTwitter) sendRemoved(to *exfe_model.Identity, identities map[uint64]*exfe_model.Identity, cross *exfe_model.Cross) {
	msg := "Removed:"
	for _, i := range identities {
		msg = fmt.Sprintf("%s %s,", msg, i.Name)
	}
	msg = msg[0:len(msg) - 1]

	s.getIdentityInfo(to)
	isFriend := s.checkFriend(to)

	if isFriend {
		msg = fmt.Sprintf("Cross %s(%s) %s", cross.Title, cross.Link(s.config.Site_url), msg)
		s.sendDM(to.Id, to.External_username, msg)
	} else {
		tweet := fmt.Sprintf("@s Cross %s(%s) %s", to.External_username, cross.Title, cross.Link(s.config.Site_url), msg)
		s.sendTweet(tweet)
	}}
