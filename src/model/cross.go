package model

import (
	"fmt"
)

type Cross struct {
	ID          uint64    `json:"id,omitempty"`
	By          Identity  `json:"by_identity,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Time        CrossTime `json:"time,omitempty"`
	Place       Place     `json:"place,omitempty"`
	Exfee       Exfee     `json:"exfee,omitempty"`
}

func (c Cross) Equal(other *Cross) bool {
	return c.ID == other.ID
}

func (c Cross) String() string {
	return fmt.Sprintf("Cross:%d", c.ID)
}

type CrossUpdate struct {
	To       Recipient `json:"to"`
	OldCross Cross     `json:"old_cross"`
	Cross    Cross     `json:"cross"`
	By       Identity  `json:"by"`
}

type CrossUpdates []CrossUpdate

func (u CrossUpdates) String() string {
	if len(u) == 0 {
		return "{updates:0}"
	}
	return fmt.Sprintf("{to:%s with:%s updates:%d}", u[0].To, u[0].Cross, len(u))
}

type CrossInvitation struct {
	To    Recipient `json:"to"`
	Cross Cross     `json:"cross"`

	Config *Config `json:"-"`
}

func (a CrossInvitation) String() string {
	return fmt.Sprintf("{to:%s cross:%d}", a.To, a.Cross.ID)
}

func (a *CrossInvitation) Parse(config *Config) (err error) {
	a.Config = config
	return nil
}

func (a CrossInvitation) ToIn(invitations []Invitation) bool {
	for _, i := range invitations {
		if a.To.SameUser(&i.Identity) {
			return true
		}
	}
	return false
}

func (a CrossInvitation) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a CrossInvitation) PublicLink() string {
	return fmt.Sprintf("%s/#!%d/%s", a.Config.SiteUrl, a.Cross.ID, a.To.Token[1:5])
}

func (a CrossInvitation) Timezone() string {
	if a.To.Timezone != "" {
		return a.To.Timezone
	}
	return a.Cross.Time.BeginAt.Timezone
}

func (a CrossInvitation) IsCreator() bool {
	return a.To.SameUser(&a.Cross.By)
}

func (a CrossInvitation) LongDescription() bool {
	if len(a.Cross.Description) > 200 {
		return true
	}
	return false
}

func (a CrossInvitation) ListInvitations() string {
	l := len(a.Cross.Exfee.Invitations)
	max := 3
	ret := ""
	for i := 0; i < 3 && i < l; i++ {
		if i > 0 {
			ret += ", "
		}
		ret += a.Cross.Exfee.Invitations[i].Identity.Name
	}
	if l > max {
		ret += "..."
	}
	return ret
}
