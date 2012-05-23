package mail

import (
	"net/smtp"
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/textproto"
	"strings"
)

type MailUser struct {
	Mail string
	Name string
}

func (m MailUser) ToString() string {
	return fmt.Sprintf("\"%s\" <%s>", m.Name, m.Mail)
}

type FilePart struct {
	Name    string
	Content []byte
}

type Mail struct {
	To        []MailUser
	From      MailUser
	Subject   string
	Text      string
	Html      string
	FileParts []FilePart
}

func (m *Mail) String() string {
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

func (m *Mail) MakeMessage() (textproto.MIMEHeader, string) {
	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)
	defer w.Close()

	header := textproto.MIMEHeader{}
	header.Add("Content-Type", "text/plain; charset=utf-8")
	w1, err := w.CreatePart(header)
	if err != nil {
		log.Printf("Create multipart plain text fail: %s", err.Error())
		return nil, ""
	}
	w1.Write([]byte(m.Text))

	header = textproto.MIMEHeader{}
	header.Add("Content-Type", "text/html; charset=utf-8")
	w1, err = w.CreatePart(header)
	if err != nil {
		log.Printf("Create multipart html fail: %s", err.Error())
		return nil, ""
	}
	w1.Write([]byte(m.Html))

	w.Close()
	header = textproto.MIMEHeader{}
	header.Add("Content-Type", fmt.Sprintf("multipart/alternative; boundary=\"%s\"", w.Boundary()))

	return header, buf.String()
}

func (m *Mail) MakeMessageWithAttachments() (textproto.MIMEHeader, string) {
	header, message := m.MakeMessage()
	if header == nil {
		log.Printf("Can't create message part")
		return nil, ""
	}
	if len(m.FileParts) == 0 {
		return header, message
	}

	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)
	defer w.Close()

	messagePart, err := w.CreatePart(header)
	if err != nil {
		w.Close()
		log.Printf("Create multipart message part fail: %s", err.Error())
		return nil, ""
	}
	messagePart.Write([]byte(message))

	for _, f := range m.FileParts {
		w1, err := w.CreateFormFile(f.Name, f.Name)
		if err != nil {
			log.Printf("Create multipart file(%s) fail: %s", f.Name, err.Error())
			return nil, ""
		}
		w1.Write([]byte(f.Content))
	}
	w.Close()

	header = textproto.MIMEHeader{}
	header.Add("Content-Type", fmt.Sprintf("multipart/mixed; boundary=\"%s\"", w.Boundary()))

	return header, buf.String()
}

func (m *Mail) Body() []byte {
	header := textproto.MIMEHeader{}
	body := ""
	if len(m.FileParts) == 0 {
		header, body = m.MakeMessage()
	} else {
		header, body = m.MakeMessageWithAttachments()
	}

	return []byte(fmt.Sprintf("Content-Type: %s\r\nFrom: %s\r\nSubject: %s\r\nTo: %s\r\n%s",
		header.Get("Content-Type"),
		m.From.ToString(),
		m.Subject,
		m.ToHeader(),
		body))
}

func (m *Mail) SendSMTP(server string, auth smtp.Auth) error {
	return smtp.SendMail(server, auth, m.From.Mail, m.ToMail(), m.Body())
}
