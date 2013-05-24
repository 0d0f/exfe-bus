package main

import (
	"broker"
	"delayrepo"
	"encoding/base64"
	"fmt"
	"github.com/googollee/go-rest"
	"io"
	"io/ioutil"
	"logger"
	"model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	rest.RegisterMarshaller("plain/text", new(PlainText))
}

type Queue struct {
	rest.Service `prefix:"/v3/queue" mime:"plain/text"`

	config  *model.Config
	timeout time.Duration

	Push   rest.Processor `path:"/:merge_key/:method/*service" method:"POST"`
	Delete rest.Processor `path:"/:merge_key/:method/*service" method:"DELETE"`
	timer  *delayrepo.Timer
}

func NewQueue(config *model.Config, redis *broker.RedisPool) (*Queue, error) {
	ret := &Queue{
		config:  config,
		timeout: time.Second * 30,
	}

	logger.NOTICE("launching timer")
	storage := broker.NewQueueRedisStorage("exfe:v3:queue", redis)
	timer, err := delayrepo.NewTimer(storage, ret.timeout)
	if err != nil {
		return nil, err
	}
	ret.timer = timer
	go timer.Serve(ret)

	return ret, nil
}

func (q *Queue) Do(key string, datas [][]byte) {
	fl := logger.FUNC(key)
	defer fl.Quit()

	splits := strings.Split(key, ",")
	if len(splits) != 3 {
		logger.ERROR("pop error key: %s", key)
		return
	}
	method, service, mergeKey := splits[0], splits[1], splits[2]
	needMerge := mergeKey[0] != '-'

	args := []byte("[")
	for _, data := range datas {
		if needMerge {
			args = append(args, data...)
			args = append(args, []byte(",")...)
		} else {
			go func() {
				resp, err := broker.Http(method, service, "application/json", data)
				if err != nil {
					logger.ERROR("%s %s: %s, with %s", method, service, err, string(data))
				} else {
					resp.Body.Close()
				}
			}()
		}
	}
	if needMerge && len(args) > 1 {
		args[len(args)-1] = byte(']')
		go func() {
			resp, err := broker.Http(method, service, "application/json", args)
			if err != nil {
				logger.ERROR("%s %s: %s, with %s", method, service, err, string(args))
			} else {
				resp.Body.Close()
			}
		}()
	}
}

func (q *Queue) OnError(err error) {
	logger.ERROR("queue error: %s", err)
}

func (q *Queue) Quit() {
	logger.NOTICE("kill timer")
	q.timer.Quit()
}

// example:
// POST to bus://exfe_service/message with merge_key 123, always send on 1366615888, data is {"abc":123}
// > curl -v "http://127.0.0.1:23334/v3/queue/123/POST/exfe_service/message?update=always&ontime=1366615888" -d '{"abc":123}'
//
// if no merge(send one by one), set merge_key to "-"
func (q Queue) HandlePush(data string) {
	method, service, mergeKey := q.Vars()["method"], q.Vars()["service"], q.Vars()["merge_key"]
	if method == "" {
		q.Error(http.StatusBadRequest, q.GetError(1, "need method"))
		return
	}
	if service == "" {
		q.Error(http.StatusBadRequest, q.GetError(2, "need service"))
		return
	}
	if mergeKey == "" {
		q.Error(http.StatusBadRequest, q.GetError(3, "invalid mergeKey: (empty)"))
		return
	}
	b, err := base64.URLEncoding.DecodeString(mergeKey)
	if err == nil {
		mergeKey = string(b)
	}
	b, err = base64.URLEncoding.DecodeString(service)
	if err != nil {
		q.Error(http.StatusBadRequest, q.GetError(4, fmt.Sprintf("service(%s) invalid: %s", service, err)))
		return
	}
	service = string(b)
	fl := logger.FUNC(method, service, mergeKey)
	defer fl.Quit()

	query := q.Request().URL.Query()
	updateType, ontimeStr := query.Get("update"), query.Get("ontime")
	ontime, err := strconv.ParseInt(ontimeStr, 10, 64)
	if err != nil {
		q.Error(http.StatusBadRequest, q.GetError(5, fmt.Sprintf("invalid ontime: %s", ontimeStr)))
		return
	}
	if ontime == 0 {
		ontime = time.Now().Unix()
	}

	err = q.timer.Push(delayrepo.UpdateType(updateType), ontime, fmt.Sprintf("%s,%s,%s", method, service, mergeKey), []byte(data))
	if err != nil {
		q.Error(http.StatusInternalServerError, q.GetError(7, err.Error()))
		return
	}
}

func (q Queue) HandleDelete() {
	method, service, mergeKey := q.Vars()["method"], q.Vars()["service"], q.Vars()["merge_key"]
	if method == "" {
		q.Error(http.StatusBadRequest, q.GetError(1, "need method"))
		return
	}
	if service == "" {
		q.Error(http.StatusBadRequest, q.GetError(2, "need service"))
		return
	}
	if mergeKey == "" {
		q.Error(http.StatusBadRequest, q.GetError(3, "invalid mergeKey: (empty)"))
		return
	}
	b, err := base64.URLEncoding.DecodeString(mergeKey)
	if err == nil {
		mergeKey = string(b)
	}
	b, err = base64.URLEncoding.DecodeString(service)
	if err != nil {
		q.Error(http.StatusBadRequest, q.GetError(4, fmt.Sprintf("service(%s) invalid: %s", service, err)))
		return
	}
	service = string(b)
	fl := logger.FUNC(method, service, mergeKey)
	defer fl.Quit()

	err = q.timer.Delete(fmt.Sprintf("%s,%s,%s", method, service, mergeKey))
	if err != nil {
		q.Error(http.StatusInternalServerError, q.GetError(7, err.Error()))
		return
	}
}

type PlainText struct{}

func (p PlainText) Unmarshal(r io.Reader, v interface{}) error {
	ps, ok := v.(*string)
	if !ok {
		return fmt.Errorf("plain text only can save in string")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	*ps = string(b)
	return nil
}

func (p PlainText) Marshal(w io.Writer, v interface{}) error {
	_, err := w.Write([]byte(fmt.Sprintf("%s", v)))
	return err
}

type TextError string

func (t TextError) Error() string {
	return string(t)
}

func (p PlainText) Error(code int, message string) error {
	return TextError(fmt.Sprintf("(%d)%s", code, message))
}
