package gobus

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type TestServer struct {
}

func (s *TestServer) Double(meta *HTTPMeta, args int, reply *int) error {
	*reply = args * 2
	return nil
}

func (s *TestServer) Triple(meta *HTTPMeta, args *int, reply *int) error {
	*reply = *args * 3
	return nil
}

func (s *TestServer) Error(meta *HTTPMeta, arg int, reply *int) error {
	return fmt.Errorf("inner error")
}

func (s *TestServer) POST(meta *HTTPMeta, arg int, reply *int) error {
	*reply = arg * 4
	meta.Response.WriteHeader(http.StatusCreated)
	return nil
}

func TestJSONServer(t *testing.T) {
	s := NewJSONServer()
	server := new(TestServer)
	s.Register(server)
	h := &http.Server{
		Addr:    "127.0.0.1:1234",
		Handler: s,
	}
	go h.ListenAndServe()

	{
		buf := bytes.NewBufferString("1")
		resp, err := http.Post("http://127.0.0.1:1234/TestServer?method=Double", "application/json", buf)
		if err != nil {
			t.Fatalf("http post error: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("http should respond 200, got: %s", resp.Status)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("read body error: %s", err)
		}
		if got, expect := string(body), "2\n"; got != expect {
			t.Errorf("expect: (%s), got: (%s)", expect, got)
		}
	}

	{
		buf := bytes.NewBufferString("2")
		resp, err := http.Post("http://127.0.0.1:1234/TestServer?method=Triple", "application/json", buf)
		if err != nil {
			t.Fatalf("http post error: %s", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("http should respond 200, got: %s", resp.Status)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("read body error: %s", err)
		}
		if got, expect := string(body), "6\n"; got != expect {
			t.Errorf("expect: (%s), got: (%s)", expect, got)
		}
	}

	{
		buf := bytes.NewBufferString("2")
		resp, err := http.Post("http://127.0.0.1:1234/TestServer?method=Error", "application/json", buf)
		if err != nil {
			t.Fatalf("http post error: %s", err)
		}
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("http should respond 500, got: %s", resp.Status)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("read body error: %s", err)
		}
		if got, expect := string(body), "\"inner error\"\n"; got != expect {
			t.Errorf("expect: (%s), got: (%s)", expect, got)
		}
	}

	{
		buf := bytes.NewBufferString("3")
		resp, err := http.Post("http://127.0.0.1:1234/TestServer", "application/json", buf)
		if err != nil {
			t.Fatalf("http post error: %s", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("http should respond 201, got: %s", resp.Status)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("read body error: %s", err)
		}
		if got, expect := string(body), "12\n"; got != expect {
			t.Errorf("expect: (%s), got: (%s)", expect, got)
		}
	}
}
