package formatter

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"
	"unicode/utf8"
)

type ForElement struct {
	V     interface{}
	Index int
	First bool
	Last  bool
}

func NewTemplate(name string) *template.Template {
	ret := template.New(name)

	funcs := template.FuncMap{
		"substr": func(start, l int, str string) string {
			if start+l > len(str) {
				l = len(str) - start
			}
			return str[start : start+l]
		},
		"sub": func(data interface{}, templates ...string) (string, error) {
			buf := bytes.NewBuffer(nil)
			var t *template.Template
			for _, templ := range templates {
				t = ret.Lookup(templ)
				if t != nil {
					break
				}
			}
			if t == nil {
				return "", fmt.Errorf("can't find both %v templates.", templates)
			}
			err := t.Execute(buf, data)
			return buf.String(), err
		},
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
		"limit": func(max int, str string) string {
			if max < 4 {
				max = 4
			}
			if utf8.RuneCountInString(str) <= max {
				return str
			}
			max = max - 3
			ret := bytes.NewBuffer(nil)
			for _, b := range []byte(str) {
				ret.WriteByte(b)
				bs := ret.Bytes()
				if utf8.Valid(bs) && utf8.RuneCount(bs) == max {
					break
				}
			}
			return ret.String() + "..."
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
				ret[i].First = i == 0
				ret[i].Last = i == (n - 1)
			}
			return ret, nil
		},
		"plural": func(single, multi string, length int) string {
			if length <= 1 {
				return single
			}
			return multi
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
		if i.Name() == "image_data" {
			continue
		}
		template := NewTemplate(i.Name())
		err := parseDirTemplate(template, fmt.Sprintf("%s/%s", path, i.Name()), "")
		if err != nil {
			return nil, fmt.Errorf("can't parse %s/%s: %s", path, i.Name(), err)
		}
		ret.templates[i.Name()] = template
	}
	return ret, nil
}

func (l *LocalTemplate) IsExist(lang, name string) bool {
	t, ok := l.templates[lang]
	if !ok {
		t, ok = l.templates[l.defaultLang]
		if !ok {
			return false
		}
	}
	return t.Lookup(name) != nil
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

func parseDirTemplate(t *template.Template, dir, name string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	fis, err := f.Readdir(0)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		n := fi.Name()
		if name != "" {
			n = fmt.Sprintf("%s/%s", name, fi.Name())
		}
		path := fmt.Sprintf("%s%c%s", dir, os.PathSeparator, fi.Name())
		if fi.IsDir() {
			err := parseDirTemplate(t, path, n)
			if err != nil {
				return err
			}
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		_, err = t.New(n).Parse(string(content))
		if err != nil {
			return err
		}
	}
	return nil
}
