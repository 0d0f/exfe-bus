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
		t.Fatal(err)
	}

	strategy, err := NewTimer(Always, "delay:test", redis)
	if err != nil {
		t.Fatal(err)
	}
	ontime := time.Now().Add(time.Second).Unix()
	err = strategy.Push(ontime, "123", []byte("a"))
	if err != nil {
		t.Fatal(err)
	}
	err = strategy.Push(ontime, "123", []byte("b"))
	if err != nil {
		t.Fatal(err)
	}
	err = strategy.Push(ontime, "123", []byte("c"))
	if err != nil {
		t.Fatal(err)
	}

	wait, err := strategy.NextWakeup()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("waiting:", wait)

	time.Sleep(wait)

	key, data, err := strategy.Pop()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, key, "123")
	assert.Equal(t, fmt.Sprintf("%v", data), "[[97] [98] [99]]")
}

func TestTimerUpdate(t *testing.T) {
	config := new(model.Config)
	config.Redis.MaxConnections = 1
	config.Redis.Netaddr = "127.0.0.1:6379"
	redis, err := broker.NewRedisPool(config)
	if err != nil {
		t.Fatal(err)
	}

	strategy, err := NewTimer(Always, "delay:test", redis)
	if err != nil {
		t.Fatal(err)
	}
	ontime := time.Now().Add(time.Second * 10).Unix()
	err = strategy.Push(ontime, "123", []byte("a"))
	if err != nil {
		t.Fatal(err)
	}
	wait, err := strategy.NextWakeup()
	if err != nil {
		t.Fatal(err)
	}
	if wait < time.Second {
		t.Fatalf("wait too short: %s", wait)
	}

	ontime = time.Now().Unix()
	err = strategy.Push(ontime, "123", []byte("b"))
	if err != nil {
		t.Fatal(err)
	}
	wait, err = strategy.NextWakeup()
	if err != nil {
		t.Fatal(err)
	}
	if wait > time.Second {
		t.Fatalf("wait too long: %s", wait)
	}

	time.Sleep(wait)

	key, data, err := strategy.Pop()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, key, "123")
	assert.Equal(t, fmt.Sprintf("%v", data), "[[97] [98]]")
}
