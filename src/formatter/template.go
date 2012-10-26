package formatter

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
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
		"append": func(str ...string) string {
			return strings.Join(str, "")
		},
		"column": func(max int, joiner string, content string) string {
			buf := bytes.NewBuffer(nil)
			for len(content) > max {
				buf.WriteString(content[:max])
				buf.WriteString(joiner)
				content = content[max:]
			}
			buf.WriteString(content)
			return buf.String()
		},
		"equal": func(a, b interface{}) bool {
			va := reflect.ValueOf(a)
			vb := reflect.ValueOf(b)
			return va.String() == vb.String()
		},
		"replace": func(from, to, str string) string {
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
		"limit": func(max int, str string) string {
			if max < 2 {
				max = 2
			}
			if len(str) <= max {
				return str
			}
			return str[:max-1] + "…"
		},
		"for": func(array interface{}) (interface{}, error) {
			v := reflect.ValueOf(array)
			if k := v.Kind(); k != reflect.Array && k != reflect.Slice {
				return nil, fmt.Errorf("input must array or slice")
			}
			ret := make([]ForElement, v.Len())
			for i, n := 0, v.Len(); i < n; i++ {
				ret[i].V = v.Index(i).Interface()
				ret[i].Index = i + 1
				ret[i].Last = i == (n - 1)
			}
			return ret, nil
		},
	}
	ret.Funcs(funcs)
	return ret
}

type LocalTemplate struct {
	defaultLang string
	templates   map[string]*template.Template
}

func NewLocalTemplate(path string, defaultLang string) (*LocalTemplate, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open dir %s: %s", path, err)
	}
	infos, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("can't read dir %s: %s", path, err)
	}
	ret := &LocalTemplate{
		defaultLang: defaultLang,
		templates:   make(map[string]*template.Template),
	}
	for _, i := range infos {
		if i.Name()[0] == '.' {
			continue
		}
		template := NewTemplate(i.Name())
		_, err := template.ParseGlob(fmt.Sprintf("%s/%s/*", path, i.Name()))
		if err != nil {
			return nil, fmt.Errorf("can't parse %s/%s: %s", path, i.Name(), err)
		}
		ret.templates[i.Name()] = template
	}
	return ret, nil
}

func (l *LocalTemplate) Execute(wr io.Writer, lang, name string, data interface{}) error {
	t, ok := l.templates[lang]
	if !ok {
		t, ok = l.templates[l.defaultLang]
	}
	if !ok {
		return fmt.Errorf("can't find lang %s or default %s", lang, l.defaultLang)
	}
	err := t.ExecuteTemplate(wr, name, data)
	if err != nil {
		return fmt.Errorf("execute %s error: %s", lang, err)
	}
	return nil
}
