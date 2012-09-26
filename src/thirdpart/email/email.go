package email

import (
	"model"
	"thirdpart"
)

type Email struct {
	helper thirdpart.Helper
}

const provider = "email"

func New(helper thirdpart.Helper) *Email {
	return &Email{
		helper: helper,
	}
}

func (e *Email) Provider() string {
	return provider
}

func (e *Email) MessageType() thirdpart.MessageType {
	return thirdpart.EmailMessage
}

func (e *Email) Send(to *model.Recipient, privateMessage string, publicMessage string, info *thirdpart.InfoData) (string, error) {
	return e.helper.SendEmail(to.ExternalID, privateMessage)
}
