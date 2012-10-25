package delayrepo

import (
	"github.com/googollee/godis"
)

type RedisClientBroker interface {
	Quit() error

	Del(keys ...string) (int64, error)

	Rpush(key string, value interface{}) (int64, error)
	Lrange(key string, start, stop int) (*godis.Reply, error)

	Zadd(key string, score interface{}, member interface{}) (bool, error)
	Zrem(key string, member interface{}) (bool, error)
	Zscore(key string, member interface{}) (float64, error)
	Zrange(key string, start int, stop int) (*godis.Reply, error)
	Zrangebyscore(key string, min string, max string, args ...string) (*godis.Reply, error)
}

type RedisPipeBroker interface {
	RedisClientBroker
	Multi() error
	Watch(keys ...string) error
	Exec() []*godis.Reply
}

type RedisBroker interface {
	RedisClientBroker
	NewPipeClient() RedisPipeBroker
}
