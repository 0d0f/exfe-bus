package broker

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-logger"
	"github.com/googollee/go-multiplexer"
	"github.com/googollee/godis"
	"model"
	"net"
	"time"
)

type RedisInstance_ struct {
	Conn  net.Conn
	Redis redis.Conn
	log   *logger.SubLogger
}

func (i *RedisInstance_) Ping() error {
	reply, err := redis.String(i.Redis.Do("PING"))
	if reply != "PONG" {
		err = fmt.Errorf("redis not pong.")
	}
	return err
}

func (i *RedisInstance_) Close() error {
	return i.Redis.Close()
}

func (i *RedisInstance_) Error(err error) {
	i.log.Err("%s", err)
}

type RedisPool struct {
	homo   *multiplexer.Homo
	config *model.Config
}

func NewRedisPool(config *model.Config) (*RedisPool, error) {
	if config.Redis.MaxConnections == 0 {
		return nil, fmt.Errorf("config Redis.MaxConnections should not 0!")
	}
	return &RedisPool{
		homo: multiplexer.NewHomo(func() (multiplexer.Instance, error) {
			conn, err := net.DialTimeout("tcp", config.Redis.Netaddr, NetworkTimeout)
			if err != nil {
				return nil, err
			}
			return &RedisInstance_{
				Conn:  conn,
				Redis: redis.NewConn(conn, 0, 0),
				log:   config.Log.SubPrefix("redis"),
			}, nil
		}, config.Redis.MaxConnections, -1, time.Duration(config.Redis.HeartBeatInSecond)*time.Second),
		config: config,
	}, nil
}

func (r *RedisPool) Do(f func(multiplexer.Instance)) error {
	return r.homo.Do(f)
}

func (r *RedisPool) Close() error {
	return r.homo.Close()
}

type UpdateType string

const (
	Always UpdateType = "always"
	Once              = "once"
)

type QueueRedisStorage struct {
	redis    *RedisPool
	prefix   string
	timerKey string
}

func NewQueueRedisStorage(prefix string, redis *RedisPool) *QueueRedisStorage {
	return &QueueRedisStorage{
		redis:    redis,
		prefix:   prefix,
		timerKey: fmt.Sprintf("%s:timer", prefix),
	}
}

func (s *QueueRedisStorage) Save(updateType UpdateType, ontime int64, key string, data []byte) (err error) {
	e := s.redis.Do(func(i multiplexer.Instance) {
		r := i.(*RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(NetworkTimeout))

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

func (s *QueueRedisStorage) Load(key string) (data [][]byte, err error) {
	var reply []interface{}
	e := s.redis.Do(func(i multiplexer.Instance) {
		r := i.(*RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(NetworkTimeout))

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

func (s *QueueRedisStorage) Ontime(key string) (ontime int64, err error) {
	e := s.redis.Do(func(i multiplexer.Instance) {
		r := i.(*RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(NetworkTimeout))

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

func (s *QueueRedisStorage) Next() (key string, err error) {
	var reply []interface{}
	e := s.redis.Do(func(i multiplexer.Instance) {
		r := i.(*RedisInstance_)
		r.Conn.SetDeadline(time.Now().Add(NetworkTimeout))

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

// old

type RedisInstance struct {
	redis *godis.Client
	log   *logger.SubLogger
}

func (i *RedisInstance) Ping() error {
	_, err := i.redis.Ping()
	return err
}

func (i *RedisInstance) Close() error {
	return i.redis.Quit()
}

func (i *RedisInstance) Error(err error) {
	i.log.Err("%s", err)
}

type RedisMultiplexer struct {
	homo   *multiplexer.Homo
	config *model.Config
}

func NewRedisMultiplexer(config *model.Config) *RedisMultiplexer {
	if config.Redis.MaxConnections == 0 {
		config.Log.Crit("config Redis.MaxConnections should not 0!")
		panic("config Redis.MaxConnections should not 0!")
	}
	return &RedisMultiplexer{
		homo: multiplexer.NewHomo(func() (multiplexer.Instance, error) {
			return &RedisInstance{
				redis: godis.New("tcp:"+config.Redis.Netaddr, config.Redis.Db, config.Redis.Password),
				log:   config.Log.SubPrefix("redis"),
			}, nil
		}, config.Redis.MaxConnections, -1, time.Duration(config.Redis.HeartBeatInSecond)*time.Second),
		config: config,
	}
}

func (m *RedisMultiplexer) Close() error {
	return m.homo.Close()
}

func (m *RedisMultiplexer) Quit() error {
	return nil // no quit
}

func (m *RedisMultiplexer) Get(key string) (elem godis.Elem, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		elem, err = redis.Get(key)
	})
	return
}

func (m *RedisMultiplexer) Set(key string, value interface{}) (err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		err = redis.Set(key, value)
	})
	return
}

func (m *RedisMultiplexer) Incrby(key string, increment int64) (ret int64, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Incrby(key, increment)
	})
	return
}

func (m *RedisMultiplexer) Del(keys ...string) (ret int64, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Del(keys...)
	})
	return
}

func (m *RedisMultiplexer) Rpush(key string, value interface{}) (ret int64, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Rpush(key, value)
	})
	return
}

func (m *RedisMultiplexer) Lrange(key string, start, stop int) (ret *godis.Reply, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Lrange(key, start, stop)
	})
	return
}

func (m *RedisMultiplexer) Zadd(key string, score interface{}, member interface{}) (ret bool, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Zadd(key, score, member)
	})
	return
}

func (m *RedisMultiplexer) Zrem(key string, member interface{}) (ret bool, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Zrem(key, member)
	})
	return
}

func (m *RedisMultiplexer) Zcount(key string, min float64, max float64) (ret int64, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Zcount(key, min, max)
	})
	return
}

func (m *RedisMultiplexer) Zscore(key string, member interface{}) (ret float64, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Zscore(key, member)
	})
	return
}

func (m *RedisMultiplexer) Zrange(key string, start int, stop int) (ret *godis.Reply, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Zrange(key, start, stop)
	})
	return
}

func (m *RedisMultiplexer) Zrangebyscore(key string, min string, max string, args ...string) (ret *godis.Reply, err error) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret, err = redis.Zrangebyscore(key, min, max, args...)
	})
	return
}

func (m *RedisMultiplexer) NewPipeClient() (ret RedisPipe) {
	m.homo.Do(func(i multiplexer.Instance) {
		redis := i.(*RedisInstance).redis
		ret = godis.NewPipeClientFromClient(redis)
	})
	return
}
