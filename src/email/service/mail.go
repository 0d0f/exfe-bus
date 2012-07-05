package email_service

import (
	"net/smtp"
	"log"
)

type EmailSenderService struct {
	Server string
	Auth   smtp.Auth
}

func (m *EmailSenderService) EmailSend(arg *MailArg, reply *string) error {
	log.Printf("Send mail: %s", arg)
	err := arg.SendViaSMTP(m.Server, m.Auth)

	if err != nil {
		log.Printf("Mail send failed: %s", err)
		return err
	}
	log.Print("send ok")

	*reply = "ok"
	return nil
}
