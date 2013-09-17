package gobus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-rest"
	old "github.com/googollee/go-rest/old_style"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/fcgi"
	"reflect"
)

type RouteCreater func() *Route

type Codec interface {
	Mime() string
	Decode(r io.Reader, t reflect.Type) (reflect.Value, error)
	Encode(w io.Writer, v reflect.Value) error
}

type Service interface {
	SetRoute(route RouteCreater) error
}

type Server struct {
	router   *mux.Router
	addr     string
	fallback *rest.Rest
}

func NewServer(addr string) (*Server, error) {
	router := mux.NewRouter()
	router.StrictSlash(true)
	return &Server{
		router: router,
		addr:   addr,
	}, nil
}

func (s *Server) RegisterPrefix(prefix string, handler http.Handler) error {
	s.router.PathPrefix(prefix).Handler(handler)
	return nil
}

func (s *Server) RegisterRestful(service interface{}) error {
	handler, err := old.New(service)
	if err != nil {
		return err
	}
	s.router.PathPrefix(handler.Prefix()).Handler(handler)
	return nil
}

func (s *Server) RegisterFallback(fallback *rest.Rest) {
	s.fallback = fallback
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler mux.RouteMatch
	ok := s.router.Match(r, &handler)
	if ok {
		s.router.ServeHTTP(w, r)
		return
	}
	if s.fallback != nil {
		s.fallback.ServeHTTP(w, r)
	}
}

func (s *Server) ListenAndServe() error {
	h := &http.Server{
		Addr:    s.addr,
		Handler: s,
	}
	return h.ListenAndServe()
}

func (s *Server) ListenFCGI() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return fcgi.Serve(l, s.router)
}

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfMap = reflect.TypeOf((map[string]string)(nil))

type Client struct {
	codec      Codec
	httpClient *http.Client
}

func NewClient(codec Codec) *Client {
	return &Client{
		codec:      codec,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) Do(urlStr, method string, arg interface{}, reply interface{}) error {
	var req *http.Request
	var err error
	if arg != nil {
		reqReader := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(reqReader)
		err = encoder.Encode(arg)
		if err != nil {
			return err
		}
		req, err = http.NewRequest(method, urlStr, reqReader)
	} else {
		req, err = http.NewRequest(method, urlStr, nil)
	}
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", c.codec.Mime()))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(reply)
	if err != nil {
		return err
	}
	return nil
}
