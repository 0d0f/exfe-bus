package main

import (
	"broker"
	"delayrepo"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-rest"
	"gobus"
	"model"
	"net/http"
	"strings"
	"time"
)

type Queue struct {
	rest.Service `prefix:"/v3/queue"`

	config     *model.Config
	timeout    time.Duration
	dispatcher *gobus.Dispatcher

	Timer rest.Processor `path:"/:merge_key/:method/*service" method:"POST"`
	timer *delayrepo.Timer
}

func NewQueue(config *model.Config, redis *broker.RedisPool, dispatcher *gobus.Dispatcher) (*Queue, error) {
	ret := &Queue{
		config:     config,
		timeout:    time.Second * 30,
		dispatcher: dispatcher,
	}

	config.Log.Notice("launching timer")
	storage := broker.NewQueueRedisStorage("exfe:v3:queue", redis)
	timer, err := delayrepo.NewTimer(storage)
	if err != nil {
		return nil, err
	}
	ret.timer = timer
	go timer.Serve(ret, ret.timeout)

	return ret, nil
}

func (q *Queue) Do(key string, data [][]byte) {
	splits := strings.Split(key, ",")
	if len(splits) != 3 {
		q.config.Log.Err("pop error key: %s", key)
		return
	}
	method, service, mergeKey := splits[0], "bus://"+splits[1], splits[2]

	args := make([]interface{}, 0)
	for _, d := range data {
		var arg interface{}
		err := json.Unmarshal(d, &arg)
		if err != nil {
			q.config.Log.Err("can't unmarshal %s(%+v)", err, d)
			continue
		}
		if mergeKey != "-" {
			args = append(args, arg)
		} else {
			var reply interface{}
			err := q.dispatcher.Do(service, method, arg, &reply)
			if err != nil {
				j, _ := json.Marshal(arg)
				q.config.Log.Err("call %s|%s failed(%s) with %s", service, method, err, string(j))
			}
		}
	}
	if mergeKey != "-" {
		var reply interface{}
		err := q.dispatcher.Do(service, method, args, &reply)
		if err != nil {
			j, _ := json.Marshal(args)
			q.config.Log.Err("call %s|%s failed(%s) with %s", service, method, err, string(j))
		}
	}
}

func (q *Queue) OnError(err error) {
	q.config.Log.Crit("queue error: %s", err)
}

func (q *Queue) Quit() {
	q.config.Log.Notice("kill timer")
	q.timer.Quit()
}

type QueueData struct {
	Type   delayrepo.UpdateType `json:"type"`
	Ontime int64                `json:"ontime"`
	Data   interface{}          `json:"data"`
}

// example:
// POST to bus://exfe_service/message with merge_key 123, always send on 1366615888, data is {"abc":123}
// > curl -v "http://127.0.0.1:23334/v3/queue/123/POST/exfe_service/message" -d '{"type":"always","ontime":1366615888,"data":{"abc":123}}'
//
// if no merge(send one by one), set merge_key to "-"
func (q Queue) HandleTimer(push QueueData) {
	method, service, mergeKey := q.Vars()["method"], q.Vars()["service"], q.Vars()["merge_key"]
	if method == "" {
		q.Error(http.StatusBadRequest, fmt.Errorf("need method"))
		return
	}
	if service == "" {
		q.Error(http.StatusBadRequest, fmt.Errorf("need service"))
		return
	}

	buf, err := json.Marshal(push.Data)
	if err != nil {
		q.Error(http.StatusBadRequest, err)
		return
	}
	err = q.timer.Push(push.Type, push.Ontime, fmt.Sprintf("%s,%s,%s", method, service, mergeKey), buf)
	if err != nil {
		q.Error(http.StatusInternalServerError, err)
		return
	}
}
