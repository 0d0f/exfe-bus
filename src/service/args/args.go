package args

import (
	"fmt"
	"model"
	"thirdpart"
)

type SendArg struct {
	To             *model.Recipient    `json:"to"`
	PrivateMessage string              `json:"private"`
	PublicMessage  string              `json:"public"`
	Info           *thirdpart.InfoData `json:"info"`
}

func (a SendArg) String() string {
	return fmt.Sprintf("{to:%s info:%s}", a.To, a.Info)
}

type ConversationUpdateArg []model.ConversationUpdate

func (u ConversationUpdateArg) String() string {
	if len(u) == 0 {
		return "{updates:0}"
	}
	return fmt.Sprintf("{to:%s with:%s updates}:%d", u[0].To, u[0].Cross, len(u))
}
