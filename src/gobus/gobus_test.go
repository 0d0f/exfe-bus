package gobus

import (
	"fmt"
	"github.com/googollee/go-log"
	"testing"
	"time"
)

type gobusTest struct {
	data []MyInt
}

type AddArgs struct {
	A int
	B int
}

func (t *gobusTest) Add(meta *HTTPMeta, args AddArgs, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func (t *gobusTest) DeclareQueue() Queue {
	return NewIntervalQueue("test", interval, redis, []MyInt{})
}

func (t *gobusTest) IntDelay(args []MyInt) error {
	t.data = append(t.data, args...)
	return nil
}

const gobusUrl = "http://127.0.0.1:12345"

func TestGobus(t *testing.T) {
	l, err := log.New(log.Stderr, "test gobus", log.LstdFlags)
	if err != nil {
		panic(err)
	}
	s, err := NewGobusServer(gobusUrl, l)
	if err != nil {
		t.Fatalf("create gobus server fail: %s", err)
	}
	test := new(gobusTest)
	test.data = make([]MyInt, 0)
	s.Register(test)
	s.RegisterDelayService(test)
	go s.ListenAndServe()

	c, err := NewGobusClient(fmt.Sprintf("%s/%s", gobusUrl, "gobusTest"))
	if err != nil {
		t.Fatalf("create gobus client fail: %s", err)
	}
	var reply int
	err = c.Do("Add", &AddArgs{2, 4}, &reply)
	if err != nil {
		t.Fatalf("call Add error: %s", err)
	}
	if expect, got := 6, reply; got != expect {
		t.Error("expect: %d, got: %d", expect, got)
	}

	for i := 0; i < 10; i++ {
		err = c.Send("IntDelay", MyInt(i))
	}

	time.Sleep(time.Duration(interval) * time.Second * 2)

	if got, expect := len(test.data), 10; got != expect {
		t.Errorf("expect: %d, got: %d", expect, got)
	}
	t.Logf("%+v", test.data)
	for i := 0; i < 10; i++ {
		if got, expect := int(test.data[i]), i; got != expect {
			t.Errorf("index %d expect %d, got: %d", i, expect, got)
		}
	}
}
