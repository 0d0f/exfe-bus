package exfe_model

import (
	"fmt"
)

type Cross struct {
	Id uint64
	Title string
	Description string
	Time CrossTime
	Place Place
	Attribute map[string]string
	Exfee Exfee
	Widget []interface{}
	Relative []struct {
		Id uint64
		Relation string
	}
	Type string
	By_identity Identity
}

func (c *Cross) Link(host string) string {
	return fmt.Sprintf("%s/!%s", host, c.Id)
}

func (c *Cross) LinkTo(host string, token string) string {
	return fmt.Sprintf("%s?token=%s", c.Link(host), token)
}
