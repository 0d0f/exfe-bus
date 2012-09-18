package email

import (
	"fmt"
	"net/stmp"
	"thirdpart"
)

type Email struct {
	Server string
	Auth   smtp.Auth
}

const provider = "email"

func New(host string, port int, username, password string) *Email {
	return &Email{
		Server: fmt.Sprintf("%s:%d", host, port),
		Auth:   smtp.PlainAuth("", username, password, host),
	}
}

func (e *Email) Provider() string {
	return provider
}

func (e *Email) MessageType() MessageType {
	return EmailMessage
}

func (e *Email) UpdateFriends(to *model.Identity) error {
	return fmt.Errorf("email can't update friends")
}

func (e *Email) UpdateIdentity(to *model.Identity) error {
	return fmt.Errorf("email can't update identity")
}

func (e *Email) Send(to *model.Identity, privateMessage string, publicMessage string) error {
	return smtp.SendMail(e.Server, e.Auth, to.ExternalID, privateMessage)
}
