package apn_service

import (
	"github.com/virushuo/Go-Apns"
	"fmt"
	"log/syslog"
	"encoding/json"
)

type Apn struct {
	cert string
	key string
	server string
	apn *goapns.Apn
	id uint32
	log *syslog.Writer
}

func NewApn(cert, key, server string, log *syslog.Writer) (*Apn, error) {
	apn, err := goapns.Connect(cert, key, server)
	if err != nil {
		return nil, err
	}
	ret := &Apn{
		cert: cert,
		key: key,
		server: server,
		apn: apn,
		id: 0,
		log: log,
	}
	go errorListen(ret)
	return ret, nil
}

func errorListen(apn *Apn) {
	for {
		apnerr := <-apn.apn.Errorchan
		apn.log.Err(fmt.Sprintf("Apn error: cmd %d, status %d, id %d", apnerr.Command, apnerr.Status, apnerr.Identifier))
		var err error
		err = apn.apn.Reconnect()
		if err != nil {
			apn.log.Err(fmt.Sprintf("Reconnect to apn server(%s) error: %s", apn.server, err))
			panic(err)
		}
	}
}

type ExfePush struct {
	Cid uint64
	T string
}

func (p *ExfePush) MarshalJSON() ([]byte, error) {
	t, _ := json.Marshal(p.T)
	return []byte(fmt.Sprintf("{\"cid\":\"%d\",\"t\":%s}", p.Cid, t)), nil
}

type ApnSendArg struct {
	DeviceToken string
	Alert string
	Badge uint
	Sound string
	Cid uint64
	T string
}

func (a *Apn) ApnSend(args []ApnSendArg) error {
	for _, arg := range args {
		notification := goapns.Notification{
			Device_token: arg.DeviceToken,
			Alert: arg.Alert,
			Badge: arg.Badge,
			Sound: arg.Sound,
			Args: ExfePush{
				Cid: arg.Cid,
				T: arg.T,
			},
			Identifier: a.id,
		}
		a.id++
		err := a.apn.SendNotification(&notification)
		if err != nil {
			a.log.Err(fmt.Sprintf("Send notification(%s) to device(%s) error: %s", arg.Alert, arg.DeviceToken, err))
		}
	}
	return nil
}
