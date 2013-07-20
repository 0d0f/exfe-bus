package broker

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type UpdateType string

const (
	Always UpdateType = "always"
	Once              = "once"
)

type QueueRedisStorage struct {
	redis    *redis.Pool
	prefix   string
	timerKey string
}

func NewQueueRedisStorage(prefix string, redis *redis.Pool) *QueueRedisStorage {
	return &QueueRedisStorage{
		redis:    redis,
		prefix:   prefix,
		timerKey: fmt.Sprintf("%s:timer", prefix),
	}
}

func (s *QueueRedisStorage) Save(updateType UpdateType, ontime int64, key string, data []byte) error {
	conn := s.redis.Get()
	defer conn.Close()

	if err := conn.Send("MULTI"); err != nil {
		return err
	}
	if err := conn.Send("RPUSH", fmt.Sprintf("%s:storage:%s", s.prefix, key), data); err != nil {
		return err
	}
	if err := conn.Send("ZADD", s.timerKey, ontime, key); err != nil {
		return err
	}
	if _, err := conn.Do("EXEC"); err != nil {
		return err
	}
	return nil
}

func (s *QueueRedisStorage) Load(key string) ([][]byte, error) {
	conn := s.redis.Get()
	defer conn.Close()

	storageKey := fmt.Sprintf("%s:storage:%s", s.prefix, key)
	reply, err := redis.Values(conn.Do("LRANGE", storageKey, 0, -1))
	if err != nil {
		return nil, err
	}
	if len(reply) == 0 {
		return nil, nil
	}
	if err = conn.Send("MULTI"); err != nil {
		return nil, err
	}
	if err = conn.Send("ZREM", s.timerKey, key); err != nil {
		return nil, err
	}
	if err = conn.Send("DEL", storageKey); err != nil {
		return nil, err
	}
	if _, err = conn.Do("EXEC"); err != nil {
		return nil, err
	}

	data := make([][]byte, len(reply))
	for i, d := range reply {
		data[i], err = redis.Bytes(d, nil)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s *QueueRedisStorage) Ontime(key string) (int64, error) {
	conn := s.redis.Get()
	defer conn.Close()

	ontime, err := redis.Int64(conn.Do("ZSCORE", s.timerKey, key))
	if err == redis.ErrNil {
		ontime, err = 0, nil
	}
	return ontime, err
}

func (s *QueueRedisStorage) Next() (string, error) {
	conn := s.redis.Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("ZRANGEBYSCORE", s.timerKey, "-inf", "+inf", "LIMIT", 0, 1, "WITHSCORES"))
	if err != nil {
		return "", err
	}
	if len(reply) == 0 {
		return "", nil
	}

	key, err := redis.String(reply[0], nil)

	return key, err
}
