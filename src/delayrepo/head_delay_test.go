package delayrepo

import (
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

func TestHead(t *testing.T) {
	var q Repository
	redis := NewRedis()

	log, err := logger.New(logger.Stderr, "test", logger.LstdFlags)
	if err != nil {
		panic(err)
	}
	q = NewHead("hdt", 2, redis)

	tomb := ServRepository(log.SubPrefix("serv"), q, func(key string, data [][]byte) {
		assert.Equal(t, key, "test1")
		assert.Equal(t, fmt.Sprintf("%+v", data), "[[48] [49] [50] [51] [52] [53] [54] [55] [56] [57]]")
	})

	next, err := q.NextWakeup()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, 2*time.Second)

	for i := 0; i < 10; i++ {
		q.Push("test1", []byte(fmt.Sprintf("%d", i)))
	}
	next, err = q.NextWakeup()
	assert.Equal(t, err, nil)
	time.Sleep(next * 3 / 2)

	tomb.Kill(nil)
	tomb.Wait()
}
