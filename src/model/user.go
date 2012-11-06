package model

import (
	"fmt"
	"strings"
)

type UserWelcome struct {
	ThirdpartTo
	NeedVerify bool `json:"need_verify"`
}

func (w UserWelcome) String() string {
	return fmt.Sprintf("{to:%s needverify:%v}", w.ThirdpartTo.String(), w.NeedVerify)
}

type UserWelcomes []UserWelcome

func (w UserWelcomes) String() string {
	welcome := make([]string, len(w))
	for i := range w {
		welcome[i] = w[i].String()
	}
	return fmt.Sprintf("[%s]", strings.Join(welcome, ","))
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

type UserConfirms []UserConfirm

func (u UserConfirms) String() string {
	c := make([]string, len(u))
	for i := range u {
		c[i] = u[i].String()
	}
	return fmt.Sprintf("[%s]", strings.Join(c, ","))
}
