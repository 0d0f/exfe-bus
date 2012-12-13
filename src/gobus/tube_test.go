package gobus

import (
	"github.com/googollee/go-logger"
	"testing"
)

type TubeTest1 struct {
	updateCount int
}

func (t *TubeTest1) Update(meta *HTTPMeta, arg *int, reply *int) error {
	t.updateCount++
	return nil
}

type TubeTest2 struct {
	streamCount int
}

func (t *TubeTest2) Stream(meta *HTTPMeta, arg *int, reply *int) error {
	t.streamCount++
	return nil
}

func TestTube(t *testing.T) {
	l, err := logger.New(logger.Stderr, "test gobus")
	if err != nil {
		panic(err)
	}
	bus, err := NewServer("http://127.0.0.1:23333", l)
	if err != nil {
		t.Fatalf("create gobus server fail: %s", err)
	}
	tester1 := new(TubeTest1)
	tester2 := new(TubeTest2)
	_, err = bus.Register(tester1)
	if err != nil {
		t.Fatalf("register failed: %s", err)
	}
	_, err = bus.Register(tester2)
	if err != nil {
		t.Fatalf("register failed: %s", err)
	}

	go bus.ListenAndServe()

	tube := NewTubeClient("TubeTest")
	if got, expect := tube.Name(), "TubeTest"; got != expect {
		t.Errorf("expect: %s, got: %s", expect, got)
	}
	err = tube.AddService("http://127.0.0.1:23333/TubeTest1", "Update")
	if err != nil {
		t.Fatalf("add service tester1 failed: %s", err)
	}
	err = tube.AddService("http://127.0.0.1:23333/TubeTest2", "Stream")
	if err != nil {
		t.Fatalf("add service tester2 failed: %s", err)
	}
	err = tube.Send(1)
	if err != nil {
		t.Fatalf("send failed: %s", err)
	}

	if got, expect := tester1.updateCount, 1; got != expect {
		t.Errorf("expect: %d, got: %d", expect, got)
	}
	if got, expect := tester2.streamCount, 1; got != expect {
		t.Errorf("expect: %d, got: %d", expect, got)
	}
}
