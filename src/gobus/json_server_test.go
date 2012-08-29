package gobus

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type Server struct {
	lastInstance string
	lastMethod   string
}

func (s *Server) Double(meta *HTTPMeta, args int, reply *int) error {
	*reply = args * 2
	return nil
}

func (s *Server) Triple(meta *HTTPMeta, args *int, reply *int) error {
	*reply = *args * 3
	return nil
}

func (s *Server) Error(meta *HTTPMeta, arg int, reply *int) error {
	return fmt.Errorf("inner error")
}

func (s *Server) POST(meta *HTTPMeta, arg int, reply *int) error {
	*reply = arg * 4
	meta.Response.WriteHeader(http.StatusCreated)
	return nil
}

func (s *Server) Dispatch(req *http.Request, instance, method string) error {
	s.lastInstance = instance
	s.lastMethod = method
	return nil
}

func TestJSONServer(t *testing.T) {
	s := NewJSONServer()
	server := new(Server)
	s.Register(server)
	s.SetDispatcher(server)
	h := &http.Server{
		Addr:    "127.0.0.1:1234",
		Handler: s,
	}
	go h.ListenAndServe()

	{
		buf := bytes.NewBufferString("1")
		resp, err := http.Post("http://127.0.0.1:1234/Server?method=Double", "application/json", buf)
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
		resp, err := http.Post("http://127.0.0.1:1234/Server?method=Triple", "application/json", buf)
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
		resp, err := http.Post("http://127.0.0.1:1234/Server?method=Error", "application/json", buf)
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
		resp, err := http.Post("http://127.0.0.1:1234/Server", "application/json", buf)
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

	{
		buf := bytes.NewBufferString("3")
		resp, err := http.Post("http://127.0.0.1:1234/Server?method=NoInstance", "application/json", buf)
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
		if got, expect := string(body), "true\n"; got != expect {
			t.Errorf("expect: (%s), got: (%s)", expect, got)
		}
		if got, expect := server.lastInstance, "Server"; got != expect {
			t.Errorf("expect: %s, got: %s", expect, got)
		}
		if got, expect := server.lastMethod, "NoInstance"; got != expect {
			t.Errorf("expect: %s, got: %s", expect, got)
		}
	}
}
