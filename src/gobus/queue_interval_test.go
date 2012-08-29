package gobus

import (
	"testing"
	"time"
)

func TestIntervalQueue(t *testing.T) {
	var q Queue
	q = NewIntervalQueue("iqt", interval, redis, []testData{})

	next, _ := q.NextWakeup()
	if got, expect := int((next+time.Second/2)/time.Second), interval; got != expect {
		t.Fatalf("queue next wakeup should %v, but got: %v", expect, got)
	}

	q.Push(testData{"test1", 1})
	q.Push(testData{"test2", 2})

	time.Sleep(next / 2)

	next, _ = q.NextWakeup()
	if got, expect := int((next+time.Second/2)/time.Second), interval/2; got != expect {
		t.Fatalf("queue next wakeup should %v, but got: %v", expect, got)
	}

	d, _ := q.Pop()
	data, ok := d.([]testData)
	if !ok {
		t.Fatalf("pop should be []testData")
	}
	if expect, got := 2, len(data); expect != got {
		t.Fatalf("pop data length should: %d, but got: %d", expect, got)
	}
	if expect, got := "test1", data[0].Id; got != expect {
		t.Errorf("data[0].Id expect: %s, got: %s", expect, got)
	}
	if expect, got := "test2", data[1].Id; got != expect {
		t.Errorf("data[1].Id expect: %s, got: %s", expect, got)
	}
}
