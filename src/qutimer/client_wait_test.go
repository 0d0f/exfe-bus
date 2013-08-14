package qutimer

import (
	"github.com/googollee/go-assert"
	"testing"
	"time"
)

func TestClientWait(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()
	defer clearQueue()

	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)

	begin := time.Now()
	client.Wait(time.Second)
	end := time.Now()
	assert.Equal(t, end.Sub(begin) > time.Second, true)
}

func TestClientWaitPush(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()
	defer clearQueue()

	client, err := New(redisPool, prefix, 2)
	assert.MustEqual(t, err, nil)

	go func() {
		time.Sleep(time.Second / 10)
		conn.Do("PUBLISH", pushChannel, "insert")
	}()

	begin := time.Now()
	client.Wait(time.Second)
	end := time.Now()
	assert.Equal(t, end.Sub(begin) < time.Second/5, true)
}
