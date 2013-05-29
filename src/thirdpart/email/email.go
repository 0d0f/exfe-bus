package email

import (
	"model"
	"thirdpart"
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

func (e *Email) Send(to *model.Recipient, text string) (string, error) {
	return e.Post(to.ExternalUsername, text)
}

func (e *Email) Post(id, text string) (string, error) {
	return e.helper.SendEmail(id, text)
}
