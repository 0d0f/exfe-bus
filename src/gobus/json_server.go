package gobus

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"net/http"
	"reflect"
	"strings"
)

type JSONServer struct {
	services map[string]*serviceType
	log      *logger.Logger
}

func NewJSONServer(l *logger.Logger) *JSONServer {
	return &JSONServer{
		services: make(map[string]*serviceType),
		log:      l,
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
	methodName := s.methodName(r)
	subLogger := s.log.Sub(fmt.Sprintf("[%s]", methodName))
	defer func() {
		e := json.NewEncoder(w)
		err := e.Encode(ret)
		if err != nil {
			subLogger.Crit("%s can't encode return(%s) to json: %s", methodName, ret, err)
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

	subLogger.Debug("call with %s", input.Interface())
	output, err := method.call(service.service, &HTTPMeta{r, w, subLogger}, input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		subLogger.Err("failed with input %s: %s", input.Interface(), err)
		ret = fmt.Sprintf("%s", err)
		return
	}
	ret = output.Interface()
	subLogger.Debug("return %s", ret)
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
	Log      *logger.SubLogger
}
