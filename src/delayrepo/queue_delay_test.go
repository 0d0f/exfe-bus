package queue

import (
	"testing"
	"time"
)

func TestTailDelayQueue(t *testing.T) {
	var q Queue
	q = NewTailDelayQueue("tdt", interval, redis, []testData{})

	next, _ := q.NextWakeup()
	if next != duration {
		t.Fatalf("queue next wakeup should %v, but got: %v", duration, next)
	}

	for i := 0; i < 10; i++ {
		q.Push(testData{"test1", i})
	}
	next, _ = q.NextWakeup()
	time.Sleep(next)

	q.Push(testData{"test1", 10})
	ret, err := q.Pop()
	if err != nil {
		t.Fatalf("Pop error: %s", err)
	}
	if ret != nil {
		t.Fatalf("Pop should get nothing, but got: %v", ret)
	}

	q.Push(testData{"test2", 0})
	next, _ = q.NextWakeup()
	time.Sleep(next / 2)

	q.Push(testData{"test2", 1})
	next, _ = q.NextWakeup()
	time.Sleep(next)

	ret, err = q.Pop()
	if err != nil {
		t.Fatalf("Pop error: %s", err)
	}
	t.Logf("ret: %v", ret)
	if len(ret.([]testData)) != 11 {
		t.Fatalf("Pop data error: %v", ret)
	}

	ret, err = q.Pop()
	if err != nil {
		t.Fatalf("Pop error: %s", err)
	}
	if ret != nil {
		t.Fatalf("Pop should get nothing, but got: %s", ret)
	}

	next, _ = q.NextWakeup()
	time.Sleep(next)
	ret, err = q.Pop()
	if err != nil {
		t.Fatalf("Pop error: %s", err)
	}
	t.Logf("ret: %s", ret)
	if len(ret.([]testData)) != 2 {
		t.Fatalf("Pop should not get anything, but got: %s", ret)
	}

	next, _ = q.NextWakeup()
	if next != duration {
		t.Fatalf("queue next wakeup should %v, but got: %v", duration, next)
	}
}

func TestHeadDelayQueue(t *testing.T) {
	var q Queue
	q = NewHeadDelayQueue("hdt", interval, redis, []testData{})

	next, _ := q.NextWakeup()
	if next != duration {
		t.Fatalf("queue next wakeup should %v, but got: %v", duration, next)
	}

	for i := 0; i < 10; i++ {
		q.Push(testData{"test1", i})
	}
	next, _ = q.NextWakeup()
	time.Sleep(next)

	q.Push(testData{"test1", 10})
	ret, err := q.Pop()
	if err != nil {
		t.Fatalf("Pop error: %s", err)
	}
	if ret == nil {
		t.Fatalf("Pop should get anything, but got nothing")
	}
}
