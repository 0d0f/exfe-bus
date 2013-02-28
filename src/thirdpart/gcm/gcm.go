package gcm

import (
	"encoding/json"
	"fmt"
	"formatter"
	"github.com/googollee/go-gcm"
	"model"
	"regexp"
	"strings"
	"unicode/utf8"
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

func (g *GCM) Send(to *model.Recipient, text string) (id string, err error) {
	ids := ""
	lines := strings.Split(text, "\n")
	dataStr := lines[len(lines)-1]
	lines = lines[:len(lines)-1]
	var data map[string]interface{}
	err = json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return "", err
	}

	for _, line := range lines {
		line = strings.Trim(line, " \r\n\t")
		line = tailUrlRegex.ReplaceAllString(line, "")
		line = tailQuoteUrlRegex.ReplaceAllString(line, `)\)`)
		if line == "" {
			continue
		}

		cutter, err := formatter.CutterParse(line, gcmLen)
		if err != nil {
			return "", fmt.Errorf("parse cutter error: %s", err)
		}
		message := gcm.NewMessage(to.ExternalID)
		message.SetPayload("badge", "1")
		message.SetPayload("sound", "")
		if len(data) > 0 {
			for k, v := range data {
				message.SetPayload(k, fmt.Sprintf("%v", v))
			}
		}
		message.DelayWhileIdle = true
		message.CollapseKey = "exfe"

		for _, content := range cutter.Limit(140) {
			message.SetPayload("text", content)
			resp, err := g.broker.Send(message)
			if err != nil {
				return "", fmt.Errorf("send to %s error: %s", to, err)
			}

			for i := range resp.Results {
				if resp.Results[i].RegistrationID != to.ExternalID {
					continue
				}
				if err := resp.Results[i].Error; err != "" {
					return resp.Results[i].MessageID, fmt.Errorf("send to %s error: %s", to, err)
				}
				ids = fmt.Sprintf("%s,%d", ids, resp.Results[i].MessageID)
			}
		}
	}
	if ids != "" {
		ids = ids[1:]
	}
	return ids, nil
}

func gcmLen(content string) int {
	return utf8.RuneCountInString(content)
}

var tailUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?$`)
var tailQuoteUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?\)(\\\))$`)
