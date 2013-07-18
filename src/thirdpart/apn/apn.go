package apn

import (
	"encoding/json"
	"fmt"
	"github.com/virushuo/Go-Apns"
	"logger"
	"regexp"
	"strings"
	"thirdpart"
	"time"
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
	f      thirdpart.Callback
}

type ErrorHandler func(err apns.NotificationError)

func New(broker Broker) *Apn {
	ret := &Apn{
		broker: broker,
		id:     0,
	}
	go listenError(broker.GetErrorChan(), func(err apns.NotificationError) {
		if ret.f != nil {
			ret.f(fmt.Sprintf("%d", err.Identifier), err)
		}
	})
	return ret
}

func (a *Apn) Provider() string {
	return "iOS"
}

func (a *Apn) SetPosterCallback(callback thirdpart.Callback) (time.Duration, bool) {
	a.f = callback
	return time.Second * 30, true
}

func (a *Apn) Post(from, id, text string) (string, error) {
	text = strings.Trim(text, " \r\n")
	last := strings.LastIndex(text, "\n")
	if last == -1 {
		return "", fmt.Errorf("no payload")
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
		fmt.Println("here:", err)
		return fmt.Sprint("%d", ret), fmt.Errorf("send %d error: %s", ret, err)
	}
	fmt.Println("ret:", ret)
	return fmt.Sprint("%d", ret), nil
}

func listenError(errChan <-chan error, h ErrorHandler) {
	for {
		err := <-errChan
		e, ok := err.(apns.NotificationError)
		if !ok {
			logger.ERROR("unknow err: %s", err)
			continue
		}
		h(e)
	}
}

var tailUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?$`)
