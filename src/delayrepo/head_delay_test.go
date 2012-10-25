package delayrepo

import (
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

func TestHeadDelay(t *testing.T) {
	var q Queue
	redis := NewRedis()

	q = NewHeadDelay("hdt", 2, redis)

	next, err := q.NextWakeup()
	assert.Equal(t, err, nil)
	assert.Equal(t, next, 2*time.Second)

	for i := 0; i < 10; i++ {
		q.Push("test1", []byte(fmt.Sprintf("%d", i)))
	}
	next, err = q.NextWakeup()
	assert.Equal(t, err, nil)
	time.Sleep(next)

	q.Push("test1", []byte("10"))
	ret, err := q.Pop()
	assert.Equal(t, err, nil)
	assert.Equal(t, fmt.Sprintf("%+v", ret), "[[48] [49] [50] [51] [52] [53] [54] [55] [56] [57] [49 48]]")
}
