package gobus

import (
	"testing"
)

type TubeTest1 struct {
	updateCount int
}

func (t *TubeTest1) SetRoute(route RouteCreater) {
	json := new(JSON)
	route().Methods("GET").Path("/update").HandlerFunc(Must(Method(json, t, "Update")))
}

func (t *TubeTest1) Update(params map[string]string) (int, error) {
	t.updateCount++
	return 0, nil
}

type TubeTest2 struct {
	streamCount int
}

func (t *TubeTest2) SetRoute(route RouteCreater) {
	json := new(JSON)
	route().Methods("GET").Path("/stream").HandlerFunc(Must(Method(json, t, "Stream")))
}

func (t *TubeTest2) Stream(params map[string]string) (int, error) {
	t.streamCount++
	return 0, nil
}

func TestTube(t *testing.T) {
	bus, err := NewServer("127.0.0.1:23333")
	if err != nil {
		t.Fatalf("create gobus server fail: %s", err)
	}
	tester1 := new(TubeTest1)
	tester2 := new(TubeTest2)
	err = bus.Register(tester1)
	if err != nil {
		t.Fatalf("register failed: %s", err)
	}
	err = bus.Register(tester2)
	if err != nil {
		t.Fatalf("register failed: %s", err)
	}

	go bus.ListenAndServe()

	tube := NewTubeClient("TubeTest")
	if got, expect := tube.Name(), "TubeTest"; got != expect {
		t.Errorf("expect: %s, got: %s", expect, got)
	}
	err = tube.AddService("http://127.0.0.1:23333/update", "GET")
	if err != nil {
		t.Fatalf("add service tester1 failed: %s", err)
	}
	err = tube.AddService("http://127.0.0.1:23333/stream", "GET")
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
