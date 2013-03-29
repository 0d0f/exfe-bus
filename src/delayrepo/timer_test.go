package delayrepo

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchrcom/testify/assert"
	"net"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	conn, err := net.DialTimeout("tcp", ":6379", time.Second)
	if err != nil {
		panic(err)
	}
	r := redis.NewConn(conn, 0, 0)

	queue := NewTimer("delay:queue", r)
	ontime := time.Now().Add(time.Second).Unix()
	err = queue.Push(ontime, "123", 1)
	if err != nil {
		panic(err)
	}
	err = queue.Push(ontime, "123", 2)
	if err != nil {
		panic(err)
	}
	err = queue.Push(ontime, "123", 3)
	if err != nil {
		panic(err)
	}

	wait, err := queue.NextWakeup()
	if err != nil {
		panic(err)
	}
	fmt.Println("waiting:", wait)

	time.Sleep(wait)

	key, data, err := queue.Pop()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, key, "123")
	assert.Equal(t, fmt.Sprintf("%v", data), "[[49] [50] [51]]")
}
