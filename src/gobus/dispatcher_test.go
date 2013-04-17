package gobus

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

func TestTable(t *testing.T) {
	config := `
	{
	    "bus://test1": {"_default": "http://127.0.0.1/test1"},
	    "bus://test2/sub": {
	    	"_default": "http://127.0.0.1/test2",
	    	"twitter": "http://127.0.0.2/test2",
	    	".*?u123": "http://127.0.0.3/test2"
	    }
	}`

	var route map[string]map[string]string
	err := json.Unmarshal([]byte(config), &route)
	if err != nil {
		t.Fatal(err)
	}

	table, err := NewTable(route)
	assert.Equal(t, err, nil)

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
		url, err := table.Find("bus://test1/sub", "abc")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.1/test1/sub")
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

	{
		url, err := table.Find("bus://test2/sub", "c123,u123")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.3/test2")
	}

}

func TestDispatcher(t *testing.T) {
	l, err := logger.New(logger.Stderr, "test dispatcher")
	if err != nil {
		panic(err)
	}

	const gobusUrl = "127.0.0.1:12345"
	s, _ := NewServer(gobusUrl, l)
	test := new(gobusTest)
	s.Register(test)

	go s.ListenAndServe()
	time.Sleep(time.Second / 3)

	config := `
	{
	    "bus://add": {"_default": "http://127.0.0.1:12345/add"}
	}`

	var route map[string]map[string]string
	err = json.Unmarshal([]byte(config), &route)
	assert.Equal(t, err, nil)

	table, _ := NewTable(route)
	dispatcher := NewDispatcher(table)

	{
		var reply int
		err = dispatcher.Do("bus://add", "POST", AddArgs{2, 4}, &reply)
		assert.Equal(t, err, nil, fmt.Sprintf("error: %s", err))
		assert.Equal(t, reply, 6)
	}

	{
		var reply int
		err = dispatcher.DoWithTicket("abc", "bus://add", "POST", AddArgs{2, 4}, &reply)
		assert.Equal(t, err, nil, fmt.Sprintf("error: %s", err))
		assert.Equal(t, reply, 6)
	}
}
