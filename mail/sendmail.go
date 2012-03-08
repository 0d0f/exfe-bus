package main

import (
	"flag"
	"config"
	"time"
	"gobus"
	"net/smtp"
	"log"
	"fmt"
	"strings"
)

type MailUser struct {
	Mail string
	Name string
}

func (m MailUser) ToString() string {
	return fmt.Sprintf("\"%s\" <%s>", m.Name, m.Mail)
}

type Mail struct {
	To []MailUser
	From MailUser
	Subject string
	Message string
}

func (m *Mail) GoString() string {
	return fmt.Sprintf("Mail send from %s to %s with subject: %s", m.From.ToString(), m.ToLine(), m.Subject)
}

func (m *Mail) ToLine() string {
	var users []string
	for _, m := range m.To {
		users = append(users, m.ToString())
	}
	return strings.Join(users, ", ")
}

func (m *Mail) ToHeader() string {
	var users []string
	for _, m := range m.To {
		users = append(users, m.ToString())
	}
	return strings.Join(users, ", \r\n        ")
}

func (m *Mail) ToMail() (mails []string) {
	for _, m := range m.To {
		mails = append(mails, m.Mail)
	}
	return
}

func (m *Mail) Body() []byte {
	return []byte(fmt.Sprintf("From: %s\r\nSubject: %s\r\nTo: %s\r\n\r\n%s\r\n", m.From.ToString(), m.Subject, m.ToHeader(), m.Message))
}

type MailSender struct {
	server string
	auth smtp.Auth
}

func (m *MailSender) Do(messages []interface{}) []interface{} {
	mail, ok := messages[0].(*Mail)
	if !ok {
		log.Printf("Can't convert input into Mail: %s", messages)
	}

	log.Printf("Try to send subject(%s) mail from (%s) to (%s)...", mail.Subject, mail.From.ToString(), mail.ToLine())
	err := smtp.SendMail(m.server, m.auth, mail.From.Mail, mail.ToMail(), mail.Body())

	errString := ""
	if err != nil {
		log.Printf("Mail send failed: %s", err)
		errString = err.Error()
	}

	return []interface{}{&errString}
}

func (m *MailSender) MaxJobsCount() int {
	return 1
}

func (m *MailSender) JobGenerator() interface{} {
	return &Mail{}
}

const (
	queue = "gobus:queue:mail:sender"
)

func main() {
	log.SetPrefix("[Sendmail]")
	log.Printf("Service start, queue: %s", queue)

	var configFile string
	flag.StringVar(&configFile, "config", "mail.yaml", "Specify the configuration file")
	flag.Parse()

	config := config.LoadFile(configFile)

	service := gobus.CreateService(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		queue,
		&MailSender{
			server: fmt.Sprintf("%s:%d", config.String("mail.host"), config.Int("mail.port")),
			auth: smtp.PlainAuth("", config.String("mail.user"), config.String("mail.password"), config.String("mail.host")),
		},
		config.Int("service.limit"))
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Run(time.Duration(config.Int("service.time_out")))
}
