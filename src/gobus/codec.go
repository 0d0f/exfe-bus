package gobus

import (
	"encoding/json"
	"io"
	"reflect"
)

type JSON struct{}

func (j *JSON) Mime() string {
	return "application/json"
}

func (j *JSON) Decode(r io.Reader, t reflect.Type) (reflect.Value, error) {
	decoder := json.NewDecoder(r)
	ret := reflect.New(t)
	err := decoder.Decode(ret.Interface())
	return ret, err
}

func (j *JSON) Encode(w io.Writer, v reflect.Value) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(v.Interface())
}
