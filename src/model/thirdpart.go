package model

import (
	"fmt"
	"strings"
)

type ThirdpartTo struct {
	To Recipient `json:"to"`

	Config *Config `json:"-"`
}

func (a ThirdpartTo) String() string {
	return a.To.String()
}

func (a *ThirdpartTo) Parse(config *Config) (err error) {
	a.Config = config
	return nil
}

func (a ThirdpartTo) ToIn(invitations []Invitation) bool {
	for _, i := range invitations {
		if a.To.SameUser(&i.Identity) {
			return true
		}
	}
	return false
}

func (a ThirdpartTo) Link() string {
	return fmt.Sprintf("%s/#!token=%s", a.Config.SiteUrl, a.To.Token)
}

func (a ThirdpartTo) ToRecipient() Recipient {
	return a.To
}

type ThirdpartTos []ThirdpartTo

func (t ThirdpartTos) String() string {
	c := make([]string, len(t))
	for i := range t {
		c[i] = t[i].String()
	}
	return fmt.Sprintf("[%s]", strings.Join(c, ","))
}

type DataType string

func (t DataType) String() string {
	return string(t)
}

const (
	TypeCrossInvitation DataType = "i"
	TypeCrossUpdate              = "u"
	TypeCrossRemove              = "r"
	TypeConversation             = "c"
)

type InfoData struct {
	CrossID uint64   `json:"cross_id"`
	Type    DataType `json:"type"`
}

func (i InfoData) String() string {
	return fmt.Sprintf("{cross:%d type:%s}", i.CrossID, i.Type)
}

type ThirdpartSend struct {
	To             Recipient `json:"to"`
	PrivateMessage string    `json:"private"`
	PublicMessage  string    `json:"public"`
	Info           *InfoData `json:"info"`

	Config *Config `json:"-"`
}

func (a ThirdpartSend) String() string {
	return fmt.Sprintf("{to:%s info:%s}", a.To, a.Info)
}
