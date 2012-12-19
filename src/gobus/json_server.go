package gobus

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"net/http"
	"reflect"
)

type jsonServer struct {
	name    string
	service *serviceType
	log     *logger.Logger
}

func newJSONServer(log *logger.Logger, server interface{}) *jsonServer {
	t := reflect.TypeOf(server)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()

	service := newServiceType(server)
	return &jsonServer{
		service: service,
		name:    name,
		log:     log,
	}
}

func (s *jsonServer) Name() string {
	return s.name
}

func (s *jsonServer) MethodCount() int {
	return len(s.service.methods)
}

func (s *jsonServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ret interface{}
	methodName := s.methodName(r)
	subLogger := s.log.Sub(fmt.Sprintf("[%s|%s]", r.URL.Path, methodName))
	defer func() {
		e := json.NewEncoder(w)
		err := e.Encode(ret)
		if err != nil {
			subLogger.Crit("%s encode return error(%s) with %s", methodName, err, ret)
		}
	}()
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		subLogger.Crit("%s panic: %s", methodName, r)
		ret = fmt.Sprintf("panic: %s", r)
	}()

	method, ok := s.service.methods[methodName]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		ret = fmt.Errorf("can't find method %s in service %s", methodName, s.Name())
		return
	}

	input, err := method.getInput(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ret = fmt.Sprintf("data can't call method %s: %s", methodName, err)
		return
	}

	inputElem := input.Interface()
	if input.Kind() == reflect.Ptr {
		inputElem = input.Elem().Interface()
	}
	subLogger.Debug("call with %+v", inputElem)
	output, err := method.call(s.service.service, newMeta(r, w, subLogger), input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		subLogger.Err("%s, with input %+v", err, inputElem)
		ret = err.Error()
		return
	}
	ret = output.Elem().Interface()
	subLogger.Debug("return %+v", ret)
}

func (s *jsonServer) methodName(r *http.Request) string {
	methodName := r.URL.Query().Get("method")
	if methodName == "" {
		methodName = r.Method
	}
	return methodName
}
