package email

import (
	"bytes"
	"fmt"
	"gobot"
	"io/ioutil"
	"model"
	"net/http"
	"net/mail"
	"net/url"
	"time"
)

type Email struct {
	From      *mail.Address
	To        []*mail.Address
	Subject   string
	CrossID   string
	Date      time.Time
	MessageID string
	Text      string
}

func (m *Email) ToUrlValues() url.Values {
	ret := make(url.Values)
	ret.Add("cross_id", m.CrossID)
	ret.Add("content", m.Text)
	ret.Add("external_id", m.From.Address)
	ret.Add("time", m.Date.Format("2006-01-02 15:04:05 -0700"))
	ret.Add("provider", "email")
	return ret
}

var badrequest = fmt.Errorf("email bad request")

type EmailContext struct {
	*bot.BaseContext
	mailBot *EmailBot
}

func NewEmailContext(id string, b *EmailBot) *EmailContext {
	return &EmailContext{bot.NewBaseContext(id), b}
}

func (c *EmailContext) EmailWithCrossID(input *Email) error {
	if input.CrossID == "" {
		return bot.BotNotMatched
	}
	url := fmt.Sprintf("%s/v2/gobus/PostConversation", c.mailBot.config.SiteApi)
	c.mailBot.config.Log.Debug("send to: %s, post content: %s\n", url, input.ToUrlValues().Encode())
	resp, err := http.PostForm(url, input.ToUrlValues())
	if err != nil {
		return fmt.Errorf("message(%s) send to server error: %s", input.MessageID, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("message(%s) get response body error: %s", input.MessageID, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("message(%s) send error(%s): %s", input.MessageID, resp.Status, string(body))
	}
	return nil
}

func (c *EmailContext) Default(input *Email) error {
	buf := bytes.NewBuffer(nil)
	err := c.mailBot.localTemplate.Execute(buf, "en_US", "conversation_reply.email", input)
	if err != nil {
		c.mailBot.config.Log.Crit("template(conversation_reply.email) failed: %s", err)
		return badrequest
	}

	arg := model.ThirdpartSend{
		PrivateMessage: buf.String(),
		PublicMessage:  "",
		Info: &model.InfoData{
			CrossID: 0,
			Type:    model.TypeCrossInvitation,
		},
	}
	arg.To = model.Recipient{
		Provider:         "email",
		ExternalID:       input.From.Address,
		ExternalUsername: input.From.Address,
	}
	var ids string
	err = c.mailBot.sender.Do("Send", &arg, &ids)
	if err != nil {
		c.mailBot.config.Log.Crit("send error: %s", err)
	}

	return badrequest
}
