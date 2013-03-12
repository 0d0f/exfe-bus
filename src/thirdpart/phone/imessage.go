package sms

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/googollee/go-multiplexer"
	"github.com/googollee/go-socket.io"
	"model"
	"time"
)

type IMessageConn struct {
	conn *socketio.Client
	log  *logger.SubLogger
}

func (i *IMessageConn) Ping() error {
	return nil
}

func (i *IMessageConn) Close() error {
	i.conn.SetDeadline(time.Now().Add(broker.NetworkTimeout))
	return i.conn.Close()
}

func (i *IMessageConn) Error(err error) {
	i.log.Err("%s", err)
}

type IMessage struct {
	conn *multiplexer.Homo
}

func NewIMessage(config *model.Config) (*IMessage, error) {
	log := config.Log.SubPrefix("imessage")
	homo := multiplexer.NewHomo(func() (multiplexer.Instance, error) {
		sio, err := socketio.Dial("http://www.kufuwu.com:1080/socket.io/", "http://www.kufuwu.com/p.html", broker.NetworkTimeout)
		if err != nil {
			return nil, err
		}
		return &IMessageConn{
			conn: sio,
			log:  log,
		}, nil
	}, 5, 30*time.Second, 40*time.Second)
	return &IMessage{
		conn: homo,
	}, nil
}

func (i *IMessage) Check(to string) (ret bool, err error) {
	err = i.conn.Do(func(i multiplexer.Instance) {
		imsg, ok := i.(*IMessageConn)
		if !ok {
			err = fmt.Errorf("instance %+v is not *IMessageConn", i)
			return
		}
		imsg.conn.SetDeadline(time.Now().Add(broker.NetworkTimeout))
		req := Request{
			To:      to,
			Channel: getChannel(),
			Action:  "1",
		}
		err = imsg.conn.Emit(true, "send", req)
		if err != nil {
			return
		}
		var msg socketio.Message
		err = imsg.conn.Receive(&msg)
		if err != nil {
			return
		}
		var respString string
		err = msg.ReadArguments(&respString)
		if err != nil {
			return
		}
		var resp Response
		err = json.Unmarshal([]byte(respString), &resp)
		if err != nil {
			return
		}
		ret = resp.Head.Status == 0
	})
	return
}

func (i *IMessage) Send(to string, contents []string) (id string, err error) {
	err = i.conn.Do(func(i multiplexer.Instance) {
		imsg, ok := i.(*IMessageConn)
		if !ok {
			err = fmt.Errorf("instance %+v is not *IMessageConn", i)
			return
		}
		imsg.conn.SetDeadline(time.Now().Add(broker.ProcessTimeout))
		for _, content := range contents {
			req := Request{
				To:      to,
				Channel: getChannel(),
				Action:  "2",
				Message: content,
			}
			err = imsg.conn.Emit(true, "send", req)
			if err != nil {
				return
			}
			var msg socketio.Message
			err = imsg.conn.Receive(&msg)
			if err != nil {
				return
			}
			var respString string
			err = msg.ReadArguments(&respString)
			if err != nil {
				return
			}
			var resp Response
			err = json.Unmarshal([]byte(respString), &resp)
			if err != nil {
				return
			}
			if resp.Head.Status != 0 {
				err = fmt.Errorf("%s", resp.Head.Err)
				return
			}
			id += "," + resp.Head.ID
		}
	})
	if len(id) > 0 {
		id = id[1:]
	}
	return
}

func (i *IMessage) Codes() []string {
	return nil
}

type Request struct {
	To      string `json:"to"`
	Channel string `json:"channel"`
	Action  string `json:"act"`
	Message string `json:"message,omitempty"`
}

type Response struct {
	Head struct {
		Status int    `json:"status"`
		To     string `json:"to"`
		Action string `json:"act"`
		ID     string `json:"id"`
		Err    string `json:"errmsg"`
	} `json:"head"`
}

func getChannel() string {
	return "4"
}
