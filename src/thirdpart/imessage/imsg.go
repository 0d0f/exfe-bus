package imessage

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/stathat/consistent"
	"logger"
	"model"
	"net"
	"time"
	"valve"
)

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
		Id     string `json:"id"`
		Err    string `json:"errmsg"`
	} `json:"head"`
}

type CallBack struct {
	resp Response
	err  error
}

type CallArg struct {
	request Request
	ret     chan CallBack
}

type IMessage struct {
	url     string
	org     string
	send    chan *CallArg
	cancel  chan *CallArg
	hash    *consistent.Consistent
	valves  map[string]*valve.Valve
	timeout time.Duration
}

func New(config *model.Config) (*IMessage, error) {
	hash := consistent.New()
	valves := make(map[string]*valve.Valve)
	for _, k := range config.Thirdpart.IMessage.Channels {
		hash.Add(k)
		valves[k] = valve.New(config.Thirdpart.IMessage.QueueDepth, time.Duration(config.Thirdpart.IMessage.PeriodInSecond)*time.Second)
	}
	ret := &IMessage{
		url:     config.Thirdpart.IMessage.Address,
		org:     config.Thirdpart.IMessage.Origin,
		send:    make(chan *CallArg),
		cancel:  make(chan *CallArg),
		hash:    hash,
		valves:  valves,
		timeout: broker.NetworkTimeout,
	}
	go ret.Serve()
	return ret, nil
}

func (im *IMessage) Provider() string {
	return "imessage"
}

func (im *IMessage) Serve() {
	for {
		sio, err := socketio.Dial(im.url, im.org, broker.NetworkTimeout)
		if err != nil {
			logger.ERROR("can't connect imessage server: %s", err)
			select {
			case arg := <-im.send:
				arg.ret <- CallBack{err: err}
			case <-time.After(im.timeout):
			}
			continue
		}
		syncId := make(map[string]*CallArg)
		eventId := make(map[string]*CallArg)
		for {
			select {
			case arg := <-im.send:
				id, err := sio.Emit(true, "send", arg.request)
				if err != nil {
					arg.ret <- CallBack{err: err}
				} else {
					syncId[fmt.Sprintf("%d+", id)] = arg
				}
			case arg := <-im.cancel:
				for k, v := range syncId {
					if v == arg {
						delete(syncId, k)
					}
				}
				for k, v := range eventId {
					if v == arg {
						delete(eventId, k)
					}
				}
			default:
			}
			var msg socketio.Message
			err := sio.Receive(&msg)
			if err != nil {
				if e, ok := err.(net.Error); !ok || !e.Timeout() {
					logger.ERROR("imessage breaked:", err)
					break
				}
			}
			switch msg.Type() {
			case socketio.MessageACK:
				if arg, ok := syncId[msg.EndPoint()]; ok {
					delete(syncId, msg.EndPoint())
					var ack Response
					err := msg.ReadArguments(&ack)
					if err != nil {
						arg.ret <- CallBack{err: err}
					} else {
						eventId[ack.Head.Id] = arg
					}
				}
			case socketio.MessageEvent:
				var data string
				err := msg.ReadArguments(&data)
				if err == nil {
					var ack Response
					err = json.Unmarshal([]byte(data), &ack)
					if err == nil {
						if arg, ok := eventId[ack.Head.Id]; ok {
							delete(eventId, ack.Head.Id)
							arg.ret <- CallBack{resp: ack}
						}
					}
				}
			}
		}
	}
}

func (im *IMessage) Check(to string) (bool, error) {
	channel, err := im.hash.Get(to)
	if err != nil {
		return false, err
	}
	call := CallArg{
		request: Request{
			To:      to,
			Channel: channel,
			Action:  "1",
		},
		ret: make(chan CallBack, 1),
	}
	im.send <- &call
	select {
	case back := <-call.ret:
		if back.err != nil {
			return false, back.err
		}
		return back.resp.Head.Status == 0, nil
	case <-time.After(im.timeout):
		im.cancel <- &call
	}
	return false, fmt.Errorf("imessage check timeout")
}

func (im *IMessage) Send(to, text string) (string, error) {
	channel, err := im.hash.Get(to)
	if err != nil {
		return "", err
	}
	call := CallArg{
		request: Request{
			To:      to,
			Channel: channel,
			Action:  "2",
			Message: text,
		},
		ret: make(chan CallBack, 1),
	}
	im.send <- &call
	select {
	case back := <-call.ret:
		if back.err != nil {
			return "", back.err
		}
		return back.resp.Head.Id, nil
	case <-time.After(im.timeout):
		im.cancel <- &call
	}
	return "", fmt.Errorf("imessage check timeout")
}

func (im *IMessage) Post(to, text string) (string, error) {
	return im.Send(to, text)
}
