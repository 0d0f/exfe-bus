package gobus

import (
	"bytes"
	"fmt"
	"github.com/googollee/go-log"
	"net/http"
	"testing"
	"time"
)

type MyInt int

func (m MyInt) KeyForQueue() string {
	return fmt.Sprintf("%d", m/10)
}

type Delay struct {
	data []MyInt
}

func newDelay() *Delay {
	return &Delay{
		data: make([]MyInt, 0, 0),
	}
}

func (d *Delay) DeclareQueue() Queue {
	return NewIntervalQueue("delay_test", interval, redis, []MyInt{})
}

func (d *Delay) Delay(args []MyInt) error {
	d.data = append(d.data, args...)
	return nil
}

func TestDelayDispatcher(t *testing.T) {
	l, err := log.New(log.Stderr, "gobus test", log.LstdFlags)
	if err != nil {
		panic(err)
	}
	d := NewDelayDispatcher(l)
	s := newDelay()
	err = d.Register(s)
	if err != nil {
		panic(err)
	}
	go d.Serve()

	for i := 0; i < 10; i++ {
		buf := bytes.NewBufferString(fmt.Sprintf("%d", i))
		req, err := http.NewRequest("Delay", "http://127.0.0.1:1234", buf)
		if err != nil {
			panic(err)
		}
		err = d.Dispatch(req, "Delay", "Delay")
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Duration(interval) * time.Second * 2)

	if got, expect := len(s.data), 10; got != expect {
		t.Errorf("expect: %d, got: %d", expect, got)
	}
	t.Logf("%+v", s.data)
	for i := 0; i < 10; i++ {
		if got, expect := int(s.data[i]), i; got != expect {
			t.Errorf("index %d expect %d, got: %d", i, expect, got)
		}
	}
}

type RetryDelay struct {
	retry bool
	data  []MyInt
}

func newRetryDelay() *RetryDelay {
	return &RetryDelay{
		retry: true,
		data:  make([]MyInt, 0, 0),
	}
}

func (d *RetryDelay) DeclareQueue() Queue {
	return NewIntervalQueue("retrydelay_test", interval, redis, []MyInt{})
}

func (d *RetryDelay) Delay(args []MyInt) error {
	if d.retry {
		d.retry = false
		return NewRetryError(fmt.Errorf("need retry"))
	}
	d.data = append(d.data, args...)
	return nil
}

func TestRetryDelayDispatcher(t *testing.T) {
	l, err := log.New(log.Stderr, "gobus test", log.LstdFlags)
	d := NewDelayDispatcher(l)
	s := newRetryDelay()
	err = d.Register(s)
	if err != nil {
		panic(err)
	}
	go d.Serve()

	for i := 0; i < 10; i++ {
		buf := bytes.NewBufferString(fmt.Sprintf("%d", i))
		req, err := http.NewRequest("Delay", "http://127.0.0.1:1234", buf)
		if err != nil {
			panic(err)
		}
		err = d.Dispatch(req, "RetryDelay", "Delay")
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Duration(interval) * time.Second * 13 / 10)
	if got, expect := len(s.data), 0; got != expect {
		t.Errorf("expect: %d, got: %d", expect, got)
	}
	if got, expect := s.retry, false; got != expect {
		t.Errorf("expect: %d, got: %d", expect, got)
	}

	time.Sleep(time.Duration(interval) * time.Second * 13 / 10)
	if got, expect := len(s.data), 10; got != expect {
		t.Errorf("expect: %d, got: %d", expect, got)
	}
	t.Logf("%+v", s.data)
	for i := 0; i < 10; i++ {
		if got, expect := int(s.data[i]), i; got != expect {
			t.Errorf("index %d expect %d, got: %d", i, expect, got)
		}
	}
}
