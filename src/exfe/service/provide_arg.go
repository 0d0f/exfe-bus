package exfe_service

import (
	"bytes"
	"exfe/model"
	"fmt"
	"log"
	"text/template"
)

func getTemplateString(name string, data interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	tmpl, err := template.New("twitter").Funcs(helper).ParseFiles(fmt.Sprintf("./template/default/%s", name))
	if err != nil {
		return "", err
	}
	err = tmpl.ExecuteTemplate(buf, name, data)
	return buf.String(), err
}

type ProviderArg struct {
	Cross         *exfe_model.Cross
	Old_cross     *exfe_model.Cross
	To_identity   *exfe_model.Identity
	By_identities []*exfe_model.Identity
	Posts         []*exfe_model.Post

	Config *Config

	Accepted     []*exfe_model.Invitation
	Declined     []*exfe_model.Identity
	NewlyInvited []*exfe_model.Invitation
	Removed      []*exfe_model.Identity
}

func (a *ProviderArg) IsCrossChanged() bool {
	return a.IsTitleChanged() || a.IsTimeChanged() || a.IsPlaceChanged() ||
		((len(a.Accepted) + len(a.Declined) + len(a.NewlyInvited) + len(a.Removed)) > 0)
}

func (a *ProviderArg) CrossChangedWithPosts() bool {
	return a.IsCrossChanged() && (len(a.Posts) > 0)
}

func (a *ProviderArg) IsTitleChanged() bool {
	if a.Old_cross == nil {
		return false
	}
	return a.Cross.Title != a.Old_cross.Title
}

func (a *ProviderArg) IsTimeChanged() bool {
	if a.Old_cross == nil {
		return false
	}
	crossTime, _ := a.Cross.Time.StringInZone(a.Timezone())
	oldTime, _ := a.Old_cross.Time.StringInZone(a.Timezone())
	return crossTime != oldTime
}

func (a *ProviderArg) IsPlaceChanged() bool {
	if a.Old_cross == nil {
		return false
	}
	return (a.Cross.Place.Title != a.Old_cross.Place.Title) || (a.Cross.Place.Description != a.Old_cross.Place.Description)
}

func (a *ProviderArg) IsHost() bool {
	return a.Cross.By_identity.DiffId() == a.To_identity.DiffId()
}

func (a *ProviderArg) PublicLink() string {
	token := a.Token()
	if len(token) < 4 {
		return a.Cross.Link(a.Config.Site_url)
	}
	tk := ""
	if t := a.Token(); t != "" {
		tk = fmt.Sprintf("/%s", t[1:4])
	}
	return fmt.Sprintf("%s%s", a.Cross.Link(a.Config.Site_url), tk)
}

func (a *ProviderArg) Token() string {
	inv := a.Cross.Exfee.FindInvitation(a.To_identity)
	if inv == nil {
		return ""
	}
	return inv.Token
}

func (a *ProviderArg) Timezone() string {
	if a.To_identity.Timezone != "" {
		return a.To_identity.Timezone
	}
	return a.Cross.Time.Begin_at.Timezone
}

func (a *ProviderArg) Confirmed() bool {
	inv := a.Cross.Exfee.FindInvitation(a.To_identity)
	if inv == nil {
		return false
	}
	return inv.IsAccepted()
}

func (a *ProviderArg) ManyPosts() bool {
	return len(a.Posts) >= 11
}

func (a *ProviderArg) LongDescription() bool {
	return len(a.Cross.Description) > 200
}

func (a *ProviderArg) OldAccepted() int {
	acceptedCount := 0
	for _, i := range a.Accepted {
		acceptedCount += 1 + int(i.Mates)
	}
	return a.Cross.TotalAccepted() - acceptedCount
}

func (a *ProviderArg) TextPublicInvitation() (string, error) {
	return getTemplateString("cross_public_invitation.txt", a)
}

func (a *ProviderArg) TextPrivateInvitation() (string, error) {
	return getTemplateString("cross_private_invitation.txt", a)
}

func (a *ProviderArg) TextQuit() (string, error) {
	return getTemplateString("cross_quit.txt", a)
}

func (a *ProviderArg) TextTitleChange() (string, error) {
	return getTemplateString("cross_title_change.txt", a)
}

func (a *ProviderArg) TextCrossChange() (string, error) {
	return getTemplateString("cross_change.txt", a)
}

