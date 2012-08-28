package gobus

import (
	"bytes"
	"encoding/json"
	"reflect"
	"unicode"
	"unicode/utf8"
)

type doMethod struct {
	method reflect.Value
	arg reflect.Type
	reply reflect.Type
}

func (m *doMethod)call(meta metaType) (reply interface{}, p interface{}, err error) {
	r := reflect.New(m.reply)
	var funcRet []reflect.Value

	defer func() {
		p = recover()
		if p == nil {
			err = getErrorFromReturn(funcRet)
		}
		reply = r.Interface()
	}()

	funcRet = m.method.Call([]reflect.Value{reflect.ValueOf(meta.Arg).Elem(), r})
	return
}

type batchMethod struct {
	method reflect.Value
	arg reflect.Type
	argSlice reflect.Type
}

func (m *batchMethod) call(args reflect.Value) (p interface{}, err error) {
	var funcRet []reflect.Value
	defer func() {
		p = recover()
		if p == nil {
			err = getErrorFromReturn(funcRet)
		}
	}()

	funcRet = m.method.Call([]reflect.Value{args})
	return
}


func getMethods(service interface{}) (domap map[string]*doMethod, batchmap map[string]*batchMethod) {
	domap, batchmap = make(map[string]*doMethod), make(map[string]*batchMethod)

	t := reflect.TypeOf(service)
	v := reflect.ValueOf(service)
	for i, n := 0, t.NumMethod(); i<n; i++ {
		m := t.Method(i)

		if m.PkgPath != "" {
			// Method must be exported.
			continue
		}

		if ! methodOutIsError(&m.Type) {
			continue
		}

		switch m.Type.NumIn() {
		case 2:
			argsType := m.Type.In(1)
			if !isExportedOrBuiltinType(argsType) {
				continue
			}
			if argsType.Kind() != reflect.Slice {
				continue
			}
			batchmap[m.Name] = &batchMethod{
				method: v.Method(i),
				arg: argsType.Elem(),
				argSlice: argsType,
			}
		case 3:
			argType := m.Type.In(1)
			if !isExportedOrBuiltinType(argType) {
				continue
			}
			replyType := m.Type.In(2)
			if !isExportedOrBuiltinType(replyType) {
				continue
			}
			if replyType.Kind() != reflect.Ptr {
				continue
			}
			domap[m.Name] = &doMethod{
				method: v.Method(i),
				arg: argType,
				reply: replyType.Elem(),
			}
		default:
			continue
		}
	}
	return
}

func methodOutIsError(mtype *reflect.Type) bool {
	if (*mtype).NumOut() != 1 {
		return false
	}
	if e := (*mtype).Out(0); e != typeOfError {
		return false
	}
	return true
}

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

func valueToJson(value interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(value)
	return buf.String(), err
}

func getErrorFromReturn(ret []reflect.Value) error {
	r := ret[0].Interface()
	if r == nil {
		return nil
	}
	return r.(error)
}

func StringError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
