package broker

import (
	"github.com/googollee/go-logger"
	"github.com/googollee/godis"
	"model"
	"testing"
	"time"
)

func TestRedisMultiplexer(t *testing.T) {
	log, err := logger.New(logger.Stderr, "test")
	if err != nil {
		t.Fatalf("new logger failed: %s", err)
	}
	config := new(model.Config)
	config.Log = log
	m := NewRedisMultiplexer(config)

	{
		redis := <-m.get
		m.back <- redis
	}

	{
		redises := make([]*godis.Client, m.max)
		for i := range redises {
			redises[i] = <-m.get
		}
		select {
		case <-m.get:
			t.Errorf("should get nothing")
		case <-time.After(time.Second):
		}
		for i := range redises {
			m.back <- redises[i]
		}
	}

	{
		m.timeout = time.Second
		time.Sleep(10 * time.Second)
	}

	m.Close()
}
