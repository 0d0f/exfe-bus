package delayrepo

import (
	"broker"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

func TestHead(t *testing.T) {
	var q Repo
	redis := broker.NewRedisImp("", 0, "")

	log, err := logger.New(logger.Stderr, "test")
	if err != nil {
		panic(err)
	}
	q = NewHead("hdtest", 2, redis)

	tomb := ServRepository(log.SubPrefix("serv"), q, func(key string, data [][]byte) {
		assert.Equal(t, key, "test1")
		assert.Equal(t, fmt.Sprintf("%+v", data), "[[48] [49] [50] [51] [52] [53] [54] [55] [56] [57]]")
	})

	next, err := q.NextWakeup()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, 2*time.Second)

	{
		for i := 0; i < 10; i++ {
			q.Push("test1", []byte(fmt.Sprintf("%d", i)))
		}
		next, err = q.NextWakeup()
		assert.Equal(t, err, nil)
		assert.Equal(t, next, 2*time.Second)
		time.Sleep(next * 3 / 2)
	}

	{
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
		assert.Equal(t, next, time.Second)

		time.Sleep(next * 3 / 2)
	}

	tomb.Kill(nil)
	tomb.Wait()
}
