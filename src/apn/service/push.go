package apn_service

import (
	"github.com/virushuo/Go-Apns"
	"log"
	"time"
	"unicode/utf8"
)

type Apn struct {
	cert     string
	key      string
	server   string
	apn      *goapns.Apn
	id       uint32
	isClosed bool
	retimer  chan int
}

func NewApn(cert, key, server, rootca string) (*Apn, error) {
	apn, err := goapns.Connect(cert, key, server)
	if err != nil {
		return nil, err
	}
	ret := &Apn{
		cert:    cert,
		key:     key,
		server:  server,
		apn:     apn,
		id:      0,
		retimer: make(chan int),
	}
	go errorListen(ret)
	go closeTimer(ret)
	return ret, nil
}

func errorListen(apn *Apn) {
	for {
		apnerr := <-apn.apn.ErrorChan
		if apnerr.Command != 0 && apnerr.Status != 0 {
			log.Printf("Apn error: cmd %d, status %d, id %d", apnerr.Command, apnerr.Status, apnerr.Identifier)
			panic("apn error")
		}
	}
}

func closeTimer(apn *Apn) {
	for {
		timer := time.After(30 * time.Minute)
		select {
		case <-timer:
			log.Printf("Apn connection timeout")
			apn.isClosed = true
			apn.apn.Close()
			return
		case <-apn.retimer:
		}
	}
}

type ExfePush struct {
	Cid uint64 `json:"cid"`
	T   string `json:"t"`
}

type ApnSendArg struct {
	DeviceToken string
	Alert       string
	Badge       uint
	Sound       string
	Cid         uint64
	T           string
}

func (a *Apn) ApnSend(args []ApnSendArg) error {
	if a.isClosed {
		log.Printf("apn connection reconnect")
		a.apn.Reconnect()
		a.isClosed = false
		go closeTimer(a)
	}

	defer func() { a.retimer <- 1 }()

	payload := goapns.Payload{}
	for _, arg := range args {
		alert := []byte(arg.Alert)
		if len(alert) > 140 {
			alert = alert[:140]
			for !utf8.Valid(alert) {
				alert = alert[:len(alert)-1]
			}
		}
		payload.Aps.Alert = string(alert)
		payload.Aps.Badge = int(arg.Badge)
		payload.Aps.Sound = arg.Sound
		payload.SetCustom("args", ExfePush{
			Cid: arg.Cid,
			T:   arg.T,
		})
		notification := goapns.Notification{
			DeviceToken: arg.DeviceToken,
			Identifier:  a.id,
			Payload:     &payload,
		}
		a.id++
		log.Printf("apn send %s to %s, id %d", notification.Payload.Aps.Alert, notification.DeviceToken, notification.Identifier)
		err := a.apn.SendNotification(&notification)
		if err != nil {
			log.Printf("Send notification(%s) to device(%s) error: %s", arg.Alert, arg.DeviceToken, err)
			return err
		}
	}
	return nil
}

func (a *Apn) Close() error {
	if a.isClosed {
		return nil
	}
	a.isClosed = true
	return a.apn.Close()
}
