package notifier

import (
	"broker"
	"fmt"
	"formatter"
	"model"
)

type Exfee struct {
	localTemplate *formatter.LocalTemplate
	config        *model.Config
	platform      *broker.Platform
}

func NewExfee(localTemplate *formatter.LocalTemplate, config *model.Config, platform *broker.Platform) *Exfee {
	return &Exfee{
		localTemplate: localTemplate,
		config:        config,
		platform:      platform,
	}
}

func (e Exfee) V3Rsvp(updates []model.RsvpUpdate) error {
	arg, err := RsvpFromUpdates(updates, e.config)
	if err != nil {
		return err
	}

	to := arg.To
	text, err := GenerateContent(e.localTemplate, "v3_exfee_rsvp", to.Provider, to.Language, arg)
	if err != nil {
		return err
	}
	_, err = e.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

func (e Exfee) V3Conversation(updates []model.ConversationUpdate) error {
	arg, err := ArgFromUpdates(updates, e.config)
	if err != nil {
		return err
	}

	to := arg.To
	text, err := GenerateContent(e.localTemplate, "v3_exfee_conversation", to.Provider, to.Language, arg)
	if err != nil {
		return err
	}
	_, err = e.platform.Send(to, text)
	if err != nil {
		return err
	}
	return nil
}

func in(id *model.Invitation, ids []model.Invitation) bool {
	for _, i := range ids {
		if id.Identity.SameUser(i.Identity) {
			return true
		}
	}
	return false
}

type RsvpArg struct {
	model.ThirdpartTo
	OldExfee *model.Exfee
	Exfee    model.Exfee
	Bys      []model.Identity

	NewInvited  []model.Identity
	Removed     []model.Identity
	NewAccepted []model.Identity
	OldAccepted []model.Identity
	NewDeclined []model.Identity
}

func RsvpFromUpdates(updates []model.RsvpUpdate, config *model.Config) (*RsvpArg, error) {
	if updates == nil && len(updates) == 0 {
		return nil, fmt.Errorf("no update info")
	}

	to := updates[0].To
	bys := make([]model.Identity, 0)

Bys:
	for _, update := range updates {
		if !to.Equal(&update.To) {
			return nil, fmt.Errorf("updates not send to same recipient: %s, %s", to, update.To)
		}
		for _, i := range bys {
			if update.By.SameUser(i) {
				continue Bys
			}
		}
		bys = append(bys, update.By)
	}

	ret := &RsvpArg{
		Bys:      bys,
		OldExfee: &updates[0].OldExfee,
		Exfee:    updates[len(updates)-1].Exfee,

		NewInvited:  make([]model.Identity, 0),
		Removed:     make([]model.Identity, 0),
		NewAccepted: make([]model.Identity, 0),
		OldAccepted: make([]model.Identity, 0),
		NewDeclined: make([]model.Identity, 0),
	}
	ret.To = to
	err := ret.Parse(config)
	if err != nil {
		return nil, err
	}

	ret.Exfee.Parse()
	ret.OldExfee.Parse()

	for _, i := range ret.Exfee.Accepted {
		if !in(&i, ret.OldExfee.Accepted) {
			ret.NewAccepted = append(ret.NewAccepted, i.Identity)
		} else {
			ret.OldAccepted = append(ret.OldAccepted, i.Identity)
		}
	}
	for _, i := range ret.Exfee.Declined {
		if !in(&i, ret.OldExfee.Declined) {
			ret.NewDeclined = append(ret.NewDeclined, i.Identity)
		}
	}
	for _, i := range ret.Exfee.Invitations {
		if !in(&i, ret.OldExfee.Invitations) {
			ret.NewInvited = append(ret.NewInvited, i.Identity)
		}
	}
	for _, i := range ret.OldExfee.Invitations {
		if !in(&i, ret.Exfee.Invitations) {
			ret.Removed = append(ret.Removed, i.Identity)
		}
	}
	return ret, nil
}

type ConversationArg struct {
	model.ThirdpartTo
	Posts []*model.Post
}

func ArgFromConversations(updates []model.ConversationUpdate, config *model.Config) (*ConversationArg, error) {
	if updates == nil && len(updates) == 0 {
		return nil, fmt.Errorf("no update info")
	}

	to := updates[0].To
	posts := make([]*model.Post, len(updates))

	for i, update := range updates {
		if !to.Equal(&update.To) {
			return nil, fmt.Errorf("updates not send to same recipient: %s, %s", to, update.To)
		}
		posts[i] = &updates[i].Post
	}

	ret := &ConversationArg{
		Posts: posts,
	}
	ret.To = to
	err := ret.Parse(config)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a ConversationArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a ConversationArg) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return "+00:00"
}
