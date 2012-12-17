package gobus

import (
	"encoding/json"
	"github.com/googollee/go-logger"
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestTable(t *testing.T) {
	config := `
	{
	    "bus://test1": {"_default": "http://127.0.0.1/test1"},
	    "bus://test2/sub": {
	    	"_default": "http://127.0.0.1/test2",
	    	"twitter": "http://127.0.0.2/test2"
	    }
	}`

	var route map[string]map[string]string
	err := json.Unmarshal([]byte(config), &route)
	if err != nil {
		t.Fatal(err)
	}

	table := NewTable(route)

	{
		_, err := table.Find("bus://not_exist", "abc")
		assert.NotEqual(t, err, nil)
	}

	{
		url, err := table.Find("bus://test1", "abc")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.1/test1")
	}

	{
		url, err := table.Find("bus://test1", "abc")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.1/test1")
	}

	{
		url, err := table.Find("bus://test2/sub", "abc")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.1/test2")
	}

	{
		url, err := table.Find("bus://test2/sub", "twitter")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.2/test2")
	}
}

func TestDispatcher(t *testing.T) {
	l, err := logger.New(logger.Stderr, "test gobus")
	if err != nil {
		panic(err)
	}
	s, _ := NewServer(gobusUrl, l)
	test := new(gobusTest)
	s.Register(test)

	go s.ListenAndServe()

	config := `
	{
	    "bus://test1": {"_default": "http://127.0.0.1:12345/gobusTest"}
	}`

	var route map[string]map[string]string
	err = json.Unmarshal([]byte(config), &route)
	if err != nil {
		t.Fatal(err)
	}

	table := NewTable(route)
	dispatcher := NewDispatcher(table)

	{
		var reply int
		err = dispatcher.Do("bus://test1", "Add", AddArgs{2, 4}, &reply)
		if err != nil {
			t.Fatalf("call Add error: %s", err)
		}
		if expect, got := 6, reply; got != expect {
			t.Error("expect: %d, got: %d", expect, got)
		}
	}

	{
		var reply int
		err = dispatcher.DoWithIdentity("abc", "bus://test1", "Add", AddArgs{2, 4}, &reply)
		if err != nil {
			t.Fatalf("call Add error: %s", err)
		}
		if expect, got := 6, reply; got != expect {
			t.Error("expect: %d, got: %d", expect, got)
		}
	}
}