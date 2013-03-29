package delayrepo

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Timer struct {
	redis    redis.Conn
	prefix   string
	timerKey string
	delay    time.Duration
}

func NewTimer(prefix string, delay time.Duration, redis redis.Conn) *Timer {
	return &Timer{
		redis:    redis,
		prefix:   prefix,
		timerKey: fmt.Sprintf("%s:timer", prefix),
		delay:    delay,
	}
}

func (t *Timer) Push(ontime int64, key, data interface{}) error {
	err := t.redis.Send("MULTI")
	if err != nil {
		return err
	}
	err = t.redis.Send("RPUSH", fmt.Sprintf("%s:storage:%s", t.prefix, key), data)
	if err != nil {
		return err
	}
	err = t.redis.Send("ZADD", t.timerKey, ontime, key)
	if err != nil {
		return err
	}
	_, err = t.redis.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

func (t *Timer) Pop() (string, []interface{}, error) {
	reply, err := redis.Values(t.redis.Do("ZRANGEBYSCORE", t.timerKey, "-inf", "+inf", "LIMIT", 0, 1, "WITHSCORES"))
	if err != nil {
		return "", nil, err
	}
	if len(reply) == 0 {
		return "", nil, nil
	}

	key, err := redis.String(reply[0], nil)
	if err != nil {
		return "", nil, err
	}
	ontime, err := redis.Int64(reply[1], nil)
	if err != nil {
		return "", nil, err
	}

	if ontime > time.Now().Unix() {
		return "", nil, nil
	}

	storageKey := fmt.Sprintf("%s:storage:%s", t.prefix, key)
	reply, err = redis.Values(t.redis.Do("LRANGE", storageKey, 0, -1))
	if err != nil {
		return "", nil, err
	}
	err = t.redis.Send("MULTI")
	if err != nil {
		return "", nil, err
	}
	err = t.redis.Send("ZREM", t.timerKey, key)
	if err != nil {
		return "", nil, err
	}
	err = t.redis.Send("DEL", storageKey)
	if err != nil {
		return "", nil, err
	}
	_, err = t.redis.Do("EXEC")
	if err != nil {
		return "", nil, err
	}

	return key, reply, nil
}

func (t *Timer) NextWakeup() (time.Duration, error) {
	reply, err := redis.Values(t.redis.Do("ZRANGEBYSCORE", t.timerKey, "-inf", "+inf", "LIMIT", 0, 1, "WITHSCORES"))
	if err != nil {
		return t.delay, err
	}
	if len(reply) == 0 {
		return t.delay, nil
	}
	ontime, err := redis.Int64(reply[1], nil)
	if err != nil {
		return t.delay, err
	}
	fmt.Println(ontime)
	fmt.Println(time.Unix(ontime, 0))
	return time.Unix(ontime, 0).Sub(time.Now()), nil
}
