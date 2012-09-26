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

func (e *Email) MessageType() thirdpart.MessageType {
	return thirdpart.EmailMessage
}

func (e *Email) Send(to *model.Recipient, privateMessage string, publicMessage string, info *thirdpart.InfoData) (string, error) {
	return e.helper.SendEmail(to.ExternalID, privateMessage)
}
