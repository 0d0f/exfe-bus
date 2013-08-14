package qutimer

import (
	"bytes"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-assert"
	"testing"
)

func TestClientPush1(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	number := 1
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	overwrite := "always"
	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)

	cont := make(chan int)
	go func() {
		conn, err := redisPool.Dial()
		assert.MustEqual(t, err, nil)
		defer conn.Close()

		ps := redis.PubSubConn{conn}
		ps.Subscribe(pushChannel)
		cont <- 1
		defer func() {
			cont <- 1
		}()
		for i := 0; i < number; {
			switch n := ps.Receive().(type) {
			case redis.Message:
				assert.Equal(t, string(n.Data), "insert")
				i++
			case redis.Subscription:
			default:
				t.Fatal("should received")
			}
		}
	}()
	<-cont

	for i := 0; i < number; i++ {
		task := bytes.NewBufferString(fmt.Sprintf("%d", i))
		ontime := int64(12340 + i)
		err = client.Push(key, task, data, ontime, overwrite)
		assert.Equal(t, err, nil)
	}
	<-cont

	i64, err := redis.Int64(conn.Do("ZSCORE", timer, key))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, int64(12340))
	i, err := redis.Int(conn.Do("LLEN", queue))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, number)
	i, err = redis.Int(conn.Do("HLEN", queueData))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, len(data))
	ok, err := redis.Bool(conn.Do("EXISTS", queueOverwrite))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueLocker))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)

	conn.Do("DEL", queue)
	conn.Do("DEL", queueOverwrite)
	conn.Do("DEL", queueLocker)
	conn.Do("DEL", queueData)
	conn.Do("DEL", timer)
}

func TestClientPush10Always(t *testing.T) {
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

	cont := make(chan int)
	go func() {
		conn, err := redisPool.Dial()
		assert.MustEqual(t, err, nil)
		defer conn.Close()

		ps := redis.PubSubConn{conn}
		ps.Subscribe(pushChannel)
		cont <- 1
		defer func() {
			cont <- 1
		}()
		for i := 0; i < number; {
			switch n := ps.Receive().(type) {
			case redis.Message:
				assert.Equal(t, string(n.Data), "insert")
				i++
			case redis.Subscription:
			default:
				t.Fatal("should received")
			}
		}
	}()
	<-cont

	for i := 0; i < number; i++ {
		ontime := int64(12340 + i)
		task := bytes.NewBufferString(fmt.Sprintf("%d", i))
		err = client.Push(key, task, data, ontime, overwrite)
		assert.Equal(t, err, nil)
	}
	<-cont

	i64, err := redis.Int64(conn.Do("ZSCORE", timer, key))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, int64(12349))
	i, err := redis.Int(conn.Do("LLEN", queue))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, number)
	i, err = redis.Int(conn.Do("HLEN", queueData))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, len(data))
	ok, err := redis.Bool(conn.Do("EXISTS", queueOverwrite))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueLocker))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)

	conn.Do("DEL", queue)
	conn.Do("DEL", queueOverwrite)
	conn.Do("DEL", queueLocker)
	conn.Do("DEL", queueData)
	conn.Do("DEL", timer)
}

func TestClientPush10Once(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	number := 10
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	overwrite := "once"
	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)

	cont := make(chan int)
	go func() {
		conn, err := redisPool.Dial()
		assert.MustEqual(t, err, nil)
		defer conn.Close()

		ps := redis.PubSubConn{conn}
		ps.Subscribe(pushChannel)
		cont <- 1
		defer func() {
			cont <- 1
		}()
		for i := 0; i < number; {
			switch n := ps.Receive().(type) {
			case redis.Message:
				assert.Equal(t, string(n.Data), "insert")
				i++
			case redis.Subscription:
			default:
				t.Fatal("should received")
			}
		}
	}()
	<-cont

	for i := 0; i < number; i++ {
		ontime := int64(12340 + i)
		task := bytes.NewBufferString(fmt.Sprintf("%d", i))
		err = client.Push(key, task, data, ontime, overwrite)
		assert.Equal(t, err, nil)
	}
	<-cont

	i64, err := redis.Int64(conn.Do("ZSCORE", timer, key))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, int64(12340))
	i, err := redis.Int(conn.Do("LLEN", queue))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, number)
	i, err = redis.Int(conn.Do("HLEN", queueData))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, len(data))
	ok, err := redis.Bool(conn.Do("EXISTS", queueOverwrite))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueLocker))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)

	conn.Do("DEL", queue)
	conn.Do("DEL", queueOverwrite)
	conn.Do("DEL", queueLocker)
	conn.Do("DEL", queueData)
	conn.Do("DEL", timer)
}

func TestClientPush10OnceOverwrite(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	number := 10
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	overwrite := "once"
	conn.Do("SET", queueOverwrite, "1")
	conn.Do("ZADD", timer, 56789, key)
	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)

	cont := make(chan int)
	go func() {
		conn, err := redisPool.Dial()
		assert.MustEqual(t, err, nil)
		defer conn.Close()

		ps := redis.PubSubConn{conn}
		ps.Subscribe(pushChannel)
		cont <- 1
		defer func() {
			cont <- 1
		}()
		for i := 0; i < number; {
			switch n := ps.Receive().(type) {
			case redis.Message:
				assert.Equal(t, string(n.Data), "insert")
				i++
			case redis.Subscription:
			default:
				t.Fatal("should received")
			}
		}
	}()
	<-cont

	for i := 0; i < number; i++ {
		ontime := int64(12340 + i)
		task := bytes.NewBufferString(fmt.Sprintf("%d", i))
		err = client.Push(key, task, data, ontime, overwrite)
		assert.Equal(t, err, nil)
	}
	<-cont

	i64, err := redis.Int64(conn.Do("ZSCORE", timer, key))
	assert.Equal(t, err, nil)
	assert.Equal(t, i64, int64(12340))
	i, err := redis.Int(conn.Do("LLEN", queue))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, number)
	i, err = redis.Int(conn.Do("HLEN", queueData))
	assert.Equal(t, err, nil)
	assert.Equal(t, i, len(data))
	ok, err := redis.Bool(conn.Do("EXISTS", queueOverwrite))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)
	ok, err = redis.Bool(conn.Do("EXISTS", queueLocker))
	assert.Equal(t, err, nil)
	assert.Equal(t, ok, false)

	conn.Do("DEL", queue)
	conn.Do("DEL", queueOverwrite)
	conn.Do("DEL", queueLocker)
	conn.Do("DEL", queueData)
	conn.Do("DEL", timer)
}
