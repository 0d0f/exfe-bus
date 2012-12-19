package gobus

import (
	"bytes"
	"fmt"
	"github.com/googollee/go-logger"
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

type Resp struct {
	buf     *bytes.Buffer
	retCode int
	header  http.Header
}

func newResp() *Resp {
	return &Resp{
		buf:     bytes.NewBuffer(nil),
		retCode: http.StatusOK,
		header:  make(http.Header),
	}
}

func (r *Resp) Header() http.Header {
	return r.header
}

func (r *Resp) Write(p []byte) (int, error) {
	return r.buf.Write(p)
}

func (r *Resp) WriteHeader(code int) {
	r.retCode = code
}

func TestJSONServer(t *testing.T) {
	l, err := logger.New(logger.Stderr, "test gobus")
	if err != nil {
		panic(err)
	}
	server := new(TestServer)
	s := newJSONServer(l, server)
	count := s.MethodCount()
	if count != 4 {
		t.Fatalf("register %d methods, should be 4", count)
	}
	name := s.Name()
	if name != "TestServer" {
		t.Fatalf("server name %s, should be TestServer", name)
	}

	{
		buf := bytes.NewBufferString("1")
		r, err := http.NewRequest("POST", "http://127.0.0.1:1234?method=Double", buf)
		if err != nil {
			t.Fatalf("new request error: %s", err)
		}
		w := newResp()
		s.ServeHTTP(w, r)
		if w.retCode != 200 {
			t.Errorf("http should respond 200, got: %s", w.retCode)
		}
		if got, expect := w.buf.String(), "2\n"; got != expect {
			t.Errorf("expect: (%s), got: (%s)", expect, got)
		}
	}

	{
		buf := bytes.NewBufferString("2")
		r, err := http.NewRequest("POST", "http://127.0.0.1:1234?method=Triple", buf)
		if err != nil {
			t.Fatalf("new request error: %s", err)
		}
		w := newResp()
		s.ServeHTTP(w, r)
		if w.retCode != 200 {
			t.Errorf("http should respond 200, got: %s", w.retCode)
		}
		if got, expect := w.buf.String(), "6\n"; got != expect {
			t.Errorf("expect: (%s), got: (%s)", expect, got)
		}
	}

	{
		buf := bytes.NewBufferString("2")
		r, err := http.NewRequest("POST", "http://127.0.0.1:1234?method=Error", buf)
		if err != nil {
			t.Fatalf("new request error: %s", err)
		}
		w := newResp()
		s.ServeHTTP(w, r)
		if w.retCode != 500 {
			t.Errorf("http should respond 500, got: %s", w.retCode)
		}
		if got, expect := w.buf.String(), "\"inner error\"\n"; got != expect {
			t.Errorf("expect: (%s), got: (%s)", expect, got)
		}
	}

	{
		buf := bytes.NewBufferString("3")
		r, err := http.NewRequest("POST", "http://127.0.0.1:1234", buf)
		if err != nil {
			t.Fatalf("new request error: %s", err)
		}
		w := newResp()
		s.ServeHTTP(w, r)
		if w.retCode != 201 {
			t.Errorf("http should respond 201, got: %s", w.retCode)
		}
		if got, expect := w.buf.String(), "12\n"; got != expect {
			t.Errorf("expect: (%s), got: (%s)", expect, got)
		}
	}
}
