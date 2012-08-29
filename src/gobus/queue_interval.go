package gobus

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/godis"
	"reflect"
	"time"
)

type IntervalQueue struct {
	redis *godis.Client

	lastWakeup    time.Time
	delayInSecond int64
	redisPrefix   string
	sliceType     reflect.Type
}

func NewIntervalQueue(name string, delayInSecond int64, redis *godis.Client, sliceType interface{}) *IntervalQueue {
	return &IntervalQueue{
		redis:         redis,
		lastWakeup:    time.Now(),
		delayInSecond: delayInSecond,
		redisPrefix:   fmt.Sprintf("gobus:queue:%s", name),
		sliceType:     reflect.TypeOf(sliceType),
	}
}

func (q *IntervalQueue) Push(data QueueData) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = q.redis.Rpush(q.redisPrefix, string(buf))
	if err != nil {
		return err
	}
	return nil
}

func (q *IntervalQueue) Pop() (interface{}, error) {
	pipe := godis.NewPipeClientFromClient(q.redis)
	defer pipe.Quit()

	err := pipe.Watch(q.redisPrefix)
	if err != nil {
		return nil, err
	}
	err = pipe.Multi()
	if err != nil {
		return nil, err
	}
	_, err = pipe.Lrange(q.redisPrefix, 0, -1)
	if err != nil {
		return nil, err
	}
	_, err = pipe.Del(q.redisPrefix)
	if err != nil {
		return nil, err
	}
	r := pipe.Exec()
	if len(r) == 0 {
		return nil, QueueChangedError
	}

	ret := reflect.MakeSlice(q.sliceType, 0, 0)
	for _, reply := range r[0].Elems {
		data := reflect.New(q.sliceType.Elem())
		err := json.Unmarshal(reply.Elem, data.Interface())
		if err != nil {
			continue
		}
		ret = reflect.Append(ret, data.Elem())
	}
	q.lastWakeup = time.Now()
	return ret.Interface(), nil
}

func (q *IntervalQueue) NextWakeup() (time.Duration, error) {
	d := time.Since(q.lastWakeup)
	d = time.Duration(q.delayInSecond)*time.Second - d
	if d < 0 {
		return 0, nil
	}
	return d, nil
}
