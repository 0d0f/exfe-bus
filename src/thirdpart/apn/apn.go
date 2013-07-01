package apn

import (
	"encoding/json"
	"fmt"
	"github.com/virushuo/Go-Apns"
	"regexp"
	"strings"
)

type Broker interface {
	Send(n *apns.Notification) error
	GetErrorChan() <-chan error
}

type sendArg struct {
	id      uint
	content string
}

type Apn struct {
	broker Broker
	id     uint32
}

type ErrorHandler func(err error)

func New(broker Broker, errorHandler ErrorHandler) *Apn {
	go listenError(broker.GetErrorChan(), errorHandler)
	return &Apn{
		broker: broker,
		id:     0,
	}
}

func (a *Apn) Provider() string {
	return "iOS"
}

func (a *Apn) Post(from, id, text string) (string, error) {
	text = strings.Trim(text, " \r\n")
	last := strings.LastIndex(text, "\n")
	if last == -1 {
		return "", nil
	}
	dataStr := text[last+1:]
	var data interface{}
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return "", fmt.Errorf("last line of text(%s) can't unmarshal: %s", dataStr, err)
	}
	text = strings.Trim(text[:last], " \r\n")
	text = tailUrlRegex.ReplaceAllString(text, "")

	ret := a.id
	a.id++

	payload := apns.Payload{}
	payload.Aps.Alert.Body = text
	payload.Aps.Badge = 1
	payload.Aps.Sound = "default"
	if data != nil {
		payload.SetCustom("args", data)
	}
	notification := apns.Notification{
		DeviceToken: id,
		Identifier:  ret,
		Payload:     &payload,
	}

	err = a.broker.Send(&notification)
	if err != nil {
		return fmt.Sprint("%d", ret), fmt.Errorf("send %d error: %s", ret, err)
	}
	return fmt.Sprint("%d", ret), nil
}

func listenError(errChan <-chan error, h ErrorHandler) {
	for {
		h(<-errChan)
	}
}

var tailUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?$`)
