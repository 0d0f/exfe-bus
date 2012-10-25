package delayrepo

import (
	"broker"
	"fmt"
	"time"
)

type timeSaver interface {
	SaveTime(time int64, key string) error
}

type delay struct {
	redis broker.Redis

	delayInSecond int
	redisPrefix   string
	timeHashName  string
	timeSaver     timeSaver
}

func (q *delay) Push(key string, data []byte) (err error) {
	err = q.pushWithTime(time.Now().Unix(), key, data)
	return
}

func (q *delay) Pop() (string, [][]byte, error) {
	key, err := q.getTimeoutKey()
	if err == EmptyError {
		return "", nil, nil
	}
	datas, err := q.popFromKey(key)
	if err != nil {
		return "", nil, err
	}
	return key, datas, nil
}

func (q *delay) NextWakeup() (time.Duration, error) {
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
	return time.Duration(int64(score)+int64(q.delayInSecond)-time.Now().Unix()) * time.Second, nil
}

func (q *delay) initDelay(name string, delayInSecond int, redis broker.Redis, timeSaver timeSaver) {
	q.redis = redis
	q.delayInSecond = delayInSecond
	q.redisPrefix = fmt.Sprintf("gobus:queue:%s", name)
	q.timeHashName = fmt.Sprintf("%s:timehash", q.redisPrefix)
	q.timeSaver = timeSaver
}

func (q *delay) pushWithTime(time int64, key string, data []byte) error {
	err := q.timeSaver.SaveTime(time, key)
	if err != nil {
		return err
	}

	_, err = q.redis.Rpush(fmt.Sprintf("%s:%s", q.redisPrefix, key), data)
	if err != nil {
		return err
	}
	return nil
}

func (q *delay) getTimeoutKey() (string, error) {
	end := time.Now().Unix() - int64(q.delayInSecond)

	reply, err := q.redis.Zrangebyscore(q.timeHashName, "0", fmt.Sprintf("%d", end))
	if err != nil {
		return "", err
	}

	if len(reply.Elems) == 0 {
		return "", EmptyError
	}

	key := string(reply.Elems[0].Elem)
	_, err = q.redis.Zrem(q.timeHashName, key)
	if err != nil {
		return key, err
	}
	return key, nil
}

func (q *delay) popFromKey(key string) ([][]byte, error) {
	queueName := fmt.Sprintf("%s:%s", q.redisPrefix, key)

	pipe := q.redis.NewPipeClient()
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
		return nil, ChangedError
	}

	ret := make([][]byte, len(r[0].Elems))
	for i, _ := range r[0].Elems {
		ret[i] = r[0].Elems[i].Elem
	}
	return ret, nil
}

//////////////////////////////

type Head struct {
	delay
}

func NewHead(name string, delayInSecond int, redis broker.Redis) *Head {
	ret := new(Head)
	ret.initDelay(name, delayInSecond, redis, ret)
	return ret
}

func (q *Head) SaveTime(time int64, key string) error {
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
