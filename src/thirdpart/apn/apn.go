package apn

import (
	"fmt"
	"github.com/virushuo/Go-Apns"
	"model"
	"regexp"
	"textcutter"
	"thirdpart"
)

type Broker interface {
	Send(n *apns.Notification) error
	GetErrorChan() <-chan apns.NotificationError
}

type sendArg struct {
	id      uint
	content string
}

type Apn struct {
	broker Broker
	id     uint32
}

type ErrorHandler func(apns.NotificationError)

func New(broker Broker, errorHandler ErrorHandler) (*Apn, error) {
	go listenError(broker.GetErrorChan(), errorHandler)
	return &Apn{
		broker: broker,
		id:     0,
	}, nil
}

func (a *Apn) Provider() string {
	return "iOS"
}

func (a *Apn) MessageType() thirdpart.MessageType {
	return thirdpart.ShortMessage
}

func (a *Apn) Send(to *model.Recipient, privateMessage string, publicMessage string, data *thirdpart.InfoData) (string, error) {
	privateMessage = urlRegex.ReplaceAllString(privateMessage, "")
	cutter, err := textcutter.Parse(privateMessage, apnLen)
	if err != nil {
		return "", fmt.Errorf("parse cutter error: %s", err)
	}

	ids := ""
	for _, content := range cutter.Limit(140) {
		id := a.id
		a.id++
		ids = fmt.Sprintf("%s,%d", ids, id)

		payload := apns.Payload{}
		payload.Aps.Alert = content
		payload.Aps.Badge = 1
		payload.Aps.Sound = ""
		payload.SetCustom("args", ExfePush{
			Cid: data.CrossID,
			T:   data.Type.String(),
		})
		notification := apns.Notification{
			DeviceToken: to.ExternalID,
			Identifier:  id,
			Payload:     &payload,
		}

		err := a.broker.Send(&notification)
		if err != nil {
			return ids, fmt.Errorf("send %d error: %s", id, err)
		}
	}
	return ids[1:], nil
}

type ExfePush struct {
	Cid uint64 `json:"cid"`
	T   string `json:"t"`
}

func listenError(errChan <-chan apns.NotificationError, h ErrorHandler) {
	for {
		h(<-errChan)
	}
}

func apnLen(content string) int {
	return len([]byte(content))
}

var urlRegex = regexp.MustCompile(` *(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`)
