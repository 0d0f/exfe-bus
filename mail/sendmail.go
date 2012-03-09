package main

import (
	"bytes"
	"net/textproto"
	"flag"
	"config"
	"time"
	"gobus"
	"net/smtp"
	"log"
	"fmt"
	"strings"
	"mime/multipart"
)

type MailUser struct {
	Mail string
	Name string
}

func (m MailUser) ToString() string {
	return fmt.Sprintf("\"%s\" <%s>", m.Name, m.Mail)
}

type FilePart struct {
	Name string
	Content string
}

type Mail struct {
	To []MailUser
	From MailUser
	Subject string
	Text string
	Html string
	FileParts []FilePart
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
	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)
	defer w.Close()

	if m.Text != "" {
		header := textproto.MIMEHeader{}
		header.Add("Content-Type", "text/plain; charset=utf-8")
		w1, err := w.CreatePart(header)
		if err != nil {
			log.Printf("Create multipart plain text fail: %s", err.Error())
			return nil
		}
		w1.Write([]byte(m.Text))
	}

	if m.Html != "" {
		header := textproto.MIMEHeader{}
		header.Add("Content-Type", "text/html; charset=utf-8")
		w1, err := w.CreatePart(header)
		if err != nil {
			log.Printf("Create multipart html fail: %s", err.Error())
			return nil
		}
		w1.Write([]byte(m.Html))
	}

	for _, f := range m.FileParts {
		w1, err := w.CreateFormFile(f.Name, f.Name)
		if err != nil {
			log.Printf("Create multipart file(%s) fail: %s", f.Name, err.Error())
			return nil
		}
		w1.Write([]byte(f.Content))
	}
	w.Close()

	return []byte(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\nFrom: %s\r\nSubject: %s\r\nTo: %s\r\n\r\n%s",
		w.Boundary(),
		m.From.ToString(),
		m.Subject,
		m.ToHeader(),
		buf.String()))
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
