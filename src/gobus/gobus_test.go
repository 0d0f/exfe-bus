package gobus

import (
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"testing"
)

type gobusTest struct {
}

type AddArgs struct {
	A int
	B int
}

func (t *gobusTest) SetRoute(route RouteCreater) {
	json := new(JSON)
	route().Methods("POST").Path("/add").HandlerFunc(Must(Method(json, t, "Add_")))
	route().Methods("GET").Path("/key/{key}").HandlerFunc(Must(Method(json, t, "CheckKey_")))
}

func (t *gobusTest) Add_(params map[string]string, arg AddArgs) (int, error) {
	return (arg.A + arg.B), nil
}

func (t *gobusTest) CheckKey_(params map[string]string) (string, error) {
	return params["key"], nil
}

func TestGobus(t *testing.T) {
	addr := "127.0.0.1:1111"

	s, err := NewServer(addr)
	assert.Equal(t, err, nil)

	test := new(gobusTest)
	err = s.Register(test)
	assert.Equal(t, err, nil)

	go s.ListenAndServe()

	json := new(JSON)
	client := NewClient(json)

	{
		var reply int
		err := client.Do(fmt.Sprintf("http://%s/add", addr), "POST", AddArgs{1, 2}, &reply)
		assert.Equal(t, err, nil)
		assert.Equal(t, reply, 3)
	}

	{
		var reply string
		err := client.Do(fmt.Sprintf("http://%s/key/abcdefg", addr), "GET", nil, &reply)
		assert.Equal(t, err, nil)
		assert.Equal(t, reply, "abcdefg")
	}
}
