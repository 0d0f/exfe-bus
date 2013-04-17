package gobus

import (
	"encoding/json"
	"github.com/googollee/go-logger"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

type TubeTest1 struct {
	updateCount int
}

func (t *TubeTest1) SetRoute(route RouteCreater) error {
	json := new(JSON)
	return route().Methods("GET").Path("/update").HandlerMethod(json, t, "Update")
}

func (t *TubeTest1) Update(params map[string]string) (int, error) {
	t.updateCount++
	return 0, nil
}

type TubeTest2 struct {
	streamCount int
}

func (t *TubeTest2) SetRoute(route RouteCreater) error {
	json := new(JSON)
	return route().Methods("GET").Path("/stream").HandlerMethod(json, t, "Stream")
}

func (t *TubeTest2) Stream(params map[string]string) (int, error) {
	t.streamCount++
	return 0, nil
}

func TestTube(t *testing.T) {
	l, err := logger.New(logger.Stderr, "test tube")
	if err != nil {
		panic(err)
	}

	bus, err := NewServer("127.0.0.1:12346", l)
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
	time.Sleep(time.Second / 3)

	config := `
	{
	    "bus://update": {"_default": "http://127.0.0.1:12346/update"},
	    "bus://stream": {"_default": "http://127.0.0.1:12346/stream"}
	}`

	var route map[string]map[string]string
	err = json.Unmarshal([]byte(config), &route)
	if err != nil {
		t.Fatal(err)
	}

	table, _ := NewTable(route)
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

	err = tube.SendWithTicket("abc", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, tester1.updateCount, 2)
	assert.Equal(t, tester2.streamCount, 2)
}
