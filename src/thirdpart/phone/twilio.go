package phone

import (
	"broker"
	"encoding/json"
	"fmt"
	"model"
	"net/url"
)

type Twilio struct {
	url  string
	from string
}

func NewTwilio(config *model.Config) *Twilio {
	return &Twilio{
		url:  config.Thirdpart.Sms.Twilio.Url,
		from: config.Thirdpart.Sms.Twilio.FromPhone,
	}
}

func (t *Twilio) Codes() []string {
	return []string{
		"+1",   // American, Canada
		"+45",  // Denmark
		"+7",   // Russia
		"+852", // Hong Kong
		"+886", // Taiwan
	}
}

type twilioReply struct {
	SID string `json:"sid"`
}

func (t *Twilio) Send(phone string, content string) (string, error) {
	params := make(url.Values)
	params.Add("From", t.from)
	params.Add("To", phone)
	params.Add("Body", content)
	resp, err := broker.HttpForm(t.url, params)
	if err != nil {
		return "", fmt.Errorf("send to %s failed: %s", phone, err)
	}
	defer resp.Close()
	decoder := json.NewDecoder(resp)
	var reply twilioReply
	err = decoder.Decode(&reply)
	if err != nil {
		return "", fmt.Errorf("send to %s reply decode failed: %s", phone, err)
	}
	return reply.SID, nil
}
