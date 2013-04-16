package main

import (
	"broker"
	"delayrepo"
	"encoding/json"
	"github.com/googollee/go-rest"
	"gobus"
	"model"
	"net/http"
	"time"
)

type Queue struct {
	rest.Service `prefix:"/v3/queue"`

	config     *model.Config
	timeout    time.Duration
	dispatcher *gobus.Dispatcher

	Timer rest.Processor `path:"/timer" method:"POST"`
	timer *delayrepo.Repository
}

func NewQueue(config *model.Config, redis *broker.RedisPool, dispatcher *gobus.Dispatcher) (*Queue, error) {
	ret := &Queue{
		config:     config,
		timeout:    time.Second * 30,
		dispatcher: dispatcher,
	}

	config.Log.Notice("launching timer")
	ret.timer = delayrepo.New(delayrepo.NewTimer("bus:queue", redis), ret, ret.timeout)
	go ret.timer.Serve()

	return ret, nil
}

func (q *Queue) Do(key string, data [][]byte) {
	service, method, mergeKey, _, err := model.QueueParseKey(key)
	if err != nil {
		q.config.Log.Err("pop error: %s")
	}

	arg := make([]interface{}, 0)
	for _, d := range data {
		var i interface{}
		err := json.Unmarshal(d, &i)
		if err != nil {
			q.config.Log.Err("can't unmarshal %s(%+v)", err, d)
			continue
		}
		if mergeKey != "" {
			arg = append(arg, d)
		} else {
			var i interface{}
			err := q.dispatcher.Do(service, method, d, &i)
			if err != nil {
				j, _ := json.Marshal(arg)
				q.config.Log.Err("call %s|%s failed(%s) with %s", service, method, err, string(j))
			}
		}
	}
	if mergeKey != "" {
		var i interface{}
		err := q.dispatcher.Do(service, "POST", arg, &i)
		if err != nil {
			j, _ := json.Marshal(arg)
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

func (q Queue) HandlerTimer(push model.QueuePush) {
	err := push.Init(q.config.ExfeQueue.Priority)
	if err != nil {
		q.Error(http.StatusBadRequest, err)
		return
	}
	for data := range push.Each() {
		buf, err := json.Marshal(data.Data)
		if err != nil {
			q.Error(http.StatusBadRequest, err)
		}
		err = q.timer.Push(push.Ontime, data.Key, buf)
		if err != nil {
			q.Error(http.StatusInternalServerError, err)
		}
	}
}
