package gcm_service

import (
	"github.com/googollee/go-gcm"
	"log"
)

type GCM struct {
	client *gcm.Client
}

func NewGCM(key string) *GCM {
	return &GCM{
		client: gcm.New(key),
	}
}

type SendArg struct {
	DeviceID string
	Text     string
	Badge    uint
	Sound    string
	Cid      uint64
	T        string
}

func (c *GCM) Send(args []SendArg) error {
	for _, arg := range args {
		log.Printf("Sending message(%s) to device(%s)", arg.Text, arg.DeviceID)
		message := gcm.NewMessage(arg.DeviceID)
		message.AddPayload("text", arg.Text)
		message.AddPayload("badge", arg.Badge)
		message.AddPayload("sound", arg.Sound)
		message.AddPayload("cid", arg.Cid)
		message.AddPayload("t", arg.T)
		message.DelayWhileIdle = true
		message.CollapseKey = "exfe"

		resp, err := c.client.Send(message)
		if err != nil {
			log.Printf("net error: %s", err)
		} else {
			errors := resp.ErrorIndexes()
			for _, i := range errors {
				log.Printf("google report error: %s", resp.Results[i].Error)
			}
		}
	}
	return nil
}
