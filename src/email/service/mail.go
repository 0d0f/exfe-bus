package email_service

import (
	"net/smtp"
	"log"
	"gomail"
)

type EmailSenderService struct {
	Server string
	Auth   smtp.Auth
}

func (m *EmailSenderService) EmailSend(arg *gomail.Mail, reply *string) error {
	log.Printf("Send mail: subject(%s) from (%s) to (%s)", arg.Subject, arg.From.ToString(), arg.ToLine())
	to := arg.ToMail()
	body := arg.Body()
	err := smtp.SendMail(m.Server, m.Auth, arg.From.Mail, to, body)

	if err != nil {
		log.Printf("Mail send failed: %s", err)
		return err
	}
	log.Print("send ok")

	*reply = "ok"
	return nil
}
