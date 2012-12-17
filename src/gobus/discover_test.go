package gobus

import (
	"encoding/json"
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestDispatcher(t *testing.T) {
	config := `
	{
	    "bus://test1/": {"_default": "http://127.0.0.1/test1/"},
	    "bus://test2/sub": {
	    	"_default": "http://127.0.0.1/test2/",
	    	"twitter": "http://127.0.0.2/test2"
	    }
	}`

	var route map[string]map[string]string
	err := json.Unmarshal([]byte(config), &route)
	if err != nil {
		t.Fatal(err)
	}

	dispatcher := NewDispatcher(route)

	{
		_, err := dispatcher.Find("bus://not_exist/", "abc")
		assert.NotEqual(t, err, nil)
	}

	{
		url, err := dispatcher.Find("bus://test1/", "abc")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.1/test1/")
	}

	{
		url, err := dispatcher.Find("bus://test1", "abc")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.1/test1/")
	}

	{
		url, err := dispatcher.Find("bus://test2/sub", "abc")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.1/test2/")
	}

	{
		url, err := dispatcher.Find("bus://test2/sub/", "twitter")
		assert.Equal(t, err, nil)
		assert.Equal(t, url, "http://127.0.0.2/test2/")
	}
}
