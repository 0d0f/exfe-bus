package gobus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/simonz05/godis"
	"reflect"
	"time"
)

//////////////////////////////////////////

const (
	Running int = iota
	Stopped
)

type BaseService struct {
	redis      *godis.Client
	queueName  string
	status     int
	quitChan   chan int
	isQuitChan chan int
	doFunc     reflect.Value
	argType    reflect.Type
}

func initBaseService(netaddr string, db int, password, queueName string) *BaseService {
	return &BaseService{
		redis:      godis.New(netaddr, db, password),
		queueName:  fmt.Sprintf("gobus:queue:%s", queueName),
		status:     Stopped,
		quitChan:   make(chan int),
		isQuitChan: make(chan int),
	}
}

func (s *BaseService) Close() error {
	if s.IsRunning() {
		s.Stop()
	}
	return s.redis.Quit()
}

func (s *BaseService) IsRunning() bool {
	return s.status == Running
}

func (s *BaseService) Stop() error {
	if !s.IsRunning() {
		return fmt.Errorf("BaseService has stopped")
	}

	s.quitChan <- 1
	<-s.isQuitChan
	return nil
}

func (s *BaseService) Clear() {
	s.redis.Del(s.queueName)
	s.redis.Del(s.queueName + ":idcount")
	keys, _ := s.redis.Keys(s.queueName + ":*")
	for _, k := range keys {
		s.redis.Del(k)
	}
}

func (s *BaseService) GetArg() (meta metaType, err error) {
	retBytes, err := s.redis.Lpop(s.queueName)
	if err != nil {
		return
	}

	meta.Arg = reflect.New(s.argType).Interface()

	err = jsonToValue(retBytes, &meta)
	if s.IsErr(err, "JSON decode text(%s)", string(retBytes)) {
		return
	}

	return
}

func (s *BaseService) IsErr(err error, format string, a ...interface{}) bool {
	if err != nil {
		fmt.Printf("gobus server(%s) ", s.queueName)
		fmt.Printf(format, a...)
		fmt.Println(" error:", err)
	}
	return err != nil
}

type Service struct {
	BaseService
	replyType  reflect.Type
}

func CreateService(netaddr string, db int, password, queueName string, job interface{}) *Service {
	ret := &Service{
		BaseService: initBaseService(netaddr, db, password, queueName),
	}

	v := reflect.ValueOf(job)
	ret.doFunc := v.MethodByName("Do")
	doType := ret.doFunc.Type()
	ret.argType := doType.In(0)
	ret.replyType := doType.In(1)
	return ret
}

func (s *Service) Run(timeOut time.Duration) {
	s.status = Running

Loop:
	for {
		select {
		case <-s.quitChan:
			break Loop
		case <-time.After(timeOut):
			s.handleQueue()
		}
	}
	s.status = Stopped
	s.isQuitChan <- 1
}

func (s *Service) handleQueue() {
	meta, err := s.GetArg()
	if err != nil && err.Error() == "Nonexisting key" {
		return
	}
	if s.isErr(err, "Get job from queue fail") {
		return
	}
	s.doJob(meta)
}

func (s *Service) doJob(meta metaType) {
	reply := reflect.New(s.replyType)
	ret := s.doFunc.Call([]reflect.Value{reflect.ValueOf(meta.Arg).Elem(), reply})

	if meta.NeedReply {
		r := ret[0].Interface()
		var err error
		if r != nil {
			err = r.(error)
		}
		s.sendBack(err, meta, reply.Interface())
	}
}

func (s *Service) sendBack(err error, meta metaType, reply interface{}) {
	key := meta.Id

	ret := returnType{
		Reply: reply,
	}
	if err != nil {
		ret.Error = err.Error()
	}

	str, err := valueToJson(ret)
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
	redis      *godis.Client
	queueName  string
	status     int
	quitChan   chan int
	isQuitChan chan int
	doFunc     reflect.Value
	argType    reflect.Type
	argsType   reflect.Type
}

func CreateBatchService(netaddr string, db int, password, queueName string, job interface{}) *BatchService {
	v := reflect.ValueOf(job)
	doFunc := v.MethodByName("Batch")
	doType := doFunc.Type()
	return &BatchService{
		redis:      godis.New(netaddr, db, password),
		queueName:  fmt.Sprintf("gobus:queue:%s", queueName),
		status:     Stopped,
		quitChan:   make(chan int),
		isQuitChan: make(chan int),
		doFunc:     doFunc,
		argType:    doType.In(0).Elem(),
		argsType:   doType.In(0),
	}
}

func (s *BatchService) Close() error {
	if s.IsRunning() {
		s.Stop()
	}
	return s.redis.Quit()
}

func (s *BatchService) Run(timeOut time.Duration) {
	s.status = Running

Loop:
	for {
		select {
		case <-s.quitChan:
			break Loop
		case <-time.After(timeOut):
			s.handleQueue()
		}
	}
	s.status = Stopped
	s.isQuitChan <- 1
}

func (s *BatchService) IsRunning() bool {
	return s.status == Running
}

func (s *BatchService) Stop() error {
	if !s.IsRunning() {
		return fmt.Errorf("BatchService has stopped")
	}

	s.quitChan <- 1
	<-s.isQuitChan
	return nil
}

func (s *BatchService) Clear() {
	s.redis.Del(s.queueName)
	s.redis.Del(s.queueName + ":idcount")
	keys, _ := s.redis.Keys(s.queueName + ":*")
	for _, k := range keys {
		s.redis.Del(k)
	}
}

func (s *BatchService) handleQueue() {
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
	s.doJobs(args)
}

func (s *BatchService) getArg() (meta metaType, err error) {
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

func (s *BatchService) doJobs(args reflect.Value) {
	s.doFunc.Call([]reflect.Value{args})
}

func (s *BatchService) isErr(err error, format string, a ...interface{}) bool {
	if err != nil {
		fmt.Printf("gobus server(%s) ", s.queueName)
		fmt.Printf(format, a...)
		fmt.Println(" error:", err)
	}
	return err != nil
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
	redis := godis.New(netaddr, db, password)

	return &Client{
		netaddr:   netaddr,
		db:        db,
		password:  password,
		redis:     redis,
		queueName: fmt.Sprintf("gobus:queue:%s", queueName),
	}
}

func (c *Client) Close() error {
	return c.redis.Quit()
}

func (c *Client) Do(arg interface{}, reply interface{}) error {
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

func (c *Client) Send(arg interface{}) error {
	meta, err := c.makeMeta(arg)
	if err != nil {
		return err
	}
	meta.NeedReply = false
	return c.send(meta)
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
		Error: "",
		Reply: reply,
	}

	err = jsonToValue(retBytes, ret)
	if err != nil {
		return err
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

type metaType struct {
	Id        string
	Arg       interface{}
	NeedReply bool
}

type returnType struct {
	Error string
	Reply interface{}
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