func (a *ProviderArg) TextAccepted() (string, error) {
	accepted := a.Cross.TotalAccepted()
	otherCount := accepted - len(a.Accepted)

	data := make(map[string]interface{})
	data["Arg"] = a
	data["Invitations"] = a.Accepted
	data["TotalAccepted"] = accepted
	data["OtherCount"] = otherCount
	data["HasOther"] = otherCount > 0
	data["IsOthers"] = otherCount > 1
	return getTemplateString("cross_accepted.txt", data)
}

func (a *ProviderArg) TextDeclined() (string, error) {
	data := make(map[string]interface{})
	data["Arg"] = a
	data["Identities"] = a.Declined
	return getTemplateString("cross_declined.txt", data)
}

func (a *ProviderArg) TextNewlyInvited() (string, error) {
	data := make(map[string]interface{})
	data["Arg"] = a
	data["Invitations"] = a.NewlyInvited
	return getTemplateString("cross_newly_invitations.txt", data)
}

func (a *ProviderArg) TextRemoved() (string, error) {
	data := make(map[string]interface{})
	data["Arg"] = a
	data["Identities"] = a.Removed
	return getTemplateString("cross_removed.txt", data)
}

func (a *ProviderArg) Diff(log *log.Logger) (accepted map[string]*exfe_model.Invitation, declined map[string]*exfe_model.Identity, newlyInvited map[string]*exfe_model.Invitation, removed map[string]*exfe_model.Identity) {
	oldId := make(map[string]*exfe_model.Invitation)
	newId := make(map[string]*exfe_model.Invitation)
	oldExId := make(map[string]*exfe_model.Invitation)
	newExId := make(map[string]*exfe_model.Invitation)

	accepted = make(map[string]*exfe_model.Invitation)
	declined = make(map[string]*exfe_model.Identity)
	newlyInvited = make(map[string]*exfe_model.Invitation)
	removed = make(map[string]*exfe_model.Identity)

	if a.Old_cross == nil {
		return
	}

	for i, v := range a.Old_cross.Exfee.Invitations {
		if v.Rsvp_status == "NOTIFICATION" || v.Rsvp_status == "REMOVED" {
			continue
		}
		if _, ok := oldId[v.Identity.DiffId()]; ok {
			log.Printf("more than one non-notification status in exfee %d, user id %d", a.Old_cross.Id, v.Identity.Connected_user_id)
		}
		oldId[v.Identity.DiffId()] = &a.Old_cross.Exfee.Invitations[i]
		oldExId[v.Identity.ExternalId()] = &a.Old_cross.Exfee.Invitations[i]
	}
	for i, v := range a.Cross.Exfee.Invitations {
		if v.Rsvp_status == "NOTIFICATION" || v.Rsvp_status == "REMOVED" {
			continue
		}
		if _, ok := newId[v.Identity.DiffId()]; ok {
			log.Printf("more than one non-notification status in exfee %d, user id %d", a.Old_cross.Id, v.Identity.Connected_user_id)
		}
		newId[v.Identity.DiffId()] = &a.Cross.Exfee.Invitations[i]
		newExId[v.Identity.ExternalId()] = &a.Cross.Exfee.Invitations[i]
	}

	for k, v := range newId {
		inv, ok := oldId[k]
		if !ok {
			inv, ok = oldExId[v.Identity.ExternalId()]
		}
		switch v.Rsvp_status {
		case "ACCEPTED":
			if !ok || inv.Rsvp_status != v.Rsvp_status {
				accepted[k] = v
			}
		case "DECLINED":
			if !ok || inv.Rsvp_status != v.Rsvp_status {
				declined[k] = &v.Identity
			}
		}
		if !ok {
			newlyInvited[k] = v
		}
	}
	for k, v := range oldId {
		_, ok := newId[k]
		if !ok {
			_, ok = newExId[v.Identity.ExternalId()]
		}
		if !ok {
			removed[k] = &v.Identity
		}
	}

	a.Accepted = make([]*exfe_model.Invitation, 0, 0)
	for _, v := range accepted {
		a.Accepted = append(a.Accepted, v)
	}
	a.Declined = make([]*exfe_model.Identity, 0, 0)
	for _, v := range declined {
		a.Declined = append(a.Declined, v)
	}
	a.NewlyInvited = make([]*exfe_model.Invitation, 0, 0)
	for _, v := range newlyInvited {
		a.NewlyInvited = append(a.NewlyInvited, v)
	}
	a.Removed = make([]*exfe_model.Identity, 0, 0)
	for _, v := range removed {
		a.Removed = append(a.Removed, v)
	}
	return
}
