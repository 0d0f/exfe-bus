package broker

import (
	"github.com/googollee/godis"
)

type RedisClient interface {
	Quit() error

	Get(key string) (godis.Elem, error)
	Set(key string, value interface{}) error
	Del(keys ...string) (int64, error)

	Rpush(key string, value interface{}) (int64, error)
	Lrange(key string, start, stop int) (*godis.Reply, error)

	Zadd(key string, score interface{}, member interface{}) (bool, error)
	Zrem(key string, member interface{}) (bool, error)
	Zcount(key string, min float64, max float64) (int64, error)
	Zscore(key string, member interface{}) (float64, error)
	Zrange(key string, start int, stop int) (*godis.Reply, error)
	Zrangebyscore(key string, min string, max string, args ...string) (*godis.Reply, error)
}

type RedisPipe interface {
	RedisClient
	Multi() error
	Watch(keys ...string) error
	Exec() []*godis.Reply
}

type Redis interface {
	RedisClient
	NewPipeClient() RedisPipe
}

type RedisImp struct {
	redis *godis.Client
}

func NewRedisImp() *RedisImp {
	return &RedisImp{
		redis: godis.New("", 0, ""),
	}
}

func (r *RedisImp) NewPipeClient() RedisPipe {
	return godis.NewPipeClientFromClient(r.redis)
}

func (r *RedisImp) Quit() error {
	return r.redis.Quit()
}

func (r *RedisImp) Del(keys ...string) (int64, error) {
	return r.redis.Del(keys...)
}

func (r *RedisImp) Rpush(key string, value interface{}) (int64, error) {
	return r.redis.Rpush(key, value)
}

func (r *RedisImp) Lrange(key string, start, stop int) (*godis.Reply, error) {
	return r.redis.Lrange(key, start, stop)
}

func (r *RedisImp) Zadd(key string, score interface{}, member interface{}) (bool, error) {
	return r.redis.Zadd(key, score, member)
}

func (r *RedisImp) Zrem(key string, member interface{}) (bool, error) {
	return r.redis.Zrem(key, member)
}

func (r *RedisImp) Zscore(key string, member interface{}) (float64, error) {
	return r.redis.Zscore(key, member)
}

func (r *RedisImp) Zrange(key string, start int, stop int) (*godis.Reply, error) {
	return r.redis.Zrange(key, start, stop)
}

func (r *RedisImp) Zrangebyscore(key string, min string, max string, args ...string) (*godis.Reply, error) {
	return r.redis.Zrangebyscore(key, min, max, args...)
}
