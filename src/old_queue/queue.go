package queue

import (
	"errors"
	"time"
)

var EmptyQueueError = errors.New("Empty queue.")
var QueueChangedError = errors.New("Queue changed before pop")
var QueueFullError = errors.New("Queue is full, wait and try again")

type QueueData interface {
	KeyForQueue() string
}

type Queue interface {
	Push(data QueueData) error
	Pop() (interface{}, error)
	NextWakeup() (time.Duration, error)
}
