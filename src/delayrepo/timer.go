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

type TimerStorage interface {
	Save(ontime int64, key string, data []byte) error
	Load(key string) ([][]byte, error)
	Ontime(key string) (int64, error)
	Next() (string, error)
}

type RedisStorage struct {
	redis    *broker.RedisPool
	prefix   string
	timerKey string
}

func NewTimerStorage(prefix string, redis *broker.RedisPool) *RedisStorage {
	return &RedisStorage{
		redis:    redis,
		prefix:   prefix,
		timerKey: fmt.Sprintf("%s:timer", prefix),
	}
}

func (s *RedisStorage) Save(ontime int64, key string, data []byte) (err error) {
	e := s.redis.Do(func(i multiplexer.Instance) {
		r := i.(*broker.RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(broker.NetworkTimeout))

		err = r.Redis.Send("MULTI")
		if err != nil {
			return
		}
		err = r.Redis.Send("RPUSH", fmt.Sprintf("%s:storage:%s", s.prefix, key), data)
		if err != nil {
			return
		}
		err = r.Redis.Send("ZADD", s.timerKey, ontime, key)
		if err != nil {
			return
		}
		_, err = r.Redis.Do("EXEC")
		if err != nil {
			return
		}
	})
	if e != nil {
		err = e
	}
	return
}

func (s *RedisStorage) Load(key string) (data [][]byte, err error) {
	var reply []interface{}
	e := s.redis.Do(func(i multiplexer.Instance) {
		r := i.(*broker.RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(broker.NetworkTimeout))

		storageKey := fmt.Sprintf("%s:storage:%s", s.prefix, key)
		reply, err = redis.Values(r.Redis.Do("LRANGE", storageKey, 0, -1))
		if err != nil {
			return
		}
		err = r.Redis.Send("MULTI")
		if err != nil {
			return
		}
		err = r.Redis.Send("ZREM", s.timerKey, key)
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
	if e != nil {
		err = e
	}
	if err != nil || len(reply) == 0 {
		return
	}

	data = make([][]byte, len(reply))
	for i, d := range reply {
		data[i], err = redis.Bytes(d, nil)
		if err != nil {
			return
		}
	}

	return
}

func (s *RedisStorage) Ontime(key string) (ontime int64, err error) {
	e := s.redis.Do(func(i multiplexer.Instance) {
		r := i.(*broker.RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(broker.NetworkTimeout))

		ontime, err = redis.Int64(r.Redis.Do("ZSCORE", s.timerKey, key))
	})
	if e != nil {
		err = e
	}
	if err == redis.ErrNil {
		ontime, err = 0, nil
	}
	return
}

func (s *RedisStorage) Next() (key string, err error) {
	var reply []interface{}
	e := s.redis.Do(func(i multiplexer.Instance) {
		r := i.(*broker.RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(broker.NetworkTimeout))

		reply, err = redis.Values(r.Redis.Do("ZRANGEBYSCORE", s.timerKey, "-inf", "+inf", "LIMIT", 0, 1, "WITHSCORES"))
	})
	if e != nil {
		err = e
	}
	if err != nil || len(reply) == 0 {
		return
	}

	key, err = redis.String(reply[0], nil)

	return
}

type Timer struct {
	updateType updateType
	storage    TimerStorage
}

func NewTimer(updateType updateType, storage TimerStorage) (*Timer, error) {
	switch updateType {
	case Always:
	case Never:
	default:
		return nil, fmt.Errorf("invalid update type: %s", updateType)
	}
	return &Timer{
		updateType: updateType,
		storage:    storage,
	}, nil
}

func (t *Timer) Push(ontime int64, key string, data []byte) error {
	return t.storage.Save(ontime, key, data)
}

func (t *Timer) Pop() (string, [][]byte, error) {
	key, err := t.storage.Next()
	if err != nil {
		return "", nil, err
	}
	data, err := t.storage.Load(key)
	if err != nil {
		return "", nil, err
	}
	return key, data, nil
}

func (t *Timer) NextWakeup() (time.Duration, error) {
	key, err := t.storage.Next()
	if err != nil {
		return -1, err
	}
	ontime, err := t.storage.Ontime(key)
	if err != nil {
		return -1, err
	}
	next := time.Unix(ontime, 0).Sub(time.Now())
	if next < 0 {
		next = 0
	}
	return next, nil
}
