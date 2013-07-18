package notifier

import (
	"broker"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"logger"
	"model"
	"thirdpart"
	"time"
)

type ResponseItem struct {
	FailUrl string
	FailArg interface{}
}

type ResponseSaver interface {
	Save(id string, item ResponseItem, ontime int64) error
	Load(id string) (ResponseItem, error)
}

type Response struct {
	saver  ResponseSaver
	config *model.Config
}

var response *Response

func WaitResponse(id string, ontime int64, defaultOk bool, recipient model.Recipient, u string, args interface{}) {
	response.WaitResponse(id, ontime, defaultOk, recipient, u, args)
}

func SetupResponse(config *model.Config, saver ResponseSaver) error {
	response = &Response{
		config: config,
		saver:  saver,
	}
	go func() {
		for {
			err := response.Listen()
			if err != nil {
				logger.ERROR("response listen failed: %s", err)
			}
			time.Sleep(time.Second * 10)
		}
	}()
	return nil
}

func (r *Response) WaitResponse(id string, ontime int64, defaultOk bool, recipient model.Recipient, u string, args interface{}) {
	if ontime == 0 || len(recipient.Fallbacks) == 0 {
		return
	}
	item := ResponseItem{
		FailUrl: u,
		FailArg: args,
	}
	err := r.saver.Save(id, item, ontime)
	if err != nil {
		logger.ERROR("save response item(%d) failed: %s", id, err)
		return
	}
	if !defaultOk {
		r.PushQueue(id, item.FailUrl, item.FailArg, ontime)
	}
}

func (r *Response) Listen() error {
	listenUrl := fmt.Sprintf("http://%s:%d/v3/poster", r.config.ExfeService.Addr, r.config.ExfeService.Port)
	reader, err := broker.HttpResponse(broker.Http("WATCH", listenUrl, "application/json", nil))
	if err != nil {
		return err
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var resp thirdpart.PostResponse
	for {
		err := decoder.Decode(&resp)
		if err != nil {
			logger.ERROR("can't decode from post watch: %s", err)
			continue
		}
		item, err := r.saver.Load(resp.Id)
		if err != nil {
			logger.ERROR("can't load response item(%s): %s", resp.Id, err)
			continue
		}
		if resp.Ok {
			r.DeleteQueue(resp.Id, item.FailUrl)
		} else {
			r.Do(item.FailUrl, item.FailArg)
		}
	}
}

func (r *Response) PushQueue(id, u string, arg interface{}, ontime int64) {
	queueUrl := fmt.Sprintf("http://%s:%d/v3/queue/-%s/POST/%s?ontime=%d",
		r.config.ExfeQueue.Addr, r.config.ExfeQueue.Port, id, base64.URLEncoding.EncodeToString([]byte(u)), ontime)
	b, err := json.Marshal(arg)
	if err != nil {
		logger.ERROR("can't marshal: %s with %#v", err, arg)
		return
	}
	resp, err := broker.HttpResponse(broker.Http("POST", queueUrl, "text/plain", b))
	if err != nil {
		logger.ERROR("push to queue %s failed: %s with %s", queueUrl, err, string(b))
		return
	}
	resp.Close()
}

func (r *Response) DeleteQueue(id, u string) {
	queueUrl := fmt.Sprintf("http://%s:%d/v3/queue/-%s/POST/%s",
		r.config.ExfeQueue.Addr, r.config.ExfeQueue.Port, id, base64.URLEncoding.EncodeToString([]byte(u)))
	resp, err := broker.HttpResponse(broker.Http("DELETE", queueUrl, "text/plain", nil))
	if err != nil {
		logger.ERROR("delete queue %s failed: %s with %s", queueUrl, err)
		return
	}
	resp.Close()
}

func (r *Response) Do(u string, arg interface{}) {
	b, err := json.Marshal(arg)
	if err != nil {
		logger.ERROR("can't marshal: %s with %#v", err, arg)
		return
	}
	resp, err := broker.HttpResponse(broker.Http("POST", u, "text/plain", b))
	if err != nil {
		logger.ERROR("response do %s failed: %s with %s", u, err, string(b))
		return
	}
	resp.Close()
}
