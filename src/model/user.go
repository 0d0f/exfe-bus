package model

import (
	"fmt"
	"strings"
)

type UserWelcome struct {
	To         Recipient `json:"to"`
	NeedVerify bool      `json:"need_verify"`

	Config *Config `json:"-"`
}

func (w UserWelcome) String() string {
	return fmt.Sprintf("{to:%s needverify:%v}", w.To.String(), w.NeedVerify)
}

func (w *UserWelcome) Parse(config *Config) (err error) {
	w.Config = config
	return nil
}

func (w UserWelcome) Link() string {
	return fmt.Sprintf("%s/#token=%s", w.Config.SiteUrl, w.To.Token)
}

type UserWelcomes []UserWelcome

func (w UserWelcomes) String() string {
	welcome := make([]string, len(w))
	for i := range w {
		welcome[i] = w[i].String()
	}
	return fmt.Sprintf("[%s]", strings.Join(welcome, ","))
}

type UserVerify struct {
	To       Recipient `json:"to"`
	UserName string    `json:"user_name"`

	Config *Config `json:"-"`
}

func (c UserVerify) String() string {
	return fmt.Sprintf("{to:%s by:%s}", c.To.String(), c.UserName)
}

func (c *UserVerify) Parse(config *Config) (err error) {
	c.Config = config
	return nil
}

func (c UserVerify) Link() string {
	return fmt.Sprintf("%s/#token=%s", c.Config.SiteUrl, c.To.Token)
}

type UserVerifys []UserVerify

func (u UserVerifys) String() string {
	c := make([]string, len(u))
	for i := range u {
		c[i] = u[i].String()
	}
	return fmt.Sprintf("[%s]", strings.Join(c, ","))
}
