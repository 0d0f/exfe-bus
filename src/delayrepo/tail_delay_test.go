package delayrepo

import (
	"broker"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

func TestTailDelayQueue(t *testing.T) {
	var q Repository
	redis := broker.NewRedisImp()

	log, err := logger.New(logger.Stderr, "test")
	if err != nil {
		panic(err)
	}
	q = NewTail("tdt", 2, redis)
	expectData := ""

	tomb := ServRepository(log.SubPrefix("serv"), q, func(key string, data [][]byte) {
		assert.Equal(t, key, "test1")
		assert.Equal(t, fmt.Sprintf("%+v", data), expectData)
	})

	next, err := q.NextWakeup()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, 2*time.Second)

	{
		expectData = "[[48] [49] [50] [51] [52] [53] [54] [55] [56] [57]]"
		for i := 0; i < 10; i++ {
			q.Push("test1", []byte(fmt.Sprintf("%d", i)))
		}
		next, err = q.NextWakeup()
		assert.Equal(t, err, nil)
		assert.Equal(t, next, 2*time.Second)
		time.Sleep(next)
	}

	{
		expectData = "[[48] [49] [50] [51] [52] [53] [54] [55] [56] [57]]"
		for i := 0; i < 5; i++ {
			q.Push("test1", []byte(fmt.Sprintf("%d", i)))
		}
		next, err = q.NextWakeup()
		assert.Equal(t, err, nil)
		assert.Equal(t, next, 2*time.Second)

		time.Sleep(next / 2)

		for i := 5; i < 10; i++ {
			q.Push("test1", []byte(fmt.Sprintf("%d", i)))
		}
		next, err = q.NextWakeup()
		assert.Equal(t, err, nil)
		assert.Equal(t, next, 2*time.Second)
	}

	tomb.Kill(nil)
	tomb.Wait()
}
