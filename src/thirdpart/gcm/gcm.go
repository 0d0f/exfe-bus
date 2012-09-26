package gcm

import (
	"fmt"
	"github.com/googollee/go-gcm"
	"model"
	"thirdpart"
)

type GCM struct {
	client *gcm.Client
}

func NewGCM(key string) *GCM {
	return &GCM{
		client: gcm.New(key),
	}
}

func (g *GCM) Provider() string {
	return "Android"
}

func (g *GCM) MessageType() thirdpart.MessageType {
	return thirdpart.ShortMessage
}

func (g *GCM) Send(to *model.Recipient, privateMessage string, publicMessage string, data *thirdpart.InfoData) (id string, err error) {
	message := gcm.NewMessage(to.ExternalID)
	message.AddPayload("text", privateMessage)
	message.AddPayload("badge", 1)
	message.AddPayload("sound", "")
	message.AddPayload("cid", data.CrossID)
	message.AddPayload("t", data.Type)
	message.DelayWhileIdle = true
	message.CollapseKey = "exfe"

	resp, err := g.client.Send(message)
	if err != nil {
		return "", err
	}

	for i := range resp.Results {
		if resp.Results[i].RegistrationID != to.ExternalID {
			continue
		}
		if err := resp.Results[i].Error; err != "" {
			return resp.Results[i].MessageID, fmt.Errorf(err)
		}
		return resp.Results[i].MessageID, nil
	}
	return "", fmt.Errorf("can't find result in response")
}
