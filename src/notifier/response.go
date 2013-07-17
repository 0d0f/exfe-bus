package notifier

import (
	"broker"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"logger"
	"model"
	"thirdpart"
)

type ResponseItem struct {
	DefaultOk bool
	OkAction  struct {
		TargetUrl string
		Arg       interface{}
	}
	FailAction struct {
		TargetUrl string
		Arg       interface{}
	}
}

type ResponseSaver interface {
	Save(id string, item ResponseItem, ontime int64) error
	Load(id string) (ResponseItem, error)
}

type Response struct {
	saver  ResponseSaver
	config *model.Config
}

func NewResponse(config *model.Config) (*Response, error) {
	return &Response{
		config: config,
	}, nil
}

func (r *Response) WaitResponse(id string, ontime int64, defaultOk bool, identityId, u string, args interface{}) {
	if ontime == 0 && defaultOk {
		return
	}
	item := ResponseItem{
		DefaultOk: defaultOk,
	}
	item.OkAction.TargetUrl = fmt.Sprintf("%s/v3/bus/notificatioincallback")
	item.OkAction.Arg = map[string]interface{}{
		"identity_id": identityId,
	}
	item.FailAction.TargetUrl = u
	item.FailAction.Arg = args
	err := r.saver.Save(id, item, ontime)
	if err != nil {
		logger.ERROR("save response item(%d) failed: %s", id, err)
		return
	}
	if defaultOk {
		r.PushQueue(item.OkAction.TargetUrl, item.OkAction.Arg, ontime)
	} else {
		r.PushQueue(item.FailAction.TargetUrl, item.FailAction.Arg, ontime)
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
		if item.DefaultOk {
			if !resp.Ok {
				r.DeleteQueue(item.OkAction.TargetUrl)
				r.Do(item.FailAction.TargetUrl, item.FailAction.Arg)
			}
		} else {
			if resp.Ok {
				r.DeleteQueue(item.FailAction.TargetUrl)
				r.Do(item.OkAction.TargetUrl, item.OkAction.Arg)
			}
		}
	}
}

func (r *Response) PushQueue(u string, arg interface{}, ontime int64) {
	queueUrl = fmt.Sprintf("http://%s:%d/v3/queue/-/POST/%s?ontime=%d",
		r.config.ExfeQueue.Addr, r.config.ExfeQueue.Port, base64.URLEncoding.EncodeToString(u), ontime)
}

func (r *Response) DeleteQueue(u string) {
}

func (r *Response) Do(u string, args interface{}) {
}
