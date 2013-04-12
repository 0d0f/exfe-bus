package model

import (
	"fmt"
)

type Exfee struct {
	ID          uint64       `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Invitations []Invitation `json:"invitations"`

	Accepted   []Invitation `json:"-"`
	Declined   []Invitation `json:"-"`
	Interested []Invitation `json:"-"`
	Pending    []Invitation `json:"-"`
}

func (e *Exfee) Parse() {
	e.Accepted = make([]Invitation, 0)
	e.Declined = make([]Invitation, 0)
	e.Interested = make([]Invitation, 0)
	e.Pending = make([]Invitation, 0)

	for _, i := range e.Invitations {
		switch i.RsvpStatus {
		case RsvpAccepted:
			e.Accepted = append(e.Accepted, i)
		case RsvpDeclined:
			e.Declined = append(e.Declined, i)
		case RsvpInterested:
			e.Interested = append(e.Interested, i)
		case RsvpNoresponse:
			e.Pending = append(e.Pending, i)
		}
	}
}

func (e Exfee) TotalCount() int {
	return len(e.Invitations)
}

func (e Exfee) AcceptedCount() int {
	ret := 0
	for _, i := range e.Invitations {
		if i.RsvpStatus == RsvpAccepted {
			ret++
		}
	}
	return ret
}

func (e Exfee) FindUser(userId int64) *Invitation {
	for i := range e.Invitations {
		if e.Invitations[i].Identity.UserID == userId {
			return &e.Invitations[i]
		}
	}
	return nil
}

func (e Exfee) CountPeople(invitations []Invitation) int {
	ret := 0
	for _, i := range invitations {
		ret += 1 + int(i.Mates)
	}
	return ret
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
