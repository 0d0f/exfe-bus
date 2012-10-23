package args

import (
	"model"
	"thirdpart"
)

type SendArg struct {
	To             *model.Recipient    `json:"to"`
	PrivateMessage string              `json:"private"`
	PublicMessage  string              `json:"public"`
	Info           *thirdpart.InfoData `json:"info"`
}
