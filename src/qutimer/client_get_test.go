package qutimer

import (
	"bytes"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-assert"
	"testing"
	"time"
)

func pushForGetTest(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	number := 10
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	overwrite := "always"
	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)
	for i := 0; i < number; i++ {
		task := bytes.NewBufferString(fmt.Sprintf("%d", i))
		ontime := int64(12340 + i)
		err = client.Push(key, task, data, ontime, overwrite)
		assert.Equal(t, err, nil)
	}
}

func TestClientGet1(t *testing.T) {
	pushForGetTest(t)
	conn := redisPool.Get()
	defer conn.Close()
	defer clearQueue()

	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)

	task, err := client.Get(key, 1)
	assert.Equal(t, err, nil)
	assert.MustEqual(t, len(task.Data()), 1)
	data, err := redis.String(task.Data()[0], nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, data, "0")

	i64, err := redis.Int64(conn.Do("ZSCORE", timer, key))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, time.Now().Unix()+2)
	i, err := redis.Int(conn.Do("LLEN", queue))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, 10)
	ok, err := redis.Bool(conn.Do("EXISTS", queueOverwrite))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, true)
	ok, err = redis.Bool(conn.Do("EXISTS", queueStart))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueLocker))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, true)
	i, err = redis.Int(conn.Do("HLEN", queueData))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, 2)
}

func TestClientGet2(t *testing.T) {
	pushForGetTest(t)
	conn := redisPool.Get()
	defer conn.Close()
	defer clearQueue()

	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)

	task, err := client.Get(key, 2)
	assert.Equal(t, err, nil)
	assert.MustEqual(t, len(task.Data()), 2)
	data, err := redis.String(task.Data()[0], nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, data, "0")
	data, err = redis.String(task.Data()[1], nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, data, "1")

	i64, err := redis.Int64(conn.Do("ZSCORE", timer, key))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, time.Now().Unix()+2)
	i, err := redis.Int(conn.Do("LLEN", queue))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, 10)
	ok, err := redis.Bool(conn.Do("EXISTS", queueOverwrite))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, true)
	ok, err = redis.Bool(conn.Do("EXISTS", queueStart))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueLocker))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, true)
	i, err = redis.Int(conn.Do("HLEN", queueData))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, 2)
}

func TestClientGetAll(t *testing.T) {
	pushForGetTest(t)
	conn := redisPool.Get()
	defer conn.Close()
	defer clearQueue()

	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)

	task, err := client.Get(key, -1)
	assert.Equal(t, err, nil)
	assert.MustEqual(t, len(task.Data()), 10)
	for i, ta := range task.Data() {
		data, err := redis.String(ta, nil)
		assert.Equal(t, err, nil)
		assert.Equal(t, data, fmt.Sprintf("%d", i))
	}

	i64, err := redis.Int64(conn.Do("ZSCORE", timer, key))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, time.Now().Unix()+2)
	i, err := redis.Int(conn.Do("LLEN", queue))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, 10)
	ok, err := redis.Bool(conn.Do("EXISTS", queueOverwrite))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, true)
	ok, err = redis.Bool(conn.Do("EXISTS", queueStart))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueLocker))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, true)
	i, err = redis.Int(conn.Do("HLEN", queueData))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, 2)
}

func TestClientGetLocker(t *testing.T) {
	pushForGetTest(t)
	conn := redisPool.Get()
	defer conn.Close()
	defer clearQueue()

	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)
	conn.Do("SET", queueLocker, "1")

	_, err = client.Get(key, -1)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err.Error(), "locked")

	i64, err := redis.Int64(conn.Do("ZSCORE", timer, key))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, int64(12349))
	i, err := redis.Int(conn.Do("LLEN", queue))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, 10)
	ok, err := redis.Bool(conn.Do("EXISTS", queueOverwrite))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueStart))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueLocker))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, true)
	i, err = redis.Int(conn.Do("HLEN", queueData))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, 2)
}
