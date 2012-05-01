package gobus

import (
	"time"
	"testing"
	"github.com/simonz05/godis"
)


func TestLastDelayQueue(t *testing.T) {
	interval := 5

	godis := godis.New("", 0, "")
	var queueType int
	q := NewLastDelayQueue("dt", interval, queueType, godis)

	next, _ := q.NextWakeup()
	if next != (time.Duration(interval) * time.Second) {
		t.Fatalf("queue next wakeup should %d seconds, but got: %d", interval, next)
	}

	for i := 0; i < 10; i++ {
		q.Push("googollee", i)
	}
	next, _ = q.NextWakeup()
	time.Sleep(next)

	q.Push("googollee", 10)
	ret, err := q.Pop()
	if err != nil {
		t.Fatalf("Pop error: %s", err)
	}
	if len(ret) != 0 {
		t.Fatalf("Pop should not get anything, but got: %s", ret)
	}

	q.Push("lzh", 0)
	next, _ = q.NextWakeup()
	time.Sleep(next / 2)

	q.Push("lzh", 1)
	next, _ = q.NextWakeup()
	time.Sleep(next)

	ret, err = q.Pop()
	if err != nil {
		t.Fatalf("Pop error: %s", err)
	}
	if len(ret) != 11 {
		t.Fatalf("Pop data error: %s", ret)
	}
	t.Logf("ret: %s", ret)

	ret, err = q.Pop()
	if err != nil {
		t.Fatalf("Pop error: %s", err)
	}
	if len(ret) != 0 {
		t.Fatalf("Pop should not get anything, but got: %s", ret)
	}

	next, _ = q.NextWakeup()
	time.Sleep(next)
	ret, err = q.Pop()
	if err != nil {
		t.Fatalf("Pop error: %s", err)
	}
	if len(ret) != 2 {
		t.Fatalf("Pop should not get anything, but got: %s", ret)
	}
	t.Logf("ret: %s", ret)

	next, _ = q.NextWakeup()
	if next != (time.Duration(interval) * time.Second) {
		t.Fatalf("queue next wakeup should %d seconds, but got: %d", interval, next)
	}
}
