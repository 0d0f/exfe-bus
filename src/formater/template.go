package formater

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

type ForElement struct {
	V     interface{}
	Index int
	Last  bool
}

func NewTemplate(name string) *template.Template {
	ret := template.New(name)

	funcs := template.FuncMap{
		"replace": func(str, from, to string) string {
			return strings.Replace(str, from, to, -1)
		},
		"base64": func(str string) string {
			return base64.StdEncoding.EncodeToString([]byte(str))
		},
		"base64url": func(str string) string {
			return base64.URLEncoding.EncodeToString([]byte(str))
		},
		"sub": func(name string, data interface{}) (string, error) {
			buf := bytes.NewBuffer(nil)
			err := ret.ExecuteTemplate(buf, name, data)
			return buf.String(), err
		},
		"for": func(array interface{}) (interface{}, error) {
			v := reflect.ValueOf(array)
			if k := v.Kind(); k != reflect.Array && k != reflect.Slice {
				return nil, fmt.Errorf("input must array or slice")
			}
			ret := make([]ForElement, v.Len())
			for i, n := 0, v.Len(); i < n; i++ {
				ret[i].V = v.Index(i)
				ret[i].Index = i + 1
				ret[i].Last = i == (n - 1)
			}
			return ret, nil
		},
	}
	ret.Funcs(funcs)
	return ret
}
