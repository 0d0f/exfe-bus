package business

import (
	"fmt"
	"model"
)

type SummaryArg struct {
	To       *model.Recipient `json:"to"`
	OldCross *model.Cross     `json:"old_cross"`
	Cross    *model.Cross     `json:"cross"`
	Posts    []*model.Post    `json:"posts"`

	Config        *model.Config      `json:"-"`
	NewInvited    []model.Invitation `json:"-"`
	Removed       []model.Invitation `json:"-"`
	NewAccepted   []model.Invitation `json:"-"`
	OldAccepted   []model.Invitation `json:"-"`
	NewDeclined   []model.Invitation `json:"-"`
	NewInterested []model.Invitation `json:"-"`
	NewPending    []model.Invitation `json:"-"`
}

func (a *SummaryArg) Parse() {
	a.NewInvited = make([]model.Invitation, 0)
	a.Removed = make([]model.Invitation, 0)
	a.NewAccepted = make([]model.Invitation, 0)
	a.OldAccepted = make([]model.Invitation, 0)
	a.NewDeclined = make([]model.Invitation, 0)
	a.NewInterested = make([]model.Invitation, 0)
	a.NewPending = make([]model.Invitation, 0)

	a.Cross.Exfee.Parse()
	if !a.IsCrossChange() {
		return
	}

	a.OldCross.Exfee.Parse()
	for _, i := range a.Cross.Exfee.Accepted {
		if !in(&i, a.OldCross.Exfee.Accepted) {
			a.NewAccepted = append(a.NewAccepted, i)
		} else {
			a.OldAccepted = append(a.OldAccepted, i)
		}
	}
	for _, i := range a.Cross.Exfee.Declined {
		if !in(&i, a.OldCross.Exfee.Declined) {
			a.NewDeclined = append(a.NewDeclined, i)
		}
	}
	for _, i := range a.Cross.Exfee.Interested {
		if !in(&i, a.OldCross.Exfee.Interested) {
			a.NewInterested = append(a.NewInterested, i)
		}
	}
	for _, i := range a.Cross.Exfee.Pending {
		if !in(&i, a.OldCross.Exfee.Pending) {
			a.NewPending = append(a.NewPending, i)
		}
	}
	for _, i := range a.Cross.Exfee.Invitations {
		if !in(&i, a.OldCross.Exfee.Invitations) {
			a.NewInvited = append(a.NewInvited, i)
		}
	}
	for _, i := range a.OldCross.Exfee.Invitations {
		if !in(&i, a.Cross.Exfee.Invitations) {
			a.Removed = append(a.Removed, i)
		}
	}
}

func (a *SummaryArg) IsCrossChange() bool {
	return a.OldCross != nil
}

func (a *SummaryArg) IsTitleChange() bool {
	return a.IsCrossChange() && a.OldCross.Title != a.Cross.Title
}

func (a *SummaryArg) IsPlaceChanged() bool {
	if !a.IsCrossChange() {
		return false
	}
	return a.Cross.Place.Same(&a.OldCross.Place)
}

func (a *SummaryArg) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

type Cross struct {
}

func (c *Cross) Summary(arg *SummaryArg) error {
	return nil
}

func in(id *model.Invitation, ids []model.Invitation) bool {
	for _, i := range ids {
		if id.Identity.SameUser(&i.Identity) {
			return true
		}
	}
	return false
}
