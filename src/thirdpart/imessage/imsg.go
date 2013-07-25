package imessage

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/stathat/consistent"
	"logger"
	"model"
	"thirdpart"
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
		Stage  string `json:"stage"`
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
	hash    *consistent.Consistent
	valves  map[string]*valve.Valve
	timeout time.Duration
	f       thirdpart.Callback
}

func New(config *model.Config) (*IMessage, error) {
	hash := consistent.New()
	valves := make(map[string]*valve.Valve)
	for _, k := range config.Thirdpart.IMessage.Channels {
		hash.Add(k)
		valves[k] = valve.New(config.Thirdpart.IMessage.QueueDepth, time.Duration(config.Thirdpart.IMessage.PeriodInSecond)*time.Second)
		go valves[k].Serve()
	}
	ret := &IMessage{
		url:     config.Thirdpart.IMessage.Address,
		org:     config.Thirdpart.IMessage.Origin,
		send:    make(chan *CallArg),
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

func (im *IMessage) SetPosterCallback(callback thirdpart.Callback) (time.Duration, bool) {
	im.f = callback
	return 10, false
}

func (im *IMessage) Serve() {
	for {
		client, err := socketio.Dial(im.url, im.org)
		if err != nil {
			logger.ERROR("can't connect imessage server: %s", err)
			select {
			case arg := <-im.send:
				arg.ret <- CallBack{err: err}
			case <-time.After(im.timeout):
			}
			continue
		}
		client.On("send", func(ns *socketio.NameSpace, arg string) {
			var resp Response
			err := json.Unmarshal([]byte(arg), &resp)
			if err != nil {
				logger.ERROR("can't parse event send: %s", arg)
				return
			}
			if resp.Head.Stage == "SENT" || resp.Head.Status != 0 {
				var err error
				if resp.Head.Status != 0 {
					err = fmt.Errorf("%s", resp.Head.Err)
				}
				im.f(resp.Head.Id, err)
			}
		})

		// go client.Run()
		time.Sleep(time.Second / 10)
		for {
			arg := <-im.send
			var ret CallBack
			ret.err = client.Call("send", im.timeout, []interface{}{&ret.resp}, arg.request)
			arg.ret <- ret
		}
		client.Quit()
	}
}

type Work struct {
	call CallArg
	im   *IMessage
}

func (w Work) Do() (interface{}, error) {
	select {
	case w.im.send <- &w.call:
		back := <-w.call.ret
		if back.err != nil {
			return nil, back.err
		}
		return back.resp, nil
	case <-time.After(w.im.timeout):
	}
	return nil, fmt.Errorf("imessage timeout")
}

func (im *IMessage) Check(to string) (bool, error) {
	channel, err := im.hash.Get(to)
	if err != nil {
		return false, err
	}
	work := Work{
		call: CallArg{
			request: Request{
				To:      to,
				Channel: channel,
				Action:  "1",
			},
			ret: make(chan CallBack, 1),
		},
		im: im,
	}
	ret, err := im.valves[channel].Do(work)
	if err != nil {
		return false, err
	}
	resp, ok := ret.(Response)
	if !ok {
		return false, fmt.Errorf("not response")
	}
	return resp.Head.Status == 0, nil
}

func (im *IMessage) Send(to, text string) (string, error) {
	channel, err := im.hash.Get(to)
	if err != nil {
		return "", err
	}
	work := Work{
		call: CallArg{
			request: Request{
				To:      to,
				Channel: channel,
				Action:  "2",
				Message: text,
			},
			ret: make(chan CallBack, 1),
		},
		im: im,
	}
	ret, err := im.valves[channel].Do(work)
	if err != nil {
		return "", err
	}
	resp, ok := ret.(Response)
	if !ok {
		return "", fmt.Errorf("not response")
	}
	if resp.Head.Status != 0 {
		return "", fmt.Errorf("%s", resp.Head.Err)
	}
	return resp.Head.Id, nil
}

func (im *IMessage) Post(from, to, text string) (string, error) {
	return im.Send(to, text)
}
