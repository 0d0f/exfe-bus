package gobus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type RouteCreater func() *mux.Route

type Codec interface {
	Mime() string
	Decode(r io.Reader, t reflect.Type) (reflect.Value, error)
	Encode(w io.Writer, v reflect.Value) error
}

type Service interface {
	SetRoute(route RouteCreater)
}

type Server_ struct {
	router *mux.Router
	addr   string
}

func NewServer_(addr string) (*Server_, error) {
	router := mux.NewRouter()
	router.StrictSlash(true)
	return &Server_{
		router: router,
		addr:   addr,
	}, nil
}

func (s *Server_) Register(service Service) error {
	service.SetRoute(func() *mux.Route { return s.router.NewRoute() })
	return nil
}

func (s *Server_) ListenAndServe() error {
	h := &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}
	return h.ListenAndServe()
}

func Must(f http.HandlerFunc, err error) http.HandlerFunc {
	if err != nil {
		panic(err)
	}
	return f
}

//var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfMap = reflect.TypeOf((map[string]string)(nil))

func callMethod(f reflect.Value, args []reflect.Value) (reflect.Value, error) {
	rets := f.Call(args)
	ret, e := rets[0], rets[1].Interface()
	if e != nil {
		return ret, e.(error)
	}
	return ret, nil
}

func HandleMethod(codec Codec, arg interface{}, method string) (http.HandlerFunc, error) {
	t := reflect.TypeOf(arg)
	v := reflect.ValueOf(arg)
	m, ok := t.MethodByName(method)
	if !ok {
		return nil, fmt.Errorf("can't find method")
	}
	if m.Type.NumOut() != 2 {
		return nil, fmt.Errorf("output arg is not 2")
	}
	if m.Type.Out(1) != typeOfError {
		return nil, fmt.Errorf("second output is not error")
	}

	switch m.Type.NumIn() {
	case 3:
		if m.Type.In(1) != typeOfMap {
			return nil, fmt.Errorf("first input is not map[string]string")
		}
		inputType := m.Type.In(2)
		inputPtr := false
		if inputType.Kind() == reflect.Ptr {
			inputType = inputType.Elem()
			inputPtr = true
		}
		return func(w http.ResponseWriter, r *http.Request) {
			input, err := codec.Decode(r.Body, inputType)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			if !inputPtr {
				input = input.Elem()
			}
			ret, err := callMethod(m.Func, []reflect.Value{v, reflect.ValueOf(Params(r)), input})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Header().Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", codec.Mime()))
			err = codec.Encode(w, ret)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			return
		}, nil
	case 2:
		if m.Type.In(1) != typeOfMap {
			return nil, fmt.Errorf("first input is not map[string]string")
		}
		return func(w http.ResponseWriter, r *http.Request) {
			ret, err := callMethod(m.Func, []reflect.Value{v, reflect.ValueOf(Params(r))})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Header().Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", codec.Mime()))
			err = codec.Encode(w, ret)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			return
		}, nil
	}
	return nil, fmt.Errorf("method must have 1 or 2 input arguments")
}

func Params(r *http.Request) map[string]string {
	vars := mux.Vars(r)
	q := r.URL.Query()
	for k, _ := range q {
		vars[k] = q.Get(k)
	}
	return vars
}

type Client_ struct {
	codec      Codec
	httpClient *http.Client
}

func NewClient_(codec Codec) *Client_ {
	return &Client_{
		codec:      codec,
		httpClient: http.DefaultClient,
	}
}

func (c *Client_) Do(urlStr, method string, arg interface{}, reply interface{}) error {
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

///////////

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
