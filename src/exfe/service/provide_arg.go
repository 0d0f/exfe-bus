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

func (a *ProviderArg) IsHost() bool {
	return a.Cross.By_identity.Connected_user_id == a.To_identity.Connected_user_id
}

func (a *ProviderArg) Token() string {
	for _, invitation := range a.Cross.Exfee.Invitations {
		if invitation.Identity.Connected_user_id == a.To_identity.Connected_user_id {
			return invitation.Token
		}
	}
	return ""
}

func (a *ProviderArg) Timezone() string {
	if a.To_identity.Timezone != "" {
		return a.To_identity.Timezone
	}
	return a.Cross.Time.Begin_at.Timezone
}

func (a *ProviderArg) Confirmed() bool {
	for _, invitation := range a.Cross.Exfee.Invitations {
		if invitation.Identity.Connected_user_id == a.To_identity.Connected_user_id {
			return invitation.IsAccepted()
		}
	}
	return false
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
	data["Identities"] = a.Accepted
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

	accepted = make(map[string]*exfe_model.Invitation)
	declined = make(map[string]*exfe_model.Identity)
	newlyInvited = make(map[string]*exfe_model.Invitation)
	removed = make(map[string]*exfe_model.Identity)

	if a.Old_cross == nil {
		return
	}

	for i, v := range a.Old_cross.Exfee.Invitations {
		if v.Rsvp_status == "NOTIFICATION" {
			continue
		}
		if _, ok := oldId[v.Identity.DiffId()]; ok {
			log.Printf("more than one non-notification status in exfee %d, user id %d", a.Old_cross.Id, v.Identity.Connected_user_id)
		}
		oldId[v.Identity.DiffId()] = &a.Old_cross.Exfee.Invitations[i]
	}
	for i, v := range a.Cross.Exfee.Invitations {
		if v.Rsvp_status == "NOTIFICATION" {
			continue
		}
		if _, ok := newId[v.Identity.DiffId()]; ok {
			log.Printf("more than one non-notification status in exfee %d, user id %d", a.Old_cross.Id, v.Identity.Connected_user_id)
		}
		newId[v.Identity.DiffId()] = &a.Cross.Exfee.Invitations[i]
	}

	for k, v := range newId {
		inv, ok := oldId[k]
		if !ok {
			inv, ok = oldId[v.Identity.ExternalId()]
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
			_, ok = newId[v.Identity.UserId()]
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
