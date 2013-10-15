package timer

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-assert"
	"math/rand"
	"testing"
	"time"
)

var redisPool = &redis.Pool{
	MaxIdle:     3,
	IdleTimeout: 30 * time.Minute,
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			return nil, err
		}
		return c, nil
	},
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	},
}

func clean(prefix string) {
	conn := redisPool.Get()
	defer conn.Close()

	conn.Do("DEL", notifyName(prefix), sortedSetname(prefix), dataName(prefix))
}

func TestTimerServer(t *testing.T) {
	prefix := fmt.Sprintf("timer:test:%d.%d", time.Now().Unix(), rand.Intn(10000))
	defer clean(prefix)

	s := NewServer(redisPool, prefix, 2)
	var err error
	go func() { err = s.Serve() }()
	time.Sleep(time.Second)
	assert.MustEqual(t, err, nil, "launch timer failed: %s", err)
	s.Close()
}

func TestTimerClientSend(t *testing.T) {
	prefix := fmt.Sprintf("timer:test:%d.%d", time.Now().Unix(), rand.Intn(10000))
	defer clean(prefix)

	s := NewServer(redisPool, prefix, 2)
	go s.Serve()
	defer s.Close()

	c := NewClient(redisPool, prefix)
	err := c.Send(time.Now().Add(time.Second).Unix(), "test", "data")
	assert.Equal(t, err, nil)
}

func TestTimerClient(t *testing.T) {
	prefix := fmt.Sprintf("timer:test:%d.%d", time.Now().Unix(), rand.Intn(10000))
	defer clean(prefix)

	s := NewServer(redisPool, prefix, 2)
	defer s.Close()

	go s.Serve()

	p := make(chan int)
	r := make(chan int)

	c := NewClient(redisPool, prefix)
	go func() {
		ch := make(chan Event)
		go c.PListen("timer:*", ch)
		p <- 1
		<-r
		event := <-ch
		assert.Equal(t, event.Name, "timer:test")
		assert.Equal(t, string(event.Data), "data")
		p <- 1
	}()
	<-p
	r <- 1
	err := c.Send(time.Now().Add(time.Second).Unix(), "timer:test", "data")
	assert.Equal(t, err, nil)
	select {
	case <-p:
		t.Fail()
	default:
	}
	time.Sleep(time.Second * 2)
	select {
	case <-p:
	default:
		t.Fail()
	}
	c.Close()
}
