package gobus

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type delayMethod struct {
	name   string
	method reflect.Value
	arg    reflect.Type
}

func newDelayMethod(service reflect.Type) *delayMethod {
	for i, n := 0, service.NumMethod(); i < n; i++ {
		m := service.Method(i)

		if m.PkgPath != "" {
			// Method must be exported.
			continue
		}

		if !methodOutIsError(&m.Type) {
			continue
		}

		if m.Type.NumIn() != 2 {
			continue
		}

		argType := m.Type.In(1)
		if argType.Kind() != reflect.Slice {
			continue
		}
		if !isExportedOrBuiltinType(argType.Elem()) {
			continue
		}
		return &delayMethod{
			name:   m.Name,
			method: m.Func,
			arg:    argType,
		}

	}
	return nil
}

type serviceMethod struct {
	service reflect.Value
	queue   Queue
	method  *delayMethod
	locker  sync.Locker
}

func newServiceMethod(service DelayService) (*serviceMethod, error) {
	method := newDelayMethod(reflect.TypeOf(service))
	if method == nil {
		return nil, fmt.Errorf("can't find delay method")
	}
	return &serviceMethod{
		service: reflect.ValueOf(service),
		queue:   service.DeclareQueue(),
		method:  method,
		locker:  new(sync.Mutex),
	}, nil
}

type DelayDispatcher struct {
	services map[string]*serviceMethod
	logger   *logger.Logger
}

func NewDelayDispatcher(l *logger.Logger) *DelayDispatcher {
	return &DelayDispatcher{
		services: make(map[string]*serviceMethod),
		logger:   l,
	}
}

func (d *DelayDispatcher) Register(service DelayService) error {
	t := reflect.TypeOf(service)
	s, err := newServiceMethod(service)
	if err != nil {
		return err
	}
	name := t.Name()
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	}
	d.services[name] = s

	return nil
}

func (d *DelayDispatcher) Dispatch(r *http.Request, service, method string) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read request error: %s", err)
	}

	s, ok := d.services[service]
	if !ok {
		return fmt.Errorf("can't find service %s", service)
	}
	if method != s.method.name {
		return fmt.Errorf("can't find method %s in service %s", method, service)
	}

	m := s.method
	arg := m.arg.Elem()
	var ret reflect.Value
	if arg.Kind() == reflect.Ptr {
		ret = reflect.New(arg.Elem())
	} else {
		ret = reflect.New(arg)
	}
	err = json.Unmarshal(body, ret.Interface())
	if err != nil {
		return fmt.Errorf("decode json error: %s", err)
	}
	if arg.Kind() != reflect.Ptr {
		ret = ret.Elem()
	}
	data, ok := ret.Interface().(QueueData)
	if !ok {
		return fmt.Errorf("data can't convert to QueueData")
	}
	s.locker.Lock()
	err = s.queue.Push(data)
	s.locker.Unlock()
	if err != nil {
		return fmt.Errorf("push arg to queue error: %s", err)
	}
	return nil
}

type RetryError struct {
	err error
}

func NewRetryError(err error) *RetryError {
	return &RetryError{err}
}

func (e *RetryError) Error() string {
	return fmt.Sprintf("%s retry", e.err)
}

func (e *RetryError) String() string {
	return e.Error()
}

func (d *DelayDispatcher) Serve() error {
	if len(d.services) == 0 {
		return fmt.Errorf("no service")
	}
	for {
		var sleepTime time.Duration = -1
		for _, s := range d.services {
			s.locker.Lock()
			t, err := s.queue.NextWakeup()
			s.locker.Unlock()
			if err != nil {
				return err
			}
			if sleepTime < 0 || t < sleepTime {
				sleepTime = t
			}
			if t == 0 {
				s.locker.Lock()
				args, err := s.queue.Pop()
				s.locker.Unlock()
				if err != nil {
					d.logger.Crit("Pop fail: %s", err)
					continue
				}
				ret := s.method.method.Call([]reflect.Value{s.service, reflect.ValueOf(args)})
				e := ret[0].Interface()
				if e != nil {
					if _, ok := e.(*RetryError); ok {
						v := reflect.ValueOf(args)
						s.locker.Lock()
						for i := 0; i < v.Len(); i++ {
							data := v.Index(i).Interface().(QueueData)
							err := s.queue.Push(data)
							if err != nil {
								d.logger.Crit("push arg to queue error when retry: %s", err)
							}
						}
						s.locker.Unlock()
					} else {
						d.logger.Crit("Call fail: %s", e)
					}
					continue
				}
			}
		}
		if sleepTime < 0 {
			sleepTime = time.Second
		}
		time.Sleep(sleepTime)
	}
	return nil
}
