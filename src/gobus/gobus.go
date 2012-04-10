package gobus

import (
	"strings"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/simonz05/godis"
	"reflect"
	"time"
	"unicode"
	"unicode/utf8"
)

//////////////////////////////////////////

const (
	Running int = iota
	Stopped
)

type runner interface {
	run()
}

type baseService struct {
	redis      *godis.Client
	queueName  string
	status     int
	quitChan   chan int
	isQuitChan chan int
	doFunc     reflect.Value
	argType    reflect.Type
	r          runner
}

func (b *baseService) init(netaddr string, db int, password, queueName string, job interface{}, name string, r runner) error {
	b.redis = godis.New(netaddr, db, password)
	b.queueName = fmt.Sprintf("gobus:queue:%s", queueName)
	b.status = Stopped
	b.quitChan = make(chan int)
	b.isQuitChan = make(chan int)
	b.r = r

	v := reflect.ValueOf(job)
	b.doFunc = v.MethodByName(name)
	if b.doFunc == reflect.ValueOf(nil) {
		return fmt.Errorf("Can't find method: %s", name)
	}
	mtype := b.doFunc.Type()
	mname := mtype.Name()
	if mtype.PkgPath() != "" {
		return fmt.Errorf("Method %s must be exported.", mname)
	}
	if mtype.NumIn() < 1 {
		return fmt.Errorf("method", mname, "must has one ins at least.")
	}
	b.argType = mtype.In(0)
	if !isExportedOrBuiltinType(b.argType) {
		return fmt.Errorf(mname, "argument type not exported:", b.argType)
	}

	return nil
}

func (s *baseService) Close() error {
	if s.IsRunning() {
		s.Stop()
	}
	return s.redis.Quit()
}

func (s *baseService) IsRunning() bool {
	return s.status == Running
}

func (s *baseService) Stop() error {
	if !s.IsRunning() {
		return fmt.Errorf("baseService has stopped")
	}

	s.quitChan <- 1
	<-s.isQuitChan
	return nil
}

func (s *baseService) Clear() {
	s.redis.Del(s.queueName)
	s.redis.Del(s.queueName + ":idcount")
	keys, _ := s.redis.Keys(s.queueName + ":*")
	for _, k := range keys {
		s.redis.Del(k)
	}
}

func (s *baseService) Serve(timeOut time.Duration) {
	s.status = Running

Loop:
	for {
		select {
		case <-s.quitChan:
			break Loop
		case <-time.After(timeOut):
			s.r.run()
		}
	}
	s.status = Stopped
	s.isQuitChan <- 1
}

func (s *baseService) getArg() (meta metaType, err error) {
	retBytes, err := s.redis.Lpop(s.queueName)
	if err != nil {
		return
	}

	meta.Arg = reflect.New(s.argType).Interface()

	err = jsonToValue(retBytes, &meta)
	if s.isErr(err, "JSON decode text(%s)", string(retBytes)) {
		return
	}

	return
}

func (s *baseService) isErr(err error, format string, a ...interface{}) bool {
	if err != nil {
		fmt.Printf("gobus server(%s) ", s.queueName)
		fmt.Printf(format, a...)
		fmt.Println(" error:", err)
	}
	return err != nil
}

///////////////////////////////////////////////////

type Service struct {
	baseService
	replyType reflect.Type
}

func CreateService(netaddr string, db int, password, queueName string, job interface{}) (*Service, error) {
	ret := &Service{}
	err := ret.init(netaddr, db, password, queueName, job, "Do", ret)
	if err != nil {
		return nil, err
	}

	mtype := ret.doFunc.Type()
	if mtype.NumIn() != 2 {
		return nil, fmt.Errorf("method Do has wrong number of ins:", mtype.NumIn())
	}
	replyType := mtype.In(1)
	if replyType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("method Do reply type not a pointer:", replyType)
	}
	if !isExportedOrBuiltinType(replyType) {
		return nil, fmt.Errorf("method Do reply type not exported:", replyType)
	}
	ret.replyType = replyType.Elem()
	if mtype.NumOut() != 1 {
		return nil, fmt.Errorf("method Do has wrong number of outs:", mtype.NumOut())
	}
	if returnType := mtype.Out(0); returnType != typeOfError {
		return nil, fmt.Errorf("method Do returns", returnType.String(), "not error")
	}
	return ret, nil
}

