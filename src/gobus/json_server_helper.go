package gobus

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"unicode"
	"unicode/utf8"
)

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfHTTPMeta = reflect.TypeOf((*HTTPMeta)(nil)).Elem()

func methodOutIsError(mtype *reflect.Type) bool {
	if (*mtype).NumOut() != 1 {
		return false
	}
	if e := (*mtype).Out(0); e != typeOfError {
		return false
	}
	return true
}

type methodType struct {
	method reflect.Value
	arg    reflect.Type
	reply  reflect.Type
}

func newMethodType(service reflect.Type) map[string]*methodType {
	ret := make(map[string]*methodType)

	for i, n := 0, service.NumMethod(); i < n; i++ {
		m := service.Method(i)

		if m.PkgPath != "" {
			// Method must be exported.
			continue
		}

		if !methodOutIsError(&m.Type) {
			continue
		}

		if m.Type.NumIn() != 4 {
			continue
		}

		reqType := m.Type.In(1)
		if reqType.Kind() != reflect.Ptr {
			continue
		}
		if reqType.Elem() != typeOfHTTPMeta {
			continue
		}
		argType := m.Type.In(2)
		if !isExportedOrBuiltinType(argType) {
			continue
		}
		replyType := m.Type.In(3)
		if !isExportedOrBuiltinType(replyType) {
			continue
		}
		if replyType.Kind() != reflect.Ptr {
			continue
		}
		ret[m.Name] = &methodType{
			method: m.Func,
			arg:    argType,
			reply:  replyType.Elem(),
		}
	}

	return ret
}

func (m *methodType) getInput(r *http.Request) (ret reflect.Value, err error) {
	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		err = e
		return
	}
	if m.arg.Kind() == reflect.Ptr {
		ret = reflect.New(m.arg.Elem())
	} else {
		ret = reflect.New(m.arg)
	}
	err = json.Unmarshal(body, ret.Interface())
	if m.arg.Kind() != reflect.Ptr {
		ret = ret.Elem()
	}
	return
}

func (m *methodType) call(service reflect.Value, meta *HTTPMeta, arg reflect.Value) (reflect.Value, error) {
	reply := reflect.New(m.reply)
	rets := m.method.Call([]reflect.Value{service, reflect.ValueOf(meta), arg, reply})
	ret := rets[0].Interface()
	if ret != nil {
		return reply, ret.(error)
	}
	return reply, nil
}

type serviceType struct {
	service reflect.Value
	methods map[string]*methodType
}

func newServiceType(service interface{}) *serviceType {
	v := reflect.ValueOf(service)
	return &serviceType{
		service: v,
		methods: newMethodType(v.Type()),
	}
}
