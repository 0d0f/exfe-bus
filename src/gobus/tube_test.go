package gobus

import (
	"encoding/json"
	"github.com/stretchrcom/testify/assert"
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
	assert.Equal(t, err, nil)
	err = bus.Register(tester2)
	assert.Equal(t, err, nil)

	go bus.ListenAndServe()

	config := `
	{
	    "bus://update": {"_default": "http://127.0.0.1:23333/update"},
	    "bus://stream": {"_default": "http://127.0.0.1:23333/stream"}
	}`

	var route map[string]map[string]string
	err = json.Unmarshal([]byte(config), &route)
	if err != nil {
		t.Fatal(err)
	}

	table := NewTable(route)
	dispatcher := NewDispatcher(table)

	tube := NewTubeClient(dispatcher)
	err = tube.AddService("bus://update", "GET")
	assert.Equal(t, err, nil)
	err = tube.AddService("bus://stream", "GET")
	assert.Equal(t, err, nil)

	err = tube.Send(1)
	assert.Equal(t, err, nil)
	assert.Equal(t, tester1.updateCount, 1)
	assert.Equal(t, tester2.streamCount, 1)

	err = tube.SendWithIdentity("abc", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, tester1.updateCount, 2)
	assert.Equal(t, tester2.streamCount, 2)
}
