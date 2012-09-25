package apn

import (
	"fmt"
	"github.com/virushuo/Go-Apns"
	"model"
	"thirdpart"
	"time"
	"unicode/utf8"
)

type sendArg struct {
	id      uint
	content string
}

type Apn struct {
	conn *apns.Apn
	id   uint32
}

type ErrorHandler func(apns.NotificationError)

func New(cert, key, server string, timeout time.Duration, errorHandler ErrorHandler) (*Apn, error) {
	apn, err := apns.New(cert, key, server, timeout)
	if err != nil {
		return nil, err
	}
	go listenError(apn, errorHandler)
	return &Apn{
		conn: apn,
		id:   0,
	}, nil
}

func (a *Apn) Provider() string {
	return "iOS"
}

func (a *Apn) MessageType() thirdpart.MessageType {
	return thirdpart.ShortMessage
}

func (a *Apn) Send(to *model.Recipient, privateMessage string, publicMessage string, data *thirdpart.InfoData) (string, error) {
	id := a.id
	a.id++

	alert := []byte(privateMessage)
	if len(alert) > 140 {
		alert = alert[:140]
		for !utf8.Valid(alert) {
			alert = alert[:len(alert)-1]
		}
	}

	payload := apns.Payload{}
	payload.Aps.Alert = string(alert)
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

	err := a.conn.Send(&notification)
	return fmt.Sprintf("%d", id), err
}

type ExfePush struct {
	Cid uint64 `json:"cid"`
	T   string `json:"t"`
}

func listenError(conn *apns.Apn, h ErrorHandler) {
	for {
		h(<-conn.ErrorChan)
	}
}
