package business

import (
	"bytes"
	"fmt"
	"formater"
	"model"
	"thirdpart"
)

type SummaryArg struct {
	To       *model.Recipient `json:"-"`
	OldCross *model.Cross     `json:"-"`
	Cross    *model.Cross     `json:"-"`
	Bys      []model.Identity `json:"-"`

	Config        *model.Config      `json:"-"`
	NewInvited    []model.Invitation `json:"-"`
	Removed       []model.Invitation `json:"-"`
	NewAccepted   []model.Invitation `json:"-"`
	OldAccepted   []model.Invitation `json:"-"`
	NewDeclined   []model.Invitation `json:"-"`
	NewInterested []model.Invitation `json:"-"`
	NewPending    []model.Invitation `json:"-"`
}

func SummaryFromUpdates(updates []model.Update, config *model.Config) (*SummaryArg, error) {
	if updates == nil && len(updates) == 0 {
		return nil, fmt.Errorf("no update info")
	}

	to := &updates[0].To
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

	ret := &SummaryArg{
		To:       to,
		Bys:      bys,
		OldCross: &updates[0].OldCross,
		Cross:    &updates[len(updates)-1].Cross,

		Config:        config,
		NewInvited:    make([]model.Invitation, 0),
		Removed:       make([]model.Invitation, 0),
		NewAccepted:   make([]model.Invitation, 0),
		OldAccepted:   make([]model.Invitation, 0),
		NewDeclined:   make([]model.Invitation, 0),
		NewInterested: make([]model.Invitation, 0),
		NewPending:    make([]model.Invitation, 0),
	}

	ret.Cross.Exfee.Parse()
	ret.OldCross.Exfee.Parse()

	for _, i := range ret.Cross.Exfee.Accepted {
		if !in(&i, ret.OldCross.Exfee.Accepted) {
			ret.NewAccepted = append(ret.NewAccepted, i)
		} else {
			ret.OldAccepted = append(ret.OldAccepted, i)
		}
	}
	for _, i := range ret.Cross.Exfee.Declined {
		if !in(&i, ret.OldCross.Exfee.Declined) {
			ret.NewDeclined = append(ret.NewDeclined, i)
		}
	}
	for _, i := range ret.Cross.Exfee.Interested {
		if !in(&i, ret.OldCross.Exfee.Interested) {
			ret.NewInterested = append(ret.NewInterested, i)
		}
	}
	for _, i := range ret.Cross.Exfee.Pending {
		if !in(&i, ret.OldCross.Exfee.Pending) {
			ret.NewPending = append(ret.NewPending, i)
		}
	}
	for _, i := range ret.Cross.Exfee.Invitations {
		if !in(&i, ret.OldCross.Exfee.Invitations) {
			ret.NewInvited = append(ret.NewInvited, i)
		}
	}
	for _, i := range ret.OldCross.Exfee.Invitations {
		if !in(&i, ret.Cross.Exfee.Invitations) {
			ret.Removed = append(ret.Removed, i)
		}
	}
	fmt.Printf("%+v\n", ret)
	return ret, nil
}

func (a *SummaryArg) TotalOldAccepted() int {
	ret := 0
	for _, e := range a.OldAccepted {
		ret += 1 + int(e.Mates)
	}
	return ret
}

func (a *SummaryArg) IsTimeChanged() bool {
	oldtime, _ := a.OldCross.Time.StringInZone(a.To.Timezone)
	time, _ := a.Cross.Time.StringInZone(a.To.Timezone)
	return oldtime != time
}

func (a *SummaryArg) IsTitleChanged() bool {
	return a.OldCross.Title != a.Cross.Title
}

func (a *SummaryArg) IsPlaceChanged() bool {
	return !a.Cross.Place.Same(&a.OldCross.Place)
}

func (a *SummaryArg) ListBy(limit int, join string) string {
	buf := bytes.NewBuffer(nil)
	for i, by := range a.Bys {
		if buf.Len() > 0 {
			buf.WriteString(join)
		}
		if i >= limit {
			buf.WriteString("etc")
			break
		}
		buf.WriteString(by.Name)
	}
	return buf.String()
}

func (a *SummaryArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

type Cross struct {
	localTemplate *formater.LocalTemplate
	config        *model.Config
}

func NewCross(localTemplate *formater.LocalTemplate, config *model.Config) *Cross {
	return &Cross{
		localTemplate: localTemplate,
		config:        config,
	}
}

func (c *Cross) Summary(updates []model.Update) error {
	_, err := c.getContent(updates)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cross) getContent(updates []model.Update) (string, error) {
	arg, err := SummaryFromUpdates(updates, c.config)
	if err != nil {
		return "", fmt.Errorf("can't parse update: %s", err)
	}
	messageType, err := thirdpart.MessageTypeFromProvider(arg.To.Provider)
	if err != nil {
		return "", err
	}
	templateName := fmt.Sprintf("cross_summary.%s", messageType)
	buf := bytes.NewBuffer(nil)
	err = c.localTemplate.Execute(buf, arg.To.Language, templateName, arg)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func in(id *model.Invitation, ids []model.Invitation) bool {
	for _, i := range ids {
		if id.Identity.SameUser(i.Identity) {
			return true
		}
	}
	return false
}