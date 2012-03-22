package gosque

import (
	"strings"
	"unicode"
	"unicode/utf8"
	"bytes"
	"reflect"
	"encoding/json"
	"fmt"
	"github.com/simonz05/godis"
	"time"
)

type Client struct {
	redis       *godis.Client
	queueFilter string
	registered  map[string]*registeredFunc
}

func CreateQueue(netaddr string, db int, password, queueFilter string) *Client {
	return &Client{
		redis:       godis.New(netaddr, db, password),
		queueFilter: queueFilter,
		registered:  make(map[string]*registeredFunc),
	}
}

func (c *Client) Close() {
	c.redis.Quit()
}

func (c *Client) Register(job interface{}) error {
	v := reflect.ValueOf(job)
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	typename := strings.ToLower(t.Name())

	f := v.MethodByName("Perform")
	if f == reflect.ValueOf(nil) {
		return fmt.Errorf("Can't find method Perform")
	}
	mtype := f.Type()
	if mtype.NumIn() < 1 {
		return fmt.Errorf("method Perform must has one ins")
	}
	argtype := mtype.In(0)
	if !argtype.Implements(reflect.TypeOf(new(argType)).Elem()) {
		return fmt.Errorf("Perform argument type not exported:", argtype)
	}
	if mtype.NumOut() != 0 {
		return fmt.Errorf("method Do has wrong number of outs:", mtype.NumOut())
	}

	c.registered[typename] = &registeredFunc{
		arg: mtype.In(0),
		function: f,
	}

	return nil
}

func (c *Client) Enqueue(name string, arg argType) error {
	data := metaType{}
	data.Class = strings.ToLower(name)
	data.Args = arg

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.Encode(data)

	q := fmt.Sprintf("resque:queue:%s",arg.Queue())
	_, err := c.redis.Rpush(q, buf.String())
	return err
}

func (c *Client) Serve(timeOut time.Duration) {
	queues := fmt.Sprintf("resque:queue:%s", c.queueFilter)
	for {
		elem, err := c.redis.Lpop(queues)
		if err != nil {
			time.Sleep(timeOut)
			continue
		}

		buffer := bytes.NewBuffer(elem)
		decoder := json.NewDecoder(buffer)

		type meta struct {
			Class string
		}
		value := meta{}
		err = decoder.Decode(&value)
		if err != nil {
			fmt.Printf("Error parse value(%s) to meta: %s\n", string(elem), err)
			continue
		}

		f := c.registered[value.Class]
		if f == nil {
			fmt.Printf("Can't find job %s registered\n", value.Class)
			continue
		}
		data := metaType{
			Args: reflect.New(f.arg).Interface(),
		}
		buffer = bytes.NewBuffer(elem)
		decoder = json.NewDecoder(buffer)
		err = decoder.Decode(&data)
		if err != nil {
			fmt.Printf("Error parse value(%s) to arg: %s\n", string(elem), err)
			continue
		}

		f.function.Call([]reflect.Value{reflect.ValueOf(data.Args).Elem()})
	}
}

type metaType struct {
	Class string
	Args interface{}
	Id   string
}

type argType interface {
	Queue() string
}

type registeredFunc struct {
	arg reflect.Type
	function reflect.Value
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func firstUpper(str string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(str[0:1]), str[1:])
}
