package broker

import (
	"fmt"
	"io"
	"io/ioutil"
)

type PlainText struct{}

func (p PlainText) Unmarshal(r io.Reader, v interface{}) error {
	ps, ok := v.(*string)
	if !ok {
		return fmt.Errorf("plain text only can save in string")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	*ps = string(b)
	return nil
}

func (p PlainText) Marshal(w io.Writer, name string, v interface{}) error {
	_, err := w.Write([]byte(fmt.Sprintf("%s", v)))
	return err
}

type TextError string

func (t TextError) Error() string {
	return string(t)
}

func (p PlainText) Error(code int, message string) error {
	return TextError(fmt.Sprintf("(%d)%s", code, message))
}
