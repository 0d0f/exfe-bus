package email

import (
	"fmt"
	"model"
	"net"
	"net/smtp"
	"strings"
	"thirdpart"
)

type Email struct {
	Server string
	Auth   smtp.Auth
	From   string
}

const provider = "email"

func New(host string, port int, username, password, from string) *Email {
	return &Email{
		Server: fmt.Sprintf("%s:%d", host, port),
		Auth:   smtp.PlainAuth("", username, password, host),
		From:   from,
	}
}

func (e *Email) Provider() string {
	return provider
}

func (e *Email) MessageType() thirdpart.MessageType {
	return thirdpart.EmailMessage
}

func (e *Email) Send(to *model.Recipient, privateMessage string, publicMessage string, info *thirdpart.InfoData) (string, error) {
	err := e.check(to.ExternalID)
	if err != nil {
		return "", err
	}
	return "", smtp.SendMail(e.Server, e.Auth, e.From, []string{to.ExternalID}, []byte(privateMessage))
}

func (e *Email) check(mail string) error {
	mail_split := strings.Split(mail, "@")
	if len(mail_split) != 2 {
		return fmt.Errorf("mail(%s) not valid.", mail)
	}
	host := mail_split[1]
	mx, err := net.LookupMX(host)
	if err != nil {
		return fmt.Errorf("lookup mail exchange fail: %s", err)
	}
	if len(mx) == 0 {
		return fmt.Errorf("can't find mail exchange for %s", mail)
	}
	s, err := smtp.Dial(fmt.Sprintf("%s:25", mx[0].Host))
	if err != nil {
		return fmt.Errorf("dial to mail exchange %s fail: %s", mx[0].Host, err)
	}
	err = s.Mail(e.From)
	if err != nil {
		return fmt.Errorf("set smtp mail fail: %s", err)
	}
	err = s.Rcpt(mail)
	if err != nil {
		return fmt.Errorf("set smtp rcpt fail: %s", err)
	}
	return nil
}
