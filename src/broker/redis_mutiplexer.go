package broker

import (
	"container/list"
	"github.com/googollee/godis"
	"launchpad.net/tomb"
	"model"
	"time"
)

type RedisMultiplexer struct {
	list    *list.List
	config  *model.Config
	get     chan *godis.Client
	back    chan *godis.Client
	timeout time.Duration
	tomb    tomb.Tomb
	max     int
}

func NewRedisMultiplexer(config *model.Config) *RedisMultiplexer {
	ret := &RedisMultiplexer{
		list:    list.New(),
		config:  config,
		get:     make(chan *godis.Client),
		back:    make(chan *godis.Client),
		timeout: time.Duration(config.Redis.HeartBeatInSecond) * time.Second,
		max:     int(config.Redis.MaxConnections),
	}
	if ret.max == 0 {
		config.Log.Crit("config Redis.MaxConnections should not 0!")
		panic("config Redis.MaxConnections should not 0!")
	}
	for i := 0; i < ret.max; i++ {
		ret.list.PushBack(ret.createRedis())
	}
	go looper(ret)
	return ret
}

func (m *RedisMultiplexer) Close() error {
	for i := 0; i < m.max; i++ {
		redis := <-m.get
		redis.Quit()
	}

	m.tomb.Kill(nil)
	err := m.tomb.Wait()
	if err != nil {
		return err
	}
	close(m.get)
	close(m.back)
	return nil
}

func (m *RedisMultiplexer) Quit() error {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return nil // no quit
}

func (m *RedisMultiplexer) Get(key string) (godis.Elem, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Get(key)
}

func (m *RedisMultiplexer) Set(key string, value interface{}) error {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Set(key, value)
}

func (m *RedisMultiplexer) Del(keys ...string) (int64, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Del(keys...)
}

func (m *RedisMultiplexer) Rpush(key string, value interface{}) (int64, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Rpush(key, value)
}

func (m *RedisMultiplexer) Lrange(key string, start, stop int) (*godis.Reply, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Lrange(key, start, stop)
}

func (m *RedisMultiplexer) Zadd(key string, score interface{}, member interface{}) (bool, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Zadd(key, score, member)
}

func (m *RedisMultiplexer) Zrem(key string, member interface{}) (bool, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Zrem(key, member)
}

func (m *RedisMultiplexer) Zcount(key string, min float64, max float64) (int64, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Zcount(key, min, max)
}

func (m *RedisMultiplexer) Zscore(key string, member interface{}) (float64, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Zscore(key, member)
}

func (m *RedisMultiplexer) Zrange(key string, start int, stop int) (*godis.Reply, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Zrange(key, start, stop)
}

func (m *RedisMultiplexer) Zrangebyscore(key string, min string, max string, args ...string) (*godis.Reply, error) {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return redis.Zrangebyscore(key, min, max, args...)
}

func (m *RedisMultiplexer) NewPipeClient() RedisPipe {
	redis := <-m.get
	defer func() { m.back <- redis }()
	return godis.NewPipeClientFromClient(redis)
}

func (m RedisMultiplexer) createRedis() *godis.Client {
	return godis.New(m.config.Redis.Netaddr, m.config.Redis.Db, m.config.Redis.Password)
}

func looper(m *RedisMultiplexer) {
	defer m.tomb.Done()

	last := time.Now()
	for {
		if elem := m.list.Front(); elem != nil {
			c := elem.Value.(*godis.Client)
			select {
			case m.get <- c:
				m.list.Remove(elem)
			case r := <-m.back:
				m.list.PushBack(r)
			case <-time.After(m.timeout):
			case <-m.tomb.Dying():
				return
			}
		} else {
			select {
			case r := <-m.back:
				m.list.PushBack(r)
			case <-time.After(m.timeout):
			case <-m.tomb.Dying():
				return
			}
		}
		if time.Now().Sub(last) >= m.timeout {
			last = time.Now()
			for i := m.list.Front(); i != nil; i = i.Next() {
				redis, ok := i.Value.(*godis.Client)
				if !ok {
					m.config.Log.Crit("value %s is not godis.Client", i.Value)
					i.Value = m.createRedis()
					continue
				}
				reply, err := redis.Ping()
				if err != nil || reply.String() != "PONG" {
					redis.Quit()
					i.Value = m.createRedis()
				}
			}
		}
	}
}
