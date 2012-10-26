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
