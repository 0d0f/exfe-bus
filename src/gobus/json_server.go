package gobus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type JSONServerDispatcher interface {
	Dispatch(req *http.Request, instance, method string) error
}

type JSONServer struct {
	services   map[string]*serviceType
	dispatcher JSONServerDispatcher
}

func NewJSONServer() *JSONServer {
	return &JSONServer{
		services:   make(map[string]*serviceType),
		dispatcher: nil,
	}
}

func (s *JSONServer) Register(arg interface{}) error {
	t := reflect.TypeOf(arg)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	s.services[t.Name()] = newServiceType(arg)
	return nil
}

func (s *JSONServer) SetDispatcher(dispatcher JSONServerDispatcher) {
	s.dispatcher = dispatcher
}

func (s *JSONServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")

	var ret interface{}
	defer func() {
		e := json.NewEncoder(w)
		err := e.Encode(ret)
		if err != nil {
			panic(err)
		}
	}()

	if len(paths) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		ret = fmt.Sprintf("can't find service to handle %s", r.URL.Path)
		return
	}

	serviceName := paths[len(paths)-1]
	service, ok := s.services[serviceName]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		ret = fmt.Sprintf("can't find service %s", serviceName)
		return
	}
	methodName := r.URL.Query().Get("method")
	if methodName == "" {
		methodName = r.Method
	}
	method, ok := service.methods[methodName]
	if !ok {
		if s.dispatcher != nil {
			err := s.dispatcher.Dispatch(r, serviceName, methodName)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ret = fmt.Sprintf("dispatch method %s in service %s error: %s", methodName, serviceName, err)
				return
			}
			ret = true
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		ret = fmt.Sprintf("can't find method %s in service %s", methodName, serviceName)
		return
	}
	input, err := method.getInput(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ret = fmt.Sprintf("post data can't call method %s: %s", methodName, err)
		return
	}
	output, err := method.call(service.service, &HTTPMeta{r, w}, input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ret = fmt.Sprintf("%s", err)
		return
	}
	ret = output.Interface()
}

type Response interface {
	Header() http.Header
	WriteHeader(int)
}

type HTTPMeta struct {
	Request  *http.Request
	Response Response
}
