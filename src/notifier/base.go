package notifier

import (
	"fmt"
	"model"
)

type ArgBase struct {
	To     model.Recipient `json:"to"`
	Config *model.Config   `json:"-"`
}

func (a ArgBase) ToIn(invitations []model.Invitation) bool {
	for _, i := range invitations {
		if a.To.SameUser(&i.Identity) {
			return true
		}
	}
	return false
}

func (a ArgBase) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}
