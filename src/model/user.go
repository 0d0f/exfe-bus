package model

import (
	"fmt"
)

type UserWelcome struct {
	ThirdpartTo
	NeedVerify bool `json:"need_verify"`
}

func (w UserWelcome) String() string {
	return fmt.Sprintf("{to:%s needverify:%v}", w.ThirdpartTo.String(), w.NeedVerify)
}

type UserConfirm struct {
	ThirdpartTo
	By Identity `json:"by"`
}

func (c UserConfirm) String() string {
	return fmt.Sprintf("{to:%s by:%s}", c.ThirdpartTo.String(), c.By.String())
}

func (a UserConfirm) NeedShowBy() bool {
	return !a.To.SameUser(&a.By)
}
