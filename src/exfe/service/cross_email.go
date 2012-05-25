package exfe_service

import (
	"exfe/model"
	"gomail"
	"fmt"
	"bytes"
	"text/template"
)

type CrossEmail struct {
	CrossProviderBase
}

func NewCrossEmail(config *Config) (ret *CrossEmail) {
	ret = &CrossEmail{
		CrossProviderBase: NewCrossProviderBase("email", config),
	}
	ret.handler = ret
	return
}

func (s *CrossEmail) Handle(to_identity *exfe_model.Identity, old_cross, cross *exfe_model.Cross) {
	s.sendNewCross(to_identity, old_cross, cross)
	s.sendCrossChange(to_identity, old_cross, cross)
}

func (s *CrossEmail) sendNewCross(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old != nil {
		return
	}

	s.sendInvitation(to, current)
}

type CrossInvitation struct {
	To *exfe_model.Invitation
	Cross *exfe_model.Cross
	Site_url string
	App_url string
}

func findInvitation(to *exfe_model.Identity, cross *exfe_model.Cross) *exfe_model.Invitation {
	for i := range cross.Exfee.Invitations {
		if to.Id == cross.Exfee.Invitations[i].Identity.Id {
			return &cross.Exfee.Invitations[i]
		}
	}
	return nil
}

func (s *CrossEmail) sendInvitation(to *exfe_model.Identity, cross *exfe_model.Cross) {
	data := CrossInvitation{findInvitation(to, cross), cross, s.config.Site_url, "appurl"}

	buf := bytes.NewBuffer(nil)
	tmpl := template.Must(template.ParseFiles("./template/default/cross_invitation_email.html"))
	err := tmpl.Execute(buf, data)
	fmt.Println("template exec:", err)
	arg := gomail.Mail{
		To: []gomail.MailUser{gomail.MailUser{to.External_id, to.Name}},
		From: gomail.MailUser{"x@exfe.com", "x@exfe.com"},
		Subject: cross.Title,
		Html: buf.String(),
	}
	fmt.Println(buf.String())
	s.client.Send("EmailSend", &arg, 5)
}

type CrossChange struct {
	To *exfe_model.Invitation
	Cross *exfe_model.Cross
	OldCross *exfe_model.Cross
	Site_url string
	App_url string
}

func (s *CrossEmail) sendCrossChange(to *exfe_model.Identity, old *exfe_model.Cross, current *exfe_model.Cross) {
	if old == nil {
		return
	}

	data := CrossChange{findInvitation(to, current), current, old, s.config.Site_url, "appurl"}

	buf := bytes.NewBuffer(nil)
	tmpl := template.Must(template.ParseFiles("./template/default/cross_change_email.html"))
	err := tmpl.Execute(buf, data)
	fmt.Println("template exec:", err)
	arg := gomail.Mail{
		To: []gomail.MailUser{gomail.MailUser{to.External_id, to.Name}},
		From: gomail.MailUser{"x@exfe.com", "x@exfe.com"},
		Subject: current.Title,
		Html: buf.String(),
	}
	fmt.Println(buf.String())
	s.client.Send("EmailSend", &arg, 5)
}
