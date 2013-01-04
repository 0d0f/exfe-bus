package gobus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
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

type Server struct {
	router *mux.Router
	addr   string
}

func NewServer(addr string) (*Server, error) {
	router := mux.NewRouter()
	router.StrictSlash(true)
	return &Server{
		router: router,
		addr:   addr,
	}, nil
}

func (s *Server) Register(service Service) error {
	service.SetRoute(func() *mux.Route { return s.router.NewRoute() })
	return nil
}

func (s *Server) ListenAndServe() error {
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

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
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
