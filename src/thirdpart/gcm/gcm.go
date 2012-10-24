package gcm

import (
	"fmt"
	"formatter"
	"github.com/googollee/go-gcm"
	"model"
	"regexp"
	"strings"
	"thirdpart"
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

func (g *GCM) Send(to *model.Recipient, privateMessage string, publicMessage string, data *thirdpart.InfoData) (id string, err error) {
	ids := ""
	privateMessage = urlRegex.ReplaceAllString(privateMessage, "")
	for _, line := range strings.Split(privateMessage, "\n") {
		line = strings.Trim(line, " \r\n\t")
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
		message.SetPayload("cid", fmt.Sprintf("%d", data.CrossID))
		message.SetPayload("t", data.Type.String())
		message.DelayWhileIdle = true
		message.CollapseKey = "exfe"

		for _, content := range cutter.Limit(140) {
			message.SetPayload("text", content)
			resp, err := g.broker.Send(message)
			if err != nil {
				return "", fmt.Errorf("send to %s error: %s", to, err)
			}

			fmt.Println(resp.Results)
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

var urlRegex = regexp.MustCompile(` *(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`)
