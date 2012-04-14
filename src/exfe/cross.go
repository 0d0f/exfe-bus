package exfe

import (
	"fmt"
)

type Cross struct {
	Id uint64
	Id_base62 string
	Type string
	Title string
	Description string
	Time CrossTime
	Place Place
	Attribute map[string]string
	Exfee []Invitation
	Widget []interface{}
	Relative []struct {
		Id uint64
		Relation string
	}
	By_identity Identity
}

func (c *Cross) Link(host string) string {
	return fmt.Sprintf("%s/!%s", host, c.Id_base62)
}

func (c *Cross) LinkTo(host string, target *Invitation) string {
	return fmt.Sprintf("%s?token=%s", c.Link(host), target.Token)
}