func (s *Service) run() {
	for {
		meta, err := s.getArg()
		if err != nil && err.Error() == "Nonexisting key" {
			return
		}
		if s.isErr(err, "Get job from queue fail") {
			return
		}
		s.doJob(meta)
	}
}

func (s *Service) doJob(meta metaType) {
	reply := reflect.New(s.replyType)
	var funcRet []reflect.Value

	defer func() {
		p := recover()
		var ret returnType
		ret.Reply = reply.Interface()
		if p != nil {
			ret.Panic = p
		} else {
			r := funcRet[0].Interface()
			if r != nil {
				ret.Error = r.(error).Error()
			}
		}
		if meta.NeedReply {
			s.sendBack(meta, &ret)
		} else if ret.Panic != nil {
			s.savePanic(meta, ret.Panic)
		} else if ret.Error != "" {
			s.saveFailed(meta, ret.Error)
		}
	}()

	funcRet = s.doFunc.Call([]reflect.Value{reflect.ValueOf(meta.Arg).Elem(), reply})
}

func (s *Service) savePanic(meta metaType, raised interface{}) {
	p := panicType{
		Meta: meta,
		Panic: raised,
	}
	json, _ := valueToJson(p)
	s.redis.Rpush("gobus:fatal", json)
}

func (s *Service) saveFailed(meta metaType, err string) {
	meta.RetryCount++
	if meta.RetryCount > meta.MaxRetry {
		s.savePanic(meta, err)
		return
	}

	failed := failedType{
		Meta: meta,
		Failed: err,
	}
	m := metaType {
		Arg: failed,
		NeedReply: false,
	}
	json, _ := valueToJson(m)
	s.redis.Rpush("gobus:failed", json)
}

func (s *Service) sendBack(meta metaType, ret *returnType) {
	key := meta.Id

	str, err := valueToJson(*ret)
	if s.isErr(err, "JSON Encode value(%v)", ret) {
		return
	}

	err = s.redis.Set(key, str)
	if s.isErr(err, "Redis SET(%s %s)", key, str) {
		return
	}

	_, err = s.redis.Publish(key, 1)
	if s.isErr(err, "Redis PUBLISH(%s, 1)", s.queueName) {
		return
	}
}

//////////////////////////////////////////

type BatchService struct {
	baseService
	argsType reflect.Type
}

func CreateBatchService(netaddr string, db int, password, queueName string, job interface{}) (*BatchService, error) {
	ret := &BatchService{}
	err := ret.init(netaddr, db, password, queueName, job, "Batch", ret)
	if err != nil {
		return nil, err
	}

	ret.argsType = ret.argType
	if ret.argsType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("method Batch arg type not Slice.")
	}
	ret.argType = ret.argsType.Elem()
	return ret, nil
}

func (s *BatchService) run() {
	args := reflect.MakeSlice(s.argsType, 0, 0)
	for {
		meta, err := s.getArg()
		if err != nil && err.Error() == "Nonexisting key" {
			break
		}
		if s.isErr(err, "Get job from queue fail") {
			break
		}
		args = reflect.Append(args, reflect.ValueOf(meta.Arg).Elem())
	}
	if (args.Len() > 0 ) { s.doJobs(args) }
}

func (s *BatchService) doJobs(args reflect.Value) {
	defer func() {
		p := recover()
		if p != nil {
			fmt.Println("Batch Panic:", p)
			panic(p)
		}
	}()
	s.doFunc.Call([]reflect.Value{args})
}

/////////////////////////////////////////

type Client struct {
	netaddr   string
	db        int
	password  string
	redis     *godis.Client
	queueName string
}

func CreateClient(netaddr string, db int, password, queueName string) *Client {
	return &Client{
		netaddr:   netaddr,
		db:        db,
		password:  password,
		queueName: fmt.Sprintf("gobus:queue:%s", queueName),
	}
}

