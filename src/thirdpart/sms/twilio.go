package sms

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"io/ioutil"
	"model"
	"net/http"
	"net/url"
)

type Twilio struct {
	url  string
	from string
	log  *logger.SubLogger
}

func NewTwilio(config *model.Config) *Twilio {
	return &Twilio{
		url:  config.Thirdpart.Sms.Twilio.Url,
		from: config.Thirdpart.Sms.Twilio.FromPhone,
		log:  config.Log.SubPrefix("sms-twilio"),
	}
}

func (t *Twilio) Codes() []string {
	return []string{"+1"}
}

type twilioReply struct {
	SID string `json:"sid"`
}

func (t *Twilio) Send(phone string, contents []string) (string, error) {
	params := make(url.Values)
	params.Add("From", t.from)
	params.Add("To", phone)
	ret := ""
	for _, content := range contents {
		params.Add("Body", content)
		resp, err := http.PostForm(t.url, params)
		if err != nil {
			return "", fmt.Errorf("send to %s failed: %s", phone, err)
		}
		if resp.StatusCode != 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("send to %s response: %s", phone, resp.Status)
			} else {
				return "", fmt.Errorf("send to %s response: %s(%s)", phone, resp.Status, string(body))
			}
		}
		decoder := json.NewDecoder(resp.Body)
		var reply twilioReply
		err = decoder.Decode(&reply)
		if err != nil {
			t.log.Err("send to %s reply decode failed: %s", phone, err)
		}
		ret += "," + reply.SID
	}
	if len(ret) > 0 {
		ret = ret[1:]
	}
	return ret, nil
}
