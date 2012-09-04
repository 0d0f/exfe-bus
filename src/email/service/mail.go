package email_service

import (
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

type EmailSenderService struct {
	SenderDomain string
	Server       string
	Auth         smtp.Auth
}

func (m *EmailSenderService) Check(mail *string, errorString *string) error {
	mail_split := strings.Split(*mail, "@")
	if len(mail_split) != 2 {
		return fmt.Errorf("mail(%s) not valid.", *mail)
	}
	host := mail_split[1]
	mx, err := net.LookupMX(host)
	if err != nil {
		return fmt.Errorf("lookup mail exchange fail: %s", err)
	}
	if len(mx) == 0 {
		return fmt.Errorf("can't find mail exchange for %s", *mail)
	}
	s, err := smtp.Dial(fmt.Sprintf("%s:25", mx[0].Host))
	if err != nil {
		return fmt.Errorf("dial to mail exchange %s fail: %s", mx[0].Host, err)
	}
	err = s.Mail(fmt.Sprintf("x@%s", m.SenderDomain))
	if err != nil {
		return fmt.Errorf("set mail from fail: %s", err)
	}
	err = s.Rcpt(*mail)
	if err != nil {
		*errorString = err.Error()
	}
	return nil
}

func (m *EmailSenderService) EmailSend(arg *MailArg, reply *string) error {
	log.Printf("Send mail: %s", arg)

	var errorString string
	to := make([]*mail.Address, 0, 0)
	for _, t := range arg.To {
		m.Check(&t.Address, &errorString)
		if errorString != "" {
			log.Printf("can't reach mail %s: %s", t.Address, errorString)
		} else {
			to = append(to, t)
		}
	}
	if len(to) == 0 {
		log.Printf("no valid receiver")
		*reply = "no valid receiver"
		return nil
	}
	arg.To = to

	err := arg.SendViaSMTP(m.Server, m.Auth)

	if err != nil {
		log.Printf("Mail send failed: %s", err)
		return err
	}
	log.Print("send ok")

	*reply = "ok"
	return nil
}
