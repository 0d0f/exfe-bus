package delayrepo

import (
	"github.com/googollee/godis"
)

type Redis struct {
	redis *godis.Client
}

func (r *Redis) NewPipeClient() RedisPipeBroker {
	return godis.NewPipeClientFromClient(r.redis)
}

func (r *Redis) Quit() error {
	return r.redis.Quit()
}

func (r *Redis) Del(keys ...string) (int64, error) {
	return r.redis.Del(keys...)
}

func (r *Redis) Rpush(key string, value interface{}) (int64, error) {
	return r.redis.Rpush(key, value)
}

func (r *Redis) Lrange(key string, start, stop int) (*godis.Reply, error) {
	return r.redis.Lrange(key, start, stop)
}

func (r *Redis) Zadd(key string, score interface{}, member interface{}) (bool, error) {
	return r.redis.Zadd(key, score, member)
}

func (r *Redis) Zrem(key string, member interface{}) (bool, error) {
	return r.redis.Zrem(key, member)
}

func (r *Redis) Zscore(key string, member interface{}) (float64, error) {
	return r.redis.Zscore(key, member)
}

func (r *Redis) Zrange(key string, start int, stop int) (*godis.Reply, error) {
	return r.redis.Zrange(key, start, stop)
}

func (r *Redis) Zrangebyscore(key string, min string, max string, args ...string) (*godis.Reply, error) {
	return r.redis.Zrangebyscore(key, min, max, args...)
}

func NewRedis() *Redis {
	return &Redis{
		redis: godis.New("", 0, ""),
	}
}
