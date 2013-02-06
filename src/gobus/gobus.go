package gobus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
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
	router *mux.Router
	addr   string
	log    *logger.Logger
}

func NewServer(addr string, log *logger.Logger) (*Server, error) {
	router := mux.NewRouter()
	router.StrictSlash(true)
	return &Server{
		router: router,
		addr:   addr,
		log:    log,
	}, nil
}

func (s *Server) Register(service Service) error {
	return service.SetRoute(func() *Route {
		return &Route{s.router.NewRoute(), s}
	})
}

func (s *Server) ListenAndServe() error {
	h := &http.Server{
		Addr:    s.addr,
		Handler: s.router,
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
