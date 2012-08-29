package gobus

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/godis"
	"reflect"
	"time"
)

type timeSaver interface {
	SaveTime(time int64, key string) error
}

type delayQueue struct {
	redis *godis.Client

	delayInSecond int64
	redisPrefix   string
	timeHashName  string
	sliceType     reflect.Type
	timeSaver     timeSaver
}

func (q *delayQueue) Push(data QueueData) (err error) {
	key := data.KeyForQueue()
	err = q.pushWithTime(time.Now().Unix(), key, data)
	return
}

func (q *delayQueue) Pop() (interface{}, error) {
	key, err := q.getTimeoutKey()
	if err == EmptyQueueError {
		return nil, nil
	}
	return q.popFromKey(key)
}

func (q *delayQueue) NextWakeup() (time.Duration, error) {
	reply, err := q.redis.Zrange(q.timeHashName, 0, -1)
	if err != nil {
		return time.Duration(q.delayInSecond) * time.Second, err
	}
	if len(reply.Elems) == 0 {
		return time.Duration(q.delayInSecond) * time.Second, nil
	}
	id := string(reply.Elems[0].Elem)
	score, err := q.redis.Zscore(q.timeHashName, id)
	if err != nil {
		return time.Duration(q.delayInSecond) * time.Second, err
	}
	return time.Duration(int64(score)+q.delayInSecond-time.Now().Unix()) * time.Second, nil
}

func (q *delayQueue) initDelayQueue(name string, delayInSecond int64, redis *godis.Client, timeSaver timeSaver, sliceType interface{}) {
	q.redis = redis
	q.delayInSecond = delayInSecond
	q.redisPrefix = fmt.Sprintf("gobus:queue:%s", name)
	q.timeHashName = fmt.Sprintf("%s:timehash", q.redisPrefix)
	q.sliceType = reflect.TypeOf(sliceType)
	q.timeSaver = timeSaver
}

func (q *delayQueue) pushWithTime(time int64, key string, data interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = q.timeSaver.SaveTime(time, key)
	if err != nil {
		return err
	}

	_, err = q.redis.Rpush(fmt.Sprintf("%s:%s", q.redisPrefix, key), string(buf))
	if err != nil {
		return err
	}
	return nil
}

func (q *delayQueue) getTimeoutKey() (string, error) {
	end := time.Now().Unix() - q.delayInSecond

	reply, err := q.redis.Zrangebyscore(q.timeHashName, "0", fmt.Sprintf("%d", end))
	if err != nil {
		return "", err
	}

	if len(reply.Elems) == 0 {
		return "", EmptyQueueError
	}

	key := string(reply.Elems[0].Elem)
	_, err = q.redis.Zrem(q.timeHashName, key)
	if err != nil {
		return key, err
	}
	return key, nil
}

func (q *delayQueue) popFromKey(key string) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s", q.redisPrefix, key)

	pipe := godis.NewPipeClientFromClient(q.redis)
	defer pipe.Quit()

	err := pipe.Watch(queueName)
	if err != nil {
		return nil, err
	}
	err = pipe.Multi()
	if err != nil {
		return nil, err
	}
	_, err = pipe.Lrange(queueName, 0, -1)
	if err != nil {
		return nil, err
	}
	_, err = pipe.Del(queueName)
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
	return ret.Interface(), nil
}

//////////////////////////////

type HeadDelayQueue struct {
	delayQueue
}

func NewHeadDelayQueue(name string, delayInSecond int64, redis *godis.Client, sliceType interface{}) *HeadDelayQueue {
	ret := new(HeadDelayQueue)
	ret.initDelayQueue(name, delayInSecond, redis, ret, sliceType)
	return ret
}

func (q *HeadDelayQueue) SaveTime(time int64, key string) error {
	score, _ := q.redis.Zscore(q.timeHashName, key)
	if score < 0 {
		score = float64(time)
		_, err := q.redis.Zadd(q.timeHashName, score, key)
		if err != nil {
			return err
		}
	}
	return nil
}

//////////////////////////////

type TailDelayQueue struct {
	delayQueue
}

func NewTailDelayQueue(name string, delayInSecond int64, redis *godis.Client, sliceType interface{}) *TailDelayQueue {
	ret := new(TailDelayQueue)
	ret.initDelayQueue(name, delayInSecond, redis, ret, sliceType)
	return ret
}

func (q *TailDelayQueue) SaveTime(time int64, key string) error {
	score := time
	_, err := q.redis.Zadd(q.timeHashName, score, key)
	if err != nil {
		return err
	}
	return nil
}
