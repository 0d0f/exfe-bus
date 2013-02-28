package model

import (
	"fmt"
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

func (a ThirdpartTo) ToRecipient() Recipient {
	return a.To
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
	To   Recipient `json:"to"`
	Text string    `json:"text"`

	Config *Config `json:"-"`
}

func (a ThirdpartSend) String() string {
	text := a.Text
	if len(text) > 10 {
		text = text[:10]
	}
	return fmt.Sprintf("{to:%s text:%s}", a.To, text)
}
