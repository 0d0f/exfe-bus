package main

import (
	"broker"
	"delayrepo"
	"encoding/base64"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-rest"
	"logger"
	"model"
	"net/http"
	"strings"
	"time"
)

func init() {
	rest.RegisterMarshaller("plain/text", new(broker.PlainText))
}

type Queue struct {
	rest.Service `prefix:"/v3/queue" mime:"plain/text"`

	config  *model.Config
	timeout time.Duration

	push   rest.SimpleNode `route:"/:merge_key/:method/*service" method:"POST"`
	delete rest.SimpleNode `route:"/:merge_key/:method/*service" method:"DELETE"`
	timer  *delayrepo.Timer
}

func NewQueue(config *model.Config, redis *redis.Pool) (*Queue, error) {
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
				resp, err := broker.HttpResponse(broker.Http(method, service, "application/json", data))
				if err != nil {
					logger.ERROR("%s %s: %s, with %s", method, service, err, string(data))
				} else {
					resp.Close()
					logger.INFO("queue", "do", method, service, mergeKey)
				}
			}()
		}
	}
	if needMerge && len(args) > 1 {
		args[len(args)-1] = byte(']')
		go func() {
			resp, err := broker.HttpResponse(broker.Http(method, service, "application/json", args))
			if err != nil {
				logger.ERROR("%s %s: %s, with %s", method, service, err, string(args))
			} else {
				resp.Close()
				logger.INFO("queue", "do", method, service, mergeKey)
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
func (q Queue) Push(ctx rest.Context, data string) {
	var method, service, mergeKey, updateType string
	var ontime int64
	ctx.Bind("method", &method)
	ctx.Bind("service", &service)
	ctx.Bind("merge_key", &mergeKey)
	ctx.Bind("update", &updateType)
	ctx.Bind("ontime", &ontime)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	if method == "" {
		ctx.Return(http.StatusBadRequest, "need method")
		return
	}
	if service == "" {
		ctx.Return(http.StatusBadRequest, "need service")
		return
	}
	if mergeKey == "" {
		ctx.Return(http.StatusBadRequest, "invalid mergeKey: (empty)")
		return
	}
	b, err := base64.URLEncoding.DecodeString(mergeKey)
	if err == nil {
		mergeKey = string(b)
	}
	b, err = base64.URLEncoding.DecodeString(service)
	if err != nil {
		ctx.Return(http.StatusBadRequest, "service(%s) invalid: %s", service, err)
		return
	}
	service = string(b)
	if updateType == "" {
		updateType = "once"
	}

	if ontime == 0 {
		ontime = time.Now().Unix()
	}

	fl := logger.FUNC(method, service, mergeKey, updateType, ontime)
	defer fl.Quit()

	err = q.timer.Push(delayrepo.UpdateType(updateType), ontime, fmt.Sprintf("%s,%s,%s", method, service, mergeKey), []byte(data))
	if err != nil {
		ctx.Return(http.StatusInternalServerError, err)
		return
	}
	logger.INFO("queue", "push", method, service, mergeKey, ontime)
}

func (q Queue) Delete(ctx rest.Context) {
	var method, service, mergeKey string
	ctx.Bind("method", &method)
	ctx.Bind("service", &service)
	ctx.Bind("merge_key", &mergeKey)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, err)
		return
	}
	if method == "" {
		ctx.Return(http.StatusBadRequest, "need method")
		return
	}
	if service == "" {
		ctx.Return(http.StatusBadRequest, "need service")
		return
	}
	if mergeKey == "" {
		ctx.Return(http.StatusBadRequest, "invalid mergeKey: (empty)")
		return
	}
	b, err := base64.URLEncoding.DecodeString(mergeKey)
	if err == nil {
		mergeKey = string(b)
	}
	b, err = base64.URLEncoding.DecodeString(service)
	if err != nil {
		ctx.Return(http.StatusBadRequest, "service(%s) invalid: %s", service, err)
		return
	}
	service = string(b)
	fl := logger.FUNC(method, service, mergeKey)
	defer fl.Quit()

	err = q.timer.Delete(fmt.Sprintf("%s,%s,%s", method, service, mergeKey))
	if err != nil {
		ctx.Return(http.StatusInternalServerError, err)
		return
	}
	logger.INFO("queue", "delete", method, service, mergeKey)
}
