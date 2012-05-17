package exfe_service

import (
	"github.com/simonz05/godis"
	"exfe/model"
	"gobus"
	"log/syslog"
	"apn/service"
	"fmt"
	"time"
	"bytes"
	"text/template"
)

const crossApnQueueName = "exfe:queue:cross:apn"

type CrossApn struct {
	queue *gobus.TailDelayQueue
	config *Config
	log *syslog.Writer
	client *gobus.Client
}

func NewCrossApn(config *Config) *CrossApn {
	arg := []OneIdentityUpdateArg{}
	redis := godis.New(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)
	log, err := syslog.New(syslog.LOG_DEBUG, "exfe.cross.apn")
	if err != nil {
		panic(err)
	}
	queue, err := gobus.NewTailDelayQueue(crossApnQueueName, config.Cross.Twitter_delay, arg, redis)
	if err != nil {
		panic(err)
	}
	return &CrossApn{
		queue: queue,
		config: config,
		log: log,
		client: gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "apn"),
	}
}

func (s *CrossApn) Serve() {
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

func (s *CrossApn) handle(args []OneIdentityUpdateArg) {
	old_cross := args[0].Old_cross
	cross := &args[len(args)-1].Cross
	to_identity := &args[0].To_identity

	s.sendNewCross(to_identity, old_cross, cross)
	s.sendCrossChange(to_identity, old_cross, cross)
	s.sendExfeeChange(to_identity, old_cross, cross)
}

func (s *CrossApn) findToken(to *exfe_model.Identity, cross *exfe_model.Cross) *string {
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

func (s *CrossApn) createInvitationData(siteUrl string, to *exfe_model.Identity, cross *exfe_model.Cross) *NewInvitationData {
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

func (s *CrossApn) sendNewCross(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old != nil {
		return
	}

	s.sendInvitation(to, current)
}

func (s *CrossApn) sendInvitation(to *exfe_model.Identity, cross *exfe_model.Cross) {
	data := s.createInvitationData(s.config.Site_url, to, cross)
	if data == nil {
		s.log.Err(fmt.Sprintf("Can't send cross %d invitation to identity %d", cross.Id, to.Id))
		return
	}

	buf := bytes.NewBuffer(nil)
	tmpl := template.Must(template.New("NewInvitation").Parse(
		"{{ if .IsHost }}You're successfully gathering this X{{ else }}Invitation{{ end }}: {{ .Title }}.{{ if .Time }} {{ .Time }}{{ end }}{{ if .Place }} at {{ .Place }}{{ end }}"))
	tmpl.Execute(buf, data)
	arg := apn_service.ApnSendArg{
		DeviceToken: to.External_id,
		Alert: buf.String(),
		Badge: 0,
		Sound: "default",
		Cid: cross.Id,
		T: "i",
	}
	s.client.Send("ApnSend", &arg, 5)
}

func (s *CrossApn) sendQuit(to *exfe_model.Identity, cross *exfe_model.Cross) {
	msg := fmt.Sprintf("You quit the Cross %s", cross.Title)
	arg := apn_service.ApnSendArg{
		DeviceToken: to.External_id,
		Alert: msg,
		Badge: 0,
		Sound: "default",
		Cid: cross.Id,
		T: "r",
	}
	s.client.Send("ApnSend", &arg, 5)
}

func (s *CrossApn) sendCrossChange(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
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
	arg := apn_service.ApnSendArg{
		DeviceToken: to.External_id,
		Alert: msg,
		Badge: 0,
		Sound: "default",
		Cid: current.Id,
		T: "u",
	}
	s.client.Send("ApnSend", &arg, 5)
}

func (s *CrossApn) sendExfeeChange(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old == nil {
		return
	}
	accepted, declined, newlyInvited, removed := newStatusUser(s.log, &old.Exfee, &current.Exfee)

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

func (s *CrossApn) sendAccepted(to *exfe_model.Identity, identities map[uint64]*exfe_model.Identity, cross *exfe_model.Cross) {
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
	arg := apn_service.ApnSendArg{
		DeviceToken: to.External_id,
		Alert: msg,
		Badge: 0,
		Sound: "default",
		Cid: cross.Id,
		T: "u",
	}
	s.client.Send("ApnSend", &arg, 5)
}

func (s *CrossApn) sendDeclined(to *exfe_model.Identity, identities map[uint64]*exfe_model.Identity, cross *exfe_model.Cross) {
	msg := "Declined:"
	for _, i := range identities {
		msg = fmt.Sprintf("%s %s,", msg, i.Name)
	}
	msg = msg[0:len(msg) - 1]

	msg = fmt.Sprintf("Cross %s %s", cross.Title, msg)
	arg := apn_service.ApnSendArg{
		DeviceToken: to.External_id,
		Alert: msg,
		Badge: 0,
		Sound: "default",
		Cid: cross.Id,
		T: "u",
	}
	s.client.Send("ApnSend", &arg, 5)
}

func (s *CrossApn) sendNewlyInvited(to *exfe_model.Identity, invitations map[uint64]*exfe_model.Invitation, cross *exfe_model.Cross) {
	msg := "Newly invited:"
	for _, i := range invitations {
		msg = fmt.Sprintf("%s %s,", msg, i.Identity.Name)
	}
	msg = msg[0:len(msg) - 1]

	msg = fmt.Sprintf("Cross %s %s", cross.Title, msg)
	arg := apn_service.ApnSendArg{
		DeviceToken: to.External_id,
		Alert: msg,
		Badge: 0,
		Sound: "default",
		Cid: cross.Id,
		T: "u",
	}
	s.client.Send("ApnSend", &arg, 5)
}

func (s *CrossApn) sendRemoved(to *exfe_model.Identity, identities map[uint64]*exfe_model.Identity, cross *exfe_model.Cross) {
	msg := "Removed:"
	for _, i := range identities {
		msg = fmt.Sprintf("%s %s,", msg, i.Name)
	}
	msg = msg[0:len(msg) - 1]

	msg = fmt.Sprintf("Cross %s %s", cross.Title, msg)
	arg := apn_service.ApnSendArg{
		DeviceToken: to.External_id,
		Alert: msg,
		Badge: 0,
		Sound: "default",
		Cid: cross.Id,
		T: "u",
	}
	s.client.Send("ApnSend", &arg, 5)
}
