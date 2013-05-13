package gobus

import (
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

type gobusTest struct {
}

type AddArgs struct {
	A int
	B int
}

func (t *gobusTest) SetRoute(route RouteCreater) error {
	json := new(JSON)
	err := route().Methods("POST").Path("/add").HandlerMethod(json, t, "Add")
	if err != nil {
		return err
	}
	err = route().Methods("GET").Path("/key/{key}").HandlerMethod(json, t, "CheckKey")
	if err != nil {
		return err
	}
	err = route().Queries("method", "Check").Path("/key").HandlerMethod(json, t, "Check")
	if err != nil {
		return err
	}
	return nil
}

func (t *gobusTest) Add(params map[string]string, arg AddArgs) (int, error) {
	return (arg.A + arg.B), nil
}

func (t *gobusTest) CheckKey(params map[string]string) (string, error) {
	return params["key"], nil
}

func (t *gobusTest) Check(params map[string]string, key string) (string, error) {
	return fmt.Sprintf("method:%s", key), nil
}

func TestGobus(t *testing.T) {
	addr := "127.0.0.1:1111"
	s, err := NewServer(addr)
	assert.Equal(t, err, nil)

	test := new(gobusTest)
	err = s.Register(test)
	assert.Equal(t, err, nil)

	go s.ListenAndServe()
	time.Sleep(time.Second / 3)

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

	{
		var reply string
		err := client.Do(fmt.Sprintf("http://%s/key?method=Check", addr), "POST", "abcde", &reply)
		assert.Equal(t, err, nil)
		assert.Equal(t, reply, "method:abcde")
	}
}
