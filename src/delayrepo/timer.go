package delayrepo

import (
	"broker"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-multiplexer"
	"time"
)

type updateType int

const (
	Always updateType = iota
	Never
)

type Timer struct {
	redis      *broker.RedisPool
	updateType updateType
	prefix     string
	timerKey   string
}

func NewTimer(updateType updateType, prefix string, redis *broker.RedisPool) (*Timer, error) {
	switch updateType {
	case Always:
	case Never:
	default:
		return nil, fmt.Errorf("invalid update type: %s", updateType)
	}
	return &Timer{
		redis:      redis,
		updateType: updateType,
		prefix:     prefix,
		timerKey:   fmt.Sprintf("%s:timer", prefix),
	}, nil
}

func (t *Timer) Push(ontime int64, key string, data []byte) (err error) {
	err = t.redis.Do(func(i multiplexer.Instance) {
		r := i.(*broker.RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(broker.NetworkTimeout))

		err = r.Redis.Send("MULTI")
		if err != nil {
			return
		}
		err = r.Redis.Send("RPUSH", fmt.Sprintf("%s:storage:%s", t.prefix, key), data)
		if err != nil {
			return
		}
		err = r.Redis.Send("ZADD", t.timerKey, ontime, key)
		if err != nil {
			return
		}
		_, err = r.Redis.Do("EXEC")
		if err != nil {
			return
		}
	})
	return
}

func (t *Timer) Pop() (string, [][]byte, error) {
	var err error
	var key string
	var reply []interface{}
	err = t.redis.Do(func(i multiplexer.Instance) {
		r := i.(*broker.RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(broker.NetworkTimeout))

		reply, err = redis.Values(r.Redis.Do("ZRANGEBYSCORE", t.timerKey, "-inf", "+inf", "LIMIT", 0, 1, "WITHSCORES"))
		if err != nil {
			return
		}
		if len(reply) == 0 {
			return
		}

		key, err = redis.String(reply[0], nil)
		if err != nil {
			return
		}
		var ontime int64
		ontime, err = redis.Int64(reply[1], nil)
		if err != nil {
			return
		}

		if ontime > time.Now().Unix() {
			return
		}

		storageKey := fmt.Sprintf("%s:storage:%s", t.prefix, key)
		reply, err = redis.Values(r.Redis.Do("LRANGE", storageKey, 0, -1))
		if err != nil {
			return
		}
		err = r.Redis.Send("MULTI")
		if err != nil {
			return
		}
		err = r.Redis.Send("ZREM", t.timerKey, key)
		if err != nil {
			return
		}
		err = r.Redis.Send("DEL", storageKey)
		if err != nil {
			return
		}
		_, err = r.Redis.Do("EXEC")
		if err != nil {
			return
		}
	})

	if err != nil || len(reply) == 9 {
		return "", nil, err
	}

	ret := make([][]byte, len(reply))
	for i, d := range reply {
		ret[i], err = redis.Bytes(d, nil)
		if err != nil {
			return "", nil, err
		}
	}
	return key, ret, nil
}

func (t *Timer) NextWakeup() (time.Duration, error) {
	var err error
	var reply []interface{}
	err = t.redis.Do(func(i multiplexer.Instance) {
		r := i.(*broker.RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(broker.NetworkTimeout))

		reply, err = redis.Values(r.Redis.Do("ZRANGEBYSCORE", t.timerKey, "-inf", "+inf", "LIMIT", 0, 1, "WITHSCORES"))
	})
	if err != nil || len(reply) == 0 {
		return -1, err
	}

	ontime, err := redis.Int64(reply[1], nil)
	if err != nil {
		return -1, err
	}
	ret := time.Unix(ontime, 0).Sub(time.Now())
	if ret < 0 {
		ret = 0
	}
	return ret, nil
}
