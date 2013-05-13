package gobus

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type Route struct {
	*mux.Route
	server *Server
}

func (r *Route) Queries(params ...string) *Route {
	r.Route.Queries(params...)
	return r
}

func (r *Route) Methods(methods ...string) *Route {
	r.Route.Methods(methods...)
	return r
}

func (r *Route) Path(tpl string) *Route {
	r.Route.Path(tpl)
	return r
}

func (r *Route) HandlerMethod(codec Codec, service interface{}, method string) error {
	t := reflect.TypeOf(service)
	v := reflect.ValueOf(service)
	m, ok := t.MethodByName(method)
	if !ok {
		return fmt.Errorf("can't find method")
	}
	if m.Type.NumOut() != 2 {
		return fmt.Errorf("output service is not 2")
	}
	if m.Type.Out(1) != typeOfError {
		return fmt.Errorf("second output is not error")
	}
	in := m.Type.NumIn()
	if in != 2 && in != 3 {
		return fmt.Errorf("method must have 1 or 2 input arguments")
	}
	if m.Type.In(1) != typeOfMap {
		return fmt.Errorf("first input is not map[string]string")
	}

	inputPtr := false
	var inputType reflect.Type
	if in == 3 {
		inputType = m.Type.In(2)
		if inputType.Kind() == reflect.Ptr {
			inputType = inputType.Elem()
			inputPtr = true
		}
	}

	r.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var err error
		var input reflect.Value
		p := params(req)
		defer func() {
			pa := recover()
			if pa != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("%+v", pa)))
				return
			}
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		}()

		args := []reflect.Value{v, reflect.ValueOf(p)}
		if in == 3 {
			input, err = codec.Decode(req.Body, inputType)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !inputPtr {
				input = input.Elem()
			}
			args = append(args, input)
		}

		rets := m.Func.Call(args)
		ret, e := rets[0], rets[1].Interface()
		if e != nil {
			err = e.(error)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = codec.Encode(w, ret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", codec.Mime()))
		return
	})
	return nil
}

func params(r *http.Request) map[string]string {
	vars := mux.Vars(r)
	q := r.URL.Query()
	for k, _ := range q {
		vars[k] = q.Get(k)
	}
	return vars
}
