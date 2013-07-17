package email

import (
	"thirdpart"
	"time"
)

type Email struct {
	helper thirdpart.Helper
}

func New(helper thirdpart.Helper) *Email {
	return &Email{
		helper: helper,
	}
}

func (e *Email) Provider() string {
	return "email"
}

func (e *Email) SetPosterCallback(callback thirdpart.Callback) (time.Duration, bool) {
	return time.Hour * 72, true
}

func (e *Email) Post(from, id, text string) (string, error) {
	return e.helper.SendEmail(id, text)
}
