package exfe_service

import (
	"exfe/model"
	"apn/service"
	"c2dm/service"
	"gobus"
	"fmt"
	"bytes"
	"text/template"
)

type CrossPush struct {
	CrossProviderBase
	android *gobus.Client
}

func NewCrossPush(config *Config) (ret *CrossPush) {
	ret = &CrossPush{
		CrossProviderBase: NewCrossProviderBase("push", config),
	}
	ret.client = gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "iOSAPN")
	ret.android = gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "Android")
	ret.handler = ret
	return
}

func (s *CrossPush) Handle(to_identity *exfe_model.Identity, old_cross, cross *exfe_model.Cross) {
	s.sendNewCross(to_identity, old_cross, cross)
	s.sendCrossChange(to_identity, old_cross, cross)
	s.sendExfeeChange(to_identity, old_cross, cross)
}

func (s *CrossPush) sendNewCross(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old != nil {
		return
	}

	s.sendInvitation(to, current)
}

func (s *CrossPush) sendInvitation(to *exfe_model.Identity, cross *exfe_model.Cross) {
	data := newInvitationData(s.log, s.config.Site_url, to, cross)
	if data == nil {
		s.log.Err(fmt.Sprintf("Can't send cross %d invitation to identity %d", cross.Id, to.Id))
		return
	}

	buf := bytes.NewBuffer(nil)
	tmpl := template.Must(template.New("NewInvitation").Parse(
		"{{ if .IsHost }}You're successfully gathering this X{{ else }}Invitation{{ end }}: {{ .Title }}.{{ if .Time }} {{ .Time }}{{ end }}{{ if .Place }} at {{ .Place }}{{ end }}"))
	tmpl.Execute(buf, data)

	s.push(to, buf.String(), "default", "i", 0, cross.Id)
}

func (s *CrossPush) sendQuit(to *exfe_model.Identity, cross *exfe_model.Cross) {
	msg := fmt.Sprintf("You quit the Cross %s", cross.Title)
	s.push(to, msg, "default", "r", 0, cross.Id)
}

func (s *CrossPush) sendCrossChange(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
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

	msg := fmt.Sprintf("Update: %s", message)
	s.push(to, msg, "default", "u", 0, current.Id)
}

func (s *CrossPush) sendExfeeChange(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
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

func (s *CrossPush) sendAccepted(to *exfe_model.Identity, identities map[uint64]*exfe_model.Identity, cross *exfe_model.Cross) {
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

	msg = fmt.Sprintf("Cross %s %s", cross.Title, msg)
	s.push(to, msg, "default", "u", 0, cross.Id)
}

func (s *CrossPush) sendDeclined(to *exfe_model.Identity, identities map[uint64]*exfe_model.Identity, cross *exfe_model.Cross) {
	msg := "Declined:"
	for _, i := range identities {
		msg = fmt.Sprintf("%s %s,", msg, i.Name)
	}
	msg = msg[0:len(msg) - 1]

	msg = fmt.Sprintf("Cross %s %s", cross.Title, msg)
	s.push(to, msg, "default", "u", 0, cross.Id)
}

func (s *CrossPush) sendNewlyInvited(to *exfe_model.Identity, invitations map[uint64]*exfe_model.Invitation, cross *exfe_model.Cross) {
	msg := "Newly invited:"
	for _, i := range invitations {
		msg = fmt.Sprintf("%s %s,", msg, i.Identity.Name)
	}
	msg = msg[0:len(msg) - 1]

	msg = fmt.Sprintf("Cross %s %s", cross.Title, msg)
	s.push(to, msg, "default", "u", 0, cross.Id)
}

func (s *CrossPush) sendRemoved(to *exfe_model.Identity, identities map[uint64]*exfe_model.Identity, cross *exfe_model.Cross) {
	msg := "Removed:"
	for _, i := range identities {
		msg = fmt.Sprintf("%s %s,", msg, i.Name)
	}
	msg = msg[0:len(msg) - 1]

	msg = fmt.Sprintf("Cross %s %s", cross.Title, msg)
	s.push(to, msg, "default", "u", 0, cross.Id)
}

func (s *CrossPush) push(to *exfe_model.Identity, message, sound, messageType string, badge uint, crossId uint64) {
	switch to.Provider {
	case "iOSAPN":
		arg := apn_service.ApnSendArg{
			DeviceToken: to.External_id,
			Alert: message,
			Badge: badge,
			Sound: sound,
			Cid: crossId,
			T: messageType,
		}
		s.client.Send("ApnSend", &arg, 5)
	case "Android":
		arg := c2dm_service.C2DMSendArg{
			DeviceID: to.External_id,
			Message: message,
		}
		s.android.Send("C2DMSend", &arg, 5)
	}
}

