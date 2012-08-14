package email_service

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
)

type FilePart struct {
	Name    string
	Content []byte
}

type MailArg struct {
	To        []*mail.Address
	From      *mail.Address
	Subject   string
	Header    textproto.MIMEHeader
	Text      string
	Html      string
	FileParts []FilePart
}

func (m *MailArg) String() string {
	return fmt.Sprintf("Mail send from %s to %s with subject: %s", m.From, m.To, m.Subject)
}

func (m *MailArg) makeMessage() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)
	defer w.Close()

	header := make(textproto.MIMEHeader)
	header.Add("Content-Type", "text/plain; charset=utf-8")
	w1, err := w.CreatePart(header)
	if err != nil {
		return nil, fmt.Errorf("Create multipart plain text fail: %s", err)
	}
	w1.Write([]byte(m.Text))

	header = make(textproto.MIMEHeader)
	header.Add("Content-Type", "text/html; charset=utf-8")
	w1, err = w.CreatePart(header)
	if err != nil {
		return nil, fmt.Errorf("Create multipart html fail: %s", err)
	}
	w1.Write([]byte(m.Html))

	w.Close()
	m.Header.Add("Content-Type", fmt.Sprintf("multipart/alternative; boundary=\"%s\"", w.Boundary()))

	return buf.Bytes(), nil
}

func (m *MailArg) makeMessageWithAttachments() ([]byte, error) {
	message, err := m.makeMessage()
	if err != nil {
		return nil, fmt.Errorf("Can't create message part")
	}
	if len(m.FileParts) == 0 {
		return message, nil
	}

	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)
	defer w.Close()

	header := make(textproto.MIMEHeader)
	header.Add("Content-Type", m.Header.Get("Content-Type"))
	messagePart, err := w.CreatePart(header)
	if err != nil {
		return nil, fmt.Errorf("Create multipart message part fail: %s", err)
	}
	messagePart.Write(message)

	for _, f := range m.FileParts {
		w1, err := w.CreateFormFile(f.Name, f.Name)
		if err != nil {
			return nil, fmt.Errorf("Create multipart file(%s) fail: %s", f.Name, err)
		}
		w1.Write([]byte(f.Content))
	}
	w.Close()

	m.Header["Content-Type"] = []string{fmt.Sprintf("multipart/mixed; boundary=\"%s\"", w.Boundary())}

	return buf.Bytes(), nil
}

func (m *MailArg) makeHeader() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for k, v := range m.Header {
		buf.Write([]byte(k))
		buf.WriteString(": ")
		buf.Write([]byte(strings.Join(v, ", ")))
		buf.WriteString("\r\n")
	}
	buf.WriteString("To: ")
	for i, t := range m.To {
		buf.WriteString(t.String())
		if i != (len(m.To) - 1) {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("\r\n")
	buf.WriteString(fmt.Sprintf("From: %s\r\n", m.From))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
	return buf.Bytes(), nil
}

func (m *MailArg) makeContent() ([]byte, error) {
	// make message first to get boundary
	var err error
	var body []byte
	if len(m.FileParts) == 0 {
		body, err = m.makeMessage()
	} else {
		body, err = m.makeMessageWithAttachments()
	}
	if err != nil {
		return nil, err
	}

	header, err := m.makeHeader()
	if err != nil {
		return nil, err
	}

	ret := header
	ret = append(ret, []byte("\r\n")...)
	ret = append(ret, body...)
	return ret, nil
}

func (m *MailArg) SendViaSMTP(server string, auth smtp.Auth) error {
	if m.Header == nil {
		m.Header = make(textproto.MIMEHeader)
	}
	mails := make([]string, len(m.To), len(m.To))
	for i, addr := range m.To {
		mails[i] = addr.Address
	}
	content, err := m.makeContent()
	if err != nil {
		return err
	}
	return smtp.SendMail(server, auth, m.From.Address, mails, content)
}
