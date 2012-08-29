package gobus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-log"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DelayService interface {
	DeclareQueue() Queue
}

type GobusServer struct {
	jsonServer      *JSONServer
	url             string
	delayDispatcher *DelayDispatcher
}

func NewGobusServer(u string, l *log.Logger) (*GobusServer, error) {
	u_, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	s := NewJSONServer()
	dispatcher := NewDelayDispatcher(l)
	s.SetDispatcher(dispatcher)
	return &GobusServer{
		jsonServer:      s,
		url:             u_.Host,
		delayDispatcher: dispatcher,
	}, nil
}

func (s *GobusServer) Register(service interface{}) error {
	return s.jsonServer.Register(service)
}

func (s *GobusServer) RegisterDelayService(service DelayService) error {
	return s.delayDispatcher.Register(service)
}

func (s *GobusServer) ListenAndServe() error {
	h := &http.Server{
		Addr:    s.url,
		Handler: s.jsonServer,
	}
	go s.delayDispatcher.Serve()
	return h.ListenAndServe()
}

type GobusClient struct {
	httpClient *http.Client
	url        *url.URL
}

func NewGobusClient(u string) (*GobusClient, error) {
	u_, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	c := new(http.Client)
	return &GobusClient{
		httpClient: c,
		url:        u_,
	}, nil
}

func (c *GobusClient) Do(method string, arg interface{}, reply interface{}) error {
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

func (c *GobusClient) Send(method string, arg interface{}) error {
	_, err := c.send(method, arg)
	return err
}

func (c *GobusClient) send(method string, arg interface{}) ([]byte, error) {
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
