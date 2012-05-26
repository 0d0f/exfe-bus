package gobus

import (
	"github.com/simonz05/godis"
	"fmt"
	"errors"
	"time"
	"encoding/json"
	"bytes"
	"reflect"
)

var EmptyQueueError = errors.New("Empty queue.")
var QueueChangedError = errors.New("Queue changed before pop")

type TailDelayQueue struct {
	redis *godis.Client
	name string
	hashName string
	delay int
	dataType reflect.Type
}

func NewTailDelayQueue(name string, delayInSeconds int, typeInstance interface{}, redis *godis.Client) (*TailDelayQueue, error) {
	t := reflect.TypeOf(typeInstance)
	if t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("typeInstance must be a slice")
	}
	return &TailDelayQueue{
		redis: redis,
		name: name,
		hashName: fmt.Sprintf("%s:timehash", name),
		delay: delayInSeconds,
		dataType: t,
	}, nil
}

func (q *TailDelayQueue) Push(id string, data interface{}) error {
	return q.PushWithTime(float64(time.Now().Unix()), id, data)
}

func (q *TailDelayQueue) PushWithTime(time float64, id string, data interface{}) error {
	score := time
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}

	_, err = q.redis.Zadd(q.hashName, score, id)
	if err != nil {
		return err
	}
	_, err = q.redis.Rpush(fmt.Sprintf("%s:%s", q.name, id), buf.String())
	if err != nil {
		return err
	}
	return nil
}

func (q *TailDelayQueue) GetTimeoutId() (string, error) {
	end := time.Now().Unix() - int64(q.delay)

	reply, err := q.redis.Zrangebyscore(q.hashName, "0", fmt.Sprintf("%d", end))
	if err != nil {
		return "", err
	}

	if len(reply.Elems) == 0 {
		return "", EmptyQueueError
	}

	id := string(reply.Elems[0].Elem)
	_, err = q.redis.Zrem(q.hashName, id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (q *TailDelayQueue) PopFromId(id string) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s", q.name, id)

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
	ret := reflect.MakeSlice(q.dataType, 0, 0)
	for _, reply := range r[0].Elems {
		buf := bytes.NewBuffer(reply.Elem)
		decoder := json.NewDecoder(buf)
		data := reflect.New(q.dataType.Elem())
		err := decoder.Decode(data.Interface())
		if err != nil {
			continue
		}
		ret = reflect.Append(ret, data.Elem())
	}
	return ret.Interface(), nil
}

func (q *TailDelayQueue) Pop() (interface{}, error) {
	id, err := q.GetTimeoutId()
	if err == EmptyQueueError {
		return nil, nil
	}
	return q.PopFromId(id)
}

func (q *TailDelayQueue) NextWakeup() (time.Duration, error) {
	reply, err := q.redis.Zrange(q.hashName, 0, -1)
	if err != nil {
		return time.Duration(q.delay) * time.Second, err
	}
	if len(reply.Elems) == 0 {
		return time.Duration(q.delay) * time.Second, nil
	}
	id := string(reply.Elems[0].Elem)
	score, err := q.redis.Zscore(q.hashName, id)
	if err != nil {
		return time.Duration(q.delay) * time.Second, err
	}
	return time.Duration(int64(score) + int64(q.delay) - time.Now().Unix()) * time.Second, nil
}

//////////////////////////////////

type HeadDelayQueue struct {
	redis *godis.Client
	name string
	hashName string
	delay int
	dataType reflect.Type
}

func NewHeadDelayQueue(name string, delayInSeconds int, typeInstance interface{}, redis *godis.Client) (*HeadDelayQueue, error) {
	t := reflect.TypeOf(typeInstance)
	if t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("typeInstance must be a slice")
	}
	return &HeadDelayQueue{
		redis: redis,
		name: name,
		hashName: fmt.Sprintf("%s:timehash", name),
		delay: delayInSeconds,
		dataType: t,
	}, nil
}

func (q *HeadDelayQueue) Push(id string, data interface{}) error {
	return q.PushWithTime(float64(time.Now().Unix()), id, data)
}

func (q *HeadDelayQueue) PushWithTime(time float64, id string, data interface{}) error {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}

	score, _ := q.redis.Zscore(q.hashName, id)
	if score < 0 {
		score = time
		_, err = q.redis.Zadd(q.hashName, score, id)
		if err != nil {
			return err
		}
	}
	_, err = q.redis.Rpush(fmt.Sprintf("%s:%s", q.name, id), buf.String())
	if err != nil {
		return err
	}
	return nil
}

func (q *HeadDelayQueue) GetTimeoutId() (string, error) {
	end := time.Now().Unix() - int64(q.delay)

	reply, err := q.redis.Zrangebyscore(q.hashName, "0", fmt.Sprintf("%d", end))
	if err != nil {
		return "", err
	}

	if len(reply.Elems) == 0 {
		return "", EmptyQueueError
	}

	id := string(reply.Elems[0].Elem)
	_, err = q.redis.Zrem(q.hashName, id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (q *HeadDelayQueue) PopFromId(id string) (interface{}, error) {
	queueName := fmt.Sprintf("%s:%s", q.name, id)

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
	ret := reflect.MakeSlice(q.dataType, 0, 0)
	for _, reply := range r[0].Elems {
		buf := bytes.NewBuffer(reply.Elem)
		decoder := json.NewDecoder(buf)
		data := reflect.New(q.dataType.Elem())
		err := decoder.Decode(data.Interface())
		if err != nil {
			continue
		}
		ret = reflect.Append(ret, data.Elem())
	}
	return ret.Interface(), nil
}

func (q *HeadDelayQueue) Pop() (interface{}, error) {
	id, err := q.GetTimeoutId()
	if err == EmptyQueueError {
		return nil, nil
	}
	return q.PopFromId(id)
}

func (q *HeadDelayQueue) NextWakeup() (time.Duration, error) {
	reply, err := q.redis.Zrange(q.hashName, 0, -1)
	if err != nil {
		return time.Duration(q.delay) * time.Second, err
	}
	if len(reply.Elems) == 0 {
		return time.Duration(q.delay) * time.Second, nil
	}
	id := string(reply.Elems[0].Elem)
	score, err := q.redis.Zscore(q.hashName, id)
	if err != nil {
		return time.Duration(q.delay) * time.Second, err
	}
	return time.Duration(int64(score) + int64(q.delay) - time.Now().Unix()) * time.Second, nil
}
