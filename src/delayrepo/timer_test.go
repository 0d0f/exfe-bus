package delayrepo

import (
	"broker"
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"model"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	config := new(model.Config)
	config.Redis.MaxConnections = 1
	config.Redis.Netaddr = "127.0.0.1:6379"
	redis, err := broker.NewRedisPool(config)
	if err != nil {
		panic(err)
	}

	strategy := NewTimer("delay:test", redis)
	ontime := time.Now().Add(time.Second).Unix()
	err = strategy.Push(ontime, "123", []byte("a"))
	if err != nil {
		panic(err)
	}
	err = strategy.Push(ontime, "123", []byte("b"))
	if err != nil {
		panic(err)
	}
	err = strategy.Push(ontime, "123", []byte("c"))
	if err != nil {
		panic(err)
	}

	wait, err := strategy.NextWakeup()
	if err != nil {
		panic(err)
	}
	fmt.Println("waiting:", wait)

	time.Sleep(wait)

	key, data, err := strategy.Pop()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, key, "123")
	assert.Equal(t, fmt.Sprintf("%v", data), "[[97] [98] [99]]")
}
