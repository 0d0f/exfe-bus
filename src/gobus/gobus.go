package gobus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Response interface {
	Header() http.Header
	WriteHeader(int)
}

type HTTPMeta struct {
	Request  *http.Request
	Response Response
	Log      *logger.SubLogger
	Vars     map[string]string
}

func newMeta(req *http.Request, resp Response, log *logger.SubLogger) *HTTPMeta {
	return &HTTPMeta{
		Request:  req,
		Response: resp,
		Log:      log,
		Vars:     mux.Vars(req),
	}
}

type Server struct {
	router *mux.Router
	url    string
	log    *logger.Logger
}

func NewServer(u string, l *logger.Logger) (*Server, error) {
	u_, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	router := mux.NewRouter()
	router.StrictSlash(true)
	return &Server{
		router: router,
		url:    u_.Host,
		log:    l,
	}, nil
}

func (s *Server) Register(service interface{}) (int, error) {
	server := newJSONServer(s.log, service)
	s.router.Handle(fmt.Sprintf("/%s", server.Name()), server)
	return server.MethodCount(), nil
}

func (s *Server) RegisterName(name string, service interface{}) (int, error) {
	server := newJSONServer(s.log, service)
	s.router.Handle(fmt.Sprintf("/%s", name), server)
	return server.MethodCount(), nil
}

func (s *Server) RegisterPath(path string, service interface{}) (int, error) {
	server := newJSONServer(s.log, service)
	s.router.Handle(path, server)
	return server.MethodCount(), nil
}

func (s *Server) ListenAndServe() error {
	h := &http.Server{
		Addr:    s.url,
		Handler: s.router,
	}
	return h.ListenAndServe()
}

type Client struct {
	httpClient *http.Client
	url        *url.URL
}

func NewClient(u string) (*Client, error) {
	u_, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	c := new(http.Client)
	return &Client{
		httpClient: c,
		url:        u_,
	}, nil
}

func (c *Client) Do(method string, arg interface{}, reply interface{}) error {
	respBody, err := c.send(method, arg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, reply)
	if err != nil {
		return fmt.Errorf("can't unmarshal response to reply: %s", err)
	}
	return nil
}

func (c *Client) send(method string, arg interface{}) ([]byte, error) {
	b, err := json.Marshal(arg)
	if err != nil {
		return nil, fmt.Errorf("can't marshal arg to json: %s", err)
	}
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest(method, c.url.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("http request fail: %s", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: %s", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response error: %s", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http response fail: (%s) %s", resp.Status, string(respBody))
	}
	return respBody, nil
}
