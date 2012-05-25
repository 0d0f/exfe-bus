package email_service

import (
	"net/smtp"
	"log/syslog"
	"gomail"
	"fmt"
)

type EmailSenderService struct {
	Server string
	Auth   smtp.Auth
	Log    *syslog.Writer
}

func (m *EmailSenderService) EmailSend(arg *gomail.Mail, reply *string) error {
	m.Log.Info(fmt.Sprintf("Send mail: subject(%s) from (%s) to (%s)", arg.Subject, arg.From.ToString(), arg.ToLine()))
	to := arg.ToMail()
	body := arg.Body()
	err := smtp.SendMail(m.Server, m.Auth, arg.From.Mail, to, body)

	if err != nil {
		m.Log.Err(fmt.Sprintf("Mail send failed: %s", err))
		return err
	}
	m.Log.Info("send ok")

	*reply = "ok"
	return nil
}
