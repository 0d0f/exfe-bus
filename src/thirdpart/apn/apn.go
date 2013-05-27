package apn

import (
	"encoding/json"
	"fmt"
	"formatter"
	"github.com/virushuo/Go-Apns"
	"model"
	"regexp"
	"strings"
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

func (a *Apn) Send(to *model.Recipient, text string) (string, error) {
	ids := ""
	lines := strings.Split(text, "\n")
	dataStr := lines[len(lines)-1]
	var data interface{}
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return "", fmt.Errorf("last line of text(%s) can't unmarshal: %s", dataStr, err)
	}
	lines = lines[:len(lines)-1]
	for _, line := range lines {
		line = strings.Trim(line, " \n\r\t")
		line = tailUrlRegex.ReplaceAllString(line, "")
		line = tailQuoteUrlRegex.ReplaceAllString(line, `)\)`)
		if line == "" {
			continue
		}

		cutter, err := formatter.CutterParse(line, apnLen)
		if err != nil {
			return "", fmt.Errorf("parse cutter error: %s", err)
		}

		for _, content := range cutter.Limit(140) {
			id := a.id
			a.id++
			ids = fmt.Sprintf("%s,%d", ids, id)

			payload := apns.Payload{}
			payload.Aps.Alert.Body = content
			payload.Aps.Badge = 1
			payload.Aps.Sound = "default"
			if data != nil {
				payload.SetCustom("args", data)
			}
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
	}
	if ids != "" {
		ids = ids[1:]
	}
	return ids, nil
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

var tailUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?$`)
var tailQuoteUrlRegex = regexp.MustCompile(` *(http|https):\/\/exfe.com(\/[\w#!:.?+=&%@!\-\/]*)?\)(\\\))$`)
