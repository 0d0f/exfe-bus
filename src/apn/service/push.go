package apn_service

import (
	"github.com/virushuo/Go-Apns"
	"log"
	"old_gobus"
	"time"
	"unicode/utf8"
)

type Apn struct {
	apn *apns.Apn
	id  uint32

	netaddr    string
	db         int
	password   string
	pushBuffer []ApnSendArg
}

func NewApn(cert, key, server, rootca, netaddr string, db int, password string) (*Apn, error) {
	apn, err := apns.New(cert, key, server, 30*time.Minute)
	if err != nil {
		return nil, err
	}
	ret := &Apn{
		apn:      apn,
		id:       0,
		netaddr:  netaddr,
		db:       db,
		password: password,
	}
	go errorListen(ret)
	return ret, nil
}

func errorListen(apn *Apn) {
	for {
		apnerr := <-apn.apn.ErrorChan
		if apnerr.Command != 0 && apnerr.Status != 0 {
			log.Printf("Apn error: cmd %d, status %d, id %d", apnerr.Command, apnerr.Status, apnerr.Identifier)
			needSend := false
			for _, push := range apn.pushBuffer {
				if push.id == apnerr.Identifier {
					needSend = true
					continue
				}
				if !needSend {
					continue
				}
				client := gobus.CreateClient(apn.netaddr, apn.db, apn.password, "iOS")
				client.Send("ApnSend", &push, 5)
			}
			panic(apnerr)
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

	id uint32
}

func (a *Apn) ApnSend(args []ApnSendArg) error {
	payload := apns.Payload{}
	a.pushBuffer = args

	for i, arg := range args {
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
		notification := apns.Notification{
			DeviceToken: arg.DeviceToken,
			Identifier:  a.id,
			Payload:     &payload,
		}
		a.pushBuffer[i].id = a.id

		log.Printf("apn send %s to %s, id %d", notification.Payload.Aps.Alert, notification.DeviceToken, notification.Identifier)
		err := a.apn.Send(&notification)
		a.id++
		if err != nil {
			log.Printf("Send notification(%s) to device(%s) error: %s", arg.Alert, arg.DeviceToken, err)
		}
	}
	return nil
}
