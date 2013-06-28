package model

import (
	"fmt"
)

type RsvpUpdate struct {
	To       Recipient `json:"to"`
	By       Identity  `json:"by"`
	Exfee    Exfee     `json:"exfee"`
	OldExfee Exfee     `json:"old_exfee"`
}

type Exfee struct {
	ID            int64        `json:"id,omitempty"`
	Name          string       `json:"name,omitempty"`
	Invitations   []Invitation `json:"invitations,omitempty"`
	ItemsCount    int          `json:"items,omitempty"`
	TotalCount    int          `json:"total,omitempty"`
	AcceptedCount int          `json:"accepted,omitempty"`

	Accepted []Invitation `json:"-"`
	Declined []Invitation `json:"-"`
	Pending  []Invitation `json:"-"`
}

func (e *Exfee) Parse() {
	e.Accepted = make([]Invitation, 0)
	e.Declined = make([]Invitation, 0)
	e.Pending = make([]Invitation, 0)

	for _, i := range e.Invitations {
		switch i.Response {
		case Accepted:
			e.Accepted = append(e.Accepted, i)
		case Declined:
			e.Declined = append(e.Declined, i)
		default:
			e.Pending = append(e.Pending, i)
		}
	}
}

func (e Exfee) FindUser(userId int64) *Invitation {
	for i := range e.Invitations {
		if e.Invitations[i].Identity.UserID == userId {
			return &e.Invitations[i]
		}
	}
	return nil
}

func (e Exfee) Equal(other *Exfee) bool {
	return e.ID == other.ID
}

func (e Exfee) FindInvitedUser(identity Identity) (Invitation, error) {
	for _, inv := range e.Invitations {
		if inv.Identity.SameUser(identity) {
			return inv, nil
		}
	}
	for _, inv := range e.Invitations {
		if inv.Identity.ExternalUsername == identity.ExternalUsername {
			return inv, nil
		}
	}
	return Invitation{}, fmt.Errorf("can't find %s", identity)
}
