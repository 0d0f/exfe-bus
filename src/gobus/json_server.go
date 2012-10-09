package gobus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type JSONServer struct {
	services map[string]*serviceType
}

func NewJSONServer() *JSONServer {
	return &JSONServer{
		services: make(map[string]*serviceType),
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

func (s *JSONServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ret interface{}
	defer func() {
		e := json.NewEncoder(w)
		err := e.Encode(ret)
		if err != nil {
			panic(err)
		}
	}()

	methodName := s.methodName(r)
	service, method, err := s.findMethod(r, methodName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ret = err
		return
	}

	input, err := method.getInput(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ret = fmt.Sprintf("data can't call method %s: %s", methodName, err)
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

func (s *JSONServer) methodName(r *http.Request) string {
	methodName := r.URL.Query().Get("method")
	if methodName == "" {
		methodName = r.Method
	}
	return methodName
}

func (s *JSONServer) findMethod(r *http.Request, methodName string) (*serviceType, *methodType, error) {
	paths := strings.Split(r.URL.Path, "/")

	if len(paths) < 1 {
		return nil, nil, fmt.Errorf("can't find service to handle %s", r.URL.Path)
	}

	serviceName := paths[len(paths)-1]
	service, ok := s.services[serviceName]
	if !ok {
		return nil, nil, fmt.Errorf("can't find service %s", serviceName)
	}

	method, ok := service.methods[methodName]
	if !ok {
		return nil, nil, fmt.Errorf("can't find method %s in service %s", methodName, serviceName)
	}

	return service, method, nil
}

type Response interface {
	Header() http.Header
	WriteHeader(int)
}

type HTTPMeta struct {
	Request  *http.Request
	Response Response
}