func (c *Client) Do(arg interface{}, reply interface{}) error {
	c.connRedis()
	defer c.closeRedis()

	meta, err := c.makeMeta(arg)
	if err != nil {
		return err
	}
	meta.NeedReply = true

	err = c.send(meta)
	if err != nil {
		return err
	}
	return c.waitReply(meta.Id, reply)
}

func (c *Client) Send(arg interface{}, maxRetry int) error {
	c.connRedis()
	defer c.closeRedis()

	meta, err := c.makeMeta(arg)
	if err != nil {
		return err
	}
	meta.NeedReply = false
	meta.MaxRetry = maxRetry
	return c.send(meta)
}

func (c *Client) connRedis() {
	c.redis = godis.New(c.netaddr, c.db, c.password)
}

func (c *Client) closeRedis() {
	c.redis.Quit()
}

func (c *Client) makeMeta(arg interface{}) (*metaType, error) {
	idCountName := c.getIdCountName()
	idCount, err := c.redis.Incr(idCountName)
	if err != nil {
		return nil, err
	}
	key := c.getId(idCount)

	value := &metaType{
		Id:  key,
		Arg: arg,
		MaxRetry: 5,
	}

	return value, nil
}

func (c *Client) send(m *metaType) error {
	j, err := valueToJson(m)
	if err != nil {
		return err
	}

	_, err = c.redis.Rpush(c.queueName, j)
	if err != nil {
		return err
	}
	_, err = c.redis.Publish(c.queueName, 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) waitReply(id string, reply interface{}) error {
	sub := godis.NewSub(c.netaddr, c.db, c.password)
	sub.Subscribe(id)
	_ = <-sub.Messages
	sub.Close()

	retBytes, err := c.redis.Get(id)
	if err != nil {
		return err
	}

	_, err = c.redis.Del(id)
	if err != nil {
		return err
	}

	ret := &returnType{
		Reply: reply,
	}

	err = jsonToValue(retBytes, ret)
	if err != nil {
		return err
	}

	if ret.Panic != nil {
		panic(ret.Panic)
	}

	if ret.Error != "" {
		return fmt.Errorf(ret.Error)
	}

	return nil
}

func (c *Client) getId(id int64) string {
	return fmt.Sprintf("%s:%d", c.queueName, id)
}

func (c *Client) getIdCountName() string {
	return fmt.Sprintf("%s:idcount", c.queueName)
}

func (c *Client) isErr(err error, format string, a ...interface{}) bool {
	if err != nil {
		fmt.Printf("gobus client(%s) ", c.queueName)
		fmt.Printf(format, a...)
		fmt.Println(" error:", err)
	}
	return err != nil
}

//////////////////////////////////////////

type RetryJob struct {
	netaddr string
	db int
	password string
}

func (j *RetryJob) Batch(args []failedType) {
	for _, arg := range args {
		sp := strings.Split(arg.Meta.Id, ":")
		queue := sp[len(sp) - 2]
		client := CreateClient(j.netaddr, j.db, j.password, queue)
		fmt.Println("Send", arg.Meta, "to", queue, "for reason", arg.Failed)

		client.connRedis()
		defer client.closeRedis()

		err := client.send(&arg.Meta)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GetDefaultRetryServer(netaddr string, db int, password string) (*BatchService, error) {
	job := RetryJob{netaddr, db, password}
	service, err := CreateBatchService(netaddr, db, password, "failed", &job)
	if err != nil {
		return nil, err
	}
	service.queueName = "gobus:failed"
	return service, nil
}

//////////////////////////////////////////

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

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

type metaType struct {
	Id         string
	Arg        interface{}
	RetryCount int
	MaxRetry   int
	NeedReply  bool
}

type returnType struct {
	Error string
	Reply interface{}
	Panic interface{}
}

type panicType struct {
	Meta metaType
	Panic interface{}
}

type failedType struct {
	Meta metaType
	Failed string
}

func jsonToValue(input []byte, value interface{}) error {
	buf := bytes.NewBuffer(input)
	decoder := json.NewDecoder(buf)
	return decoder.Decode(value)
}

func valueToJson(value interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(value)
	return buf.String(), err
}
