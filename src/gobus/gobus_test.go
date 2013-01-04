package gobus

import (
	"fmt"
	"github.com/googollee/go-logger"
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
	route().Methods("POST").Path("/add").HandlerFunc(Must(HandleMethod(json, t, "Add_")))
	route().Methods("GET").Path("/key/{key}").HandlerFunc(Must(HandleMethod(json, t, "CheckKey_")))
}

func (t *gobusTest) Add_(params map[string]string, arg AddArgs) (int, error) {
	return (arg.A + arg.B), nil
}

func (t *gobusTest) CheckKey_(params map[string]string) (string, error) {
	return params["key"], nil
}

func (t *gobusTest) Add(meta *HTTPMeta, args AddArgs, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func (t *gobusTest) CheckKey(meta *HTTPMeta, args AddArgs, reply *string) error {
	*reply = meta.Vars["key"]
	return nil
}

func TestGobus(t *testing.T) {
	addr := "127.0.0.1:1111"

	s, err := NewServer_(addr)
	assert.Equal(t, err, nil)

	test := new(gobusTest)
	err = s.Register(test)
	assert.Equal(t, err, nil)

	go s.ListenAndServe()

	json := new(JSON)
	client := NewClient_(json)

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

func TestOldGobus(t *testing.T) {
	const gobusUrl = "http://127.0.0.1:11111"
	l, err := logger.New(logger.Stderr, "test gobus")
	if err != nil {
		panic(err)
	}
	s, err := NewServer(gobusUrl, l)
	if err != nil {
		t.Fatalf("create gobus server fail: %s", err)
	}

	test := new(gobusTest)
	count, err := s.Register(test)
	if err != nil {
		t.Fatalf("register error: %s", err)
	}
	if count != 2 {
		t.Fatalf("only register %d methods, should be 2", count)
	}

	count, err = s.RegisterName("test", test)
	if err != nil {
		t.Fatalf("register error: %s", err)
	}
	if count != 2 {
		t.Fatalf("only register %d methods, should be 2", count)
	}

	count, err = s.RegisterPath("/test/{key}", test)
	if err != nil {
		t.Fatalf("register error: %s", err)
	}
	if count != 2 {
		t.Fatalf("only register %d methods, should be 2", count)
	}

	go s.ListenAndServe()

	{
		c, err := NewClient(fmt.Sprintf("%s/%s", gobusUrl, "gobusTest"))
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
	}

	{
		c, err := NewClient(fmt.Sprintf("%s/test", gobusUrl))
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
	}

	{
		c, err := NewClient(fmt.Sprintf("%s/test/abc", gobusUrl))
		if err != nil {
			t.Fatalf("create gobus client fail: %s", err)
		}
		var reply string
		err = c.Do("CheckKey", &AddArgs{2, 4}, &reply)
		if err != nil {
			t.Fatalf("call CheckKey error: %s", err)
		}
		if expect, got := "abc", reply; got != expect {
			t.Errorf("expect: %s, got: %s", expect, got)
		}
	}
}
