package imessage

import (
	"broker"
	"encoding/json"
	"fmt"
	"formatter"
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
	conn *multiplexer.Homo
	log  *logger.SubLogger
}

func New(config *model.Config) (*IMessage, error) {
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
		log:  log,
	}, nil
}

func (i *IMessage) Provider() string {
	return "imessage"
}

func (i *IMessage) Check(to string) (ret bool, err error) {
	i.conn.Do(func(i multiplexer.Instance) {
		imsg, ok := i.(*IMessageConn)
		if !ok {
			err = fmt.Errorf("instance %+v is not *IMessageConn", i)
			return
		}
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

func (i *IMessage) Send(to *model.Recipient, text string) (id string, err error) {
	phone := to.ExternalID

	if phone[:3] == "+86" {
		p := phone[3:]
		ok, err := i.Check(p)
		if err != nil {
			return "", fmt.Errorf("imessage error: %s", err)
		} else if ok {
			phone = p
			i.log.Debug("phone %s is imessage", p)
		} else {
			return "", fmt.Errorf("%s not imessage", phone)
		}
	} else {
		return "", fmt.Errorf("unsupport phone %s", phone)
	}

	lines := strings.Split(text, "\n")
	contents := make([]string, 0)
	for _, line := range lines {
		line = strings.Trim(line, " \n\r\t")
		if line == "" {
			continue
		}

		cutter, err := formatter.CutterParse(line, imsgLen)
		if err != nil {
			return "", fmt.Errorf("parse cutter error: %s", err)
		}

		for _, content := range cutter.Limit(300) {
			contents = append(contents, content)
		}
	}
	return i.SendMessage(phone, contents)
}

func imsgLen(content string) int {
	return len([]byte(content))
}

func (i *IMessage) SendMessage(to string, contents []string) (id string, err error) {
	i.conn.Do(func(i multiplexer.Instance) {
		imsg, ok := i.(*IMessageConn)
		if !ok {
			err = fmt.Errorf("instance %+v is not *IMessageConn", i)
			return
		}
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
