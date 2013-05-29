package imessage

import (
	"broker"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/googollee/go-multiplexer"
	"github.com/googollee/go-socket.io"
	"model"
	"strings"
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
	return i.conn.Close()
}

func (i *IMessageConn) Error(err error) {
	i.log.Err("%s", err)
}

type IMessage struct {
	conn    *multiplexer.Homo
	log     *logger.SubLogger
	channel string
}

func New(config *model.Config) (*IMessage, error) {
	log := config.Log.SubPrefix("imessage")
	homo := multiplexer.NewHomo(func() (multiplexer.Instance, error) {
		sio, err := socketio.Dial(config.Thirdpart.IMessage.Address, config.Thirdpart.IMessage.Origin, broker.NetworkTimeout)
		if err != nil {
			return nil, err
		}
		return &IMessageConn{
			conn: sio,
			log:  log,
		}, nil
	}, 5, 30*time.Second, 40*time.Second)
	return &IMessage{
		conn:    homo,
		log:     log,
		channel: config.Thirdpart.IMessage.Channel,
	}, nil
}

func (i *IMessage) Provider() string {
	return "imessage"
}

func (im *IMessage) Check(to string) (ret bool, err error) {
	im.conn.Do(func(i multiplexer.Instance) {
		imsg, ok := i.(*IMessageConn)
		if !ok {
			err = fmt.Errorf("instance %+v is not *IMessageConn", i)
			return
		}
		req := Request{
			To:      to,
			Channel: im.getChannel(),
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
		var resp Response
		err = msg.ReadArguments(&resp)
		if err != nil {
			return
		}
		ret = resp.Head.Status == 0
	})
	return
}

func (i *IMessage) Post(id, text string) (string, error) {
	text = strings.Trim(text, " \r\n")
	ok, err := i.Check(id)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("%s is not iMessage", id)
	}
	return i.SendMessage(id, text)
}

func (im *IMessage) SendMessage(to string, content string) (id string, err error) {
	err = im.conn.Do(func(i multiplexer.Instance) {
		imsg, ok := i.(*IMessageConn)
		if !ok {
			err = fmt.Errorf("instance %+v is not *IMessageConn", i)
			return
		}
		req := Request{
			To:      to,
			Channel: im.getChannel(),
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
		var resp Response
		err = msg.ReadArguments(&resp)
		if err != nil {
			return
		}
		if resp.Head.Status != 0 {
			err = fmt.Errorf("%s", resp.Head.Err)
			return
		}
		id = resp.Head.ID
	})
	return
}

func (i *IMessage) getChannel() string {
	return i.channel
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
