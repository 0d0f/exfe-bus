package gcm

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-gcm"
	"regexp"
	"strings"
)

type Broker interface {
	Send(message *gcm.Message) (*gcm.Response, error)
}

type GCM struct {
	broker Broker
}

func New(broker Broker) *GCM {
	return &GCM{
		broker: broker,
	}
}

func (g *GCM) Provider() string {
	return "Android"
}

func (g *GCM) Post(id, text string) (string, error) {
	text = strings.Trim(text, " \r\n")
	last := strings.LastIndex(text, "\n")
	dataStr := text[last+1:]
	var data map[string]string
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return "", fmt.Errorf("last line of text(%s) can't unmarshal: %s", dataStr, err)
	}
	text = strings.Trim(text[:last], " \r\n")
	text = tailUrlRegex.ReplaceAllString(text, "")

	message := gcm.NewMessage(id)
	message.SetPayload("badge", "1")
	message.SetPayload("sound", "")
	if len(data) > 0 {
		for k, v := range data {
			message.SetPayload(k, fmt.Sprintf("%v", v))
		}
	}
	message.DelayWhileIdle = true
	message.CollapseKey = "exfe"
	message.SetPayload("text", text)
	resp, err := g.broker.Send(message)
	if err != nil {
		return "", fmt.Errorf("send to %s@Android error: %s", id, err)
	}

	for i := range resp.Results {
		if resp.Results[i].RegistrationID != id {
			continue
		}
		if err := resp.Results[i].Error; err != "" {
			return "", fmt.Errorf("send to %s@Android error: (%s)%s", id, resp.Results[i].MessageID, err)
		}
		return resp.Results[i].MessageID, nil
	}
	return "", fmt.Errorf("parse result failed: no result")
}

var tailUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?$`)
