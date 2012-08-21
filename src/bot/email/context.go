package email

import (
	"email/service"
	"fmt"
	"gobot"
	"io/ioutil"
	"net/http"
	"net/mail"
	"net/textproto"
	"net/url"
	"time"
)

type Email struct {
	From      *mail.Address
	To        []*mail.Address
	Subject   string
	CrossId   string
	Date      time.Time
	MessageId string
	Text      string
}

func (m *Email) ToUrlValues() url.Values {
	ret := make(url.Values)
	ret.Add("cross_id", m.CrossId)
	ret.Add("content", m.Text)
	ret.Add("external_id", m.From.Address)
	ret.Add("time", m.Date.Format("2006-01-02 15:04:05 -0700"))
	ret.Add("provider", "email")
	return ret
}

type EmailContext struct {
	*bot.BaseContext
	mailBot *EmailBot
}

func NewEmailContext(id string, b *EmailBot) *EmailContext {
	return &EmailContext{bot.NewBaseContext(id), b}
}

func (c *EmailContext) EmailWithCrossID(input *Email) error {
	if input.CrossId == "" {
		return bot.BotNotMatched
	}
	url := fmt.Sprintf("%s/v2/gobus/PostConversation", c.mailBot.config.Site_api)
	fmt.Printf("send to: %s, post content: %s\n", url, input.ToUrlValues().Encode())
	resp, err := http.PostForm(url, input.ToUrlValues())
	if err != nil {
		return fmt.Errorf("message(%s) send to server error: %s", input.MessageId, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("message(%s) get response body error: %s", input.MessageId, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("message(%s) send error(%s): %s", input.MessageId, resp.Status, string(body))
	}
	return nil
}

func (c *EmailContext) Default(input *Email) error {
	body := fmt.Sprintf("Sorry for the inconvenience, but email you just sent to EXFE was not sent from an attendee identity to the X (cross). Please try again from the correct email address.\n\n--\n%s",
		input.Text)
	mailarg := &email_service.MailArg{
		To:      []*mail.Address{input.From},
		From:    &mail.Address{"x@exfe.com", "x@exfe.com"},
		Subject: fmt.Sprintf("Re: %s", input.Subject),
		Text:    body,
		Html:    fmt.Sprintf("<html><body><p>%s</p></body></html>", body),
		Header:  make(textproto.MIMEHeader),
	}
	mailarg.Header.Set("In-Reply-To", input.MessageId)
	mailarg.Header.Set("References", input.MessageId)
	c.mailBot.bus.Send("EmailSend", &mailarg, 5)
	return nil
}
