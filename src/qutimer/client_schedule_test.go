package qutimer

import (
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-assert"
	"testing"
	"time"
)

func TestClientScheduleEmpty(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()
	defer clearQueue()

	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)

	_, next, err := client.Schedule()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, 2*time.Second)
}

func TestClientScheduleNext(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()
	defer clearQueue()

	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)
	now := time.Now().Unix()

	conn.Do("ZADD", timer, now+10, "key1")

	_, next, err := client.Schedule()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, 10*time.Second)

	conn.Do("ZADD", timer, now+2, "key")

	_, next, err = client.Schedule()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, 2*time.Second)

	time.Sleep(next)

	key, next, err := client.Schedule()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, time.Duration(0))
	assert.Equal(t, key, "key")
}

func TestClientScheduleTimeout(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()
	defer clearQueue()

	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)
	now := time.Now().Unix()

	conn.Do("RPUSH", queue, "1")
	conn.Do("HSET", queueData, "key", "value")
	conn.Do("ZADD", timer, now, "key")

	ok, err := redis.Bool(conn.Do("EXISTS", queueStart))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)

	_, next, err := client.Schedule()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, time.Duration(0))

	i64, err := redis.Int64(conn.Do("GET", queueStart))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, now)

	time.Sleep(2 * time.Second)

	_, next, err = client.Schedule()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, time.Duration(0))

	i64, err = redis.Int64(conn.Do("GET", queueStart))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, now)

	time.Sleep(3 * time.Second)

	_, next, err = client.Schedule()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, 2*time.Second)

	ok, err = redis.Bool(conn.Do("EXISTS", queue))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueStart))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueData))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
}
