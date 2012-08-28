package gobus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/godis"
	"reflect"
	"time"
)

//////////////////////////////////////////

const (
	Running int = iota
	Stopped
)

type Server struct {
	redis      *godis.Client
	queueName  string
	status     int
	quitChan   chan int
	isQuitChan chan int
	domap      map[string]*doMethod
	batchmap   map[string]*batchMethod
}

func CreateServer(netaddr string, db int, password, queueName string) *Server {
	return &Server{
		redis: godis.New(netaddr, db, password),
		queueName: fmt.Sprintf("gobus:queue:%s", queueName),
		status: Stopped,
		quitChan: make(chan int),
		isQuitChan: make(chan int),
		domap: make(map[string]*doMethod),
		batchmap: make(map[string]*batchMethod),
	}
}

func (s *Server) Close() error {
	if s.IsRunning() {
		s.Stop()
	}
	return s.redis.Quit()
}

func (s *Server) IsRunning() bool {
	return s.status == Running
}

func (s *Server) Stop() error {
	if !s.IsRunning() {
		return fmt.Errorf("baseService has stopped")
	}

	s.quitChan <- 1
	<-s.isQuitChan
	return nil
}

func (s *Server) ClearQueue() {
	s.redis.Del(s.queueName)
	s.redis.Del(s.queueName + ":idcount")
	keys, _ := s.redis.Keys(s.queueName + ":*")
	for _, k := range keys {
		s.redis.Del(k)
	}
}

func (s *Server) Register(service interface{}) error {
	domap, batchmap := getMethods(service)
	for k, v := range domap {
		if _, ok := s.domap[k]; ok {
			return fmt.Errorf("Service has registered: %s(arg, reply)", k)
		}
		s.domap[k] = v
	}
	for k, v := range batchmap {
		if _, ok := s.batchmap[k]; ok {
			return fmt.Errorf("Service has registered: %s(args)", k)
		}
		s.batchmap[k] = v
	}
	return nil
}

func (s *Server) Serve(timeOut time.Duration) {
	s.status = Running

Loop:
	for {
		select {
		case <-s.quitChan:
			break Loop
		case <-time.After(timeOut):
			s.run()
		}
	}
	s.status = Stopped
	s.isQuitChan <- 1
}

func (s *Server) savePanic(meta metaType, raised interface{}) {
	p := panicType{
		Meta: meta,
		Panic: raised,
	}
	json, _ := valueToJson(p)
	s.redis.Rpush("gobus:fatal", json)
}

func (s *Server) saveFailed(meta metaType, err string) {
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
		Method: "Resend",
		Arg: failed,
		NeedReply: false,
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.Encode(m.Method)
	encoder.Encode(m)
	s.redis.Rpush("gobus:failed", buf.String())
}

func (s *Server) popQueue() (meta metaType, err error) {
	retBytes, err := s.redis.Lpop(s.queueName)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(retBytes)
	decoder := json.NewDecoder(buf)
	err = decoder.Decode(&meta.Method)
	if s.isErr(err, "JSON decode text(%s)", string(retBytes)) {
		return
	}

	if f, ok := s.domap[meta.Method]; ok {
		meta.Arg = reflect.New(f.arg).Interface()
	} else if f, ok := s.batchmap[meta.Method]; ok {
		meta.Arg = reflect.New(f.arg).Interface()
	} else {
		err = fmt.Errorf("Can't find method: %s", meta.Method)
		return
	}

	err = decoder.Decode(&meta)
	if s.isErr(err, "JSON decode text(%s)", string(retBytes)) {
		return
	}

	return
}

func (s *Server) isErr(err error, format string, a ...interface{}) bool {
	if err != nil {
		fmt.Printf("gobus server(%s) ", s.queueName)
		fmt.Printf(format, a...)
		fmt.Println(" error:", err)
	}
	return err != nil
}

func (s *Server) do(f *doMethod, meta metaType) {
	reply, p, err := f.call(meta)

	ret := returnType{
		Panic: p,
		Error: StringError(err),
		Reply: reply,
	}

	if meta.NeedReply {
		s.sendBack(meta, &ret)
	} else if ret.Panic != nil {
		s.savePanic(meta, ret.Panic)
	} else if ret.Error != "" {
		s.saveFailed(meta, ret.Error)
	}
}

func (s *Server) batch(f *batchMethod, args reflect.Value, metas []metaType) {
	p, err := f.call(args)

	if p == nil && err == nil { return }

	for _, meta := range metas {
		meta.RetryCount++
		if p != nil {
			s.savePanic(meta, p)
		} else {
			s.saveFailed(meta, err.Error())
		}
	}
}

func (s *Server) run() {
	meta, err := s.popQueue()
	for {
		if err != nil && err.Error() == "Nonexisting key" {
			break
		}
		if s.isErr(err, "Get job from queue fail") {
			break
		}

		if f, ok := s.domap[meta.Method]; ok {
			s.do(f, meta)
			meta, err = s.popQueue()
		} else if f, ok := s.batchmap[meta.Method]; ok && !meta.NeedReply {
			var metas []metaType
			var args reflect.Value
			metas, args, meta, err = s.grabSameMethod(f, meta)
			if (args.Len() > 0 ) { s.batch(f, args, metas) }
		} else {
			if meta.NeedReply {
				ret := returnType{
					Error: fmt.Sprintf("Can't find service: %s(arg, reply)", meta.Method),
				}
				s.sendBack(meta, &ret)
			} else {
				s.savePanic(meta, fmt.Sprintf("Can't find service: %s", meta.Method))
			}
			meta, err = s.popQueue()
		}
	}
}

func (s *Server) grabSameMethod(f *batchMethod, meta metaType) ([]metaType, reflect.Value, metaType, error) {
	metas := []metaType{meta}
	args := reflect.MakeSlice(f.argSlice, 0, 0)
	args = reflect.Append(args, reflect.ValueOf(meta.Arg).Elem())
	method := meta.Method
	var err error

	For: for {
		meta, err = s.popQueue()
		if err != nil && err.Error() == "Nonexisting key" {
			break For
		}
		if method != meta.Method || meta.NeedReply {
			break For
		}
		if s.isErr(err, "Get job from queue fail") {
			break For
		}
		args = reflect.Append(args, reflect.ValueOf(meta.Arg).Elem())
		metas = append(metas, meta)
	}
	return metas, args, meta, err
}

func (s *Server) sendBack(meta metaType, ret *returnType) {
	key := meta.Id

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	if ret.Panic != nil {
		encoder.Encode("Panic")
		encoder.Encode(ret.Panic)
	} else if ret.Error != "" {
		encoder.Encode("Error")
		encoder.Encode(ret.Error)
	} else {
		encoder.Encode("OK")
		err := encoder.Encode(ret.Reply)
		if s.isErr(err, "JSON Encode value(%v)", ret) {
			return
		}
	}

	err := s.redis.Set(key, buf.String())
	if s.isErr(err, "Redis SET(%s %s)", key, buf.String()) {
		return
	}

	_, err = s.redis.Publish(key, 1)
	if s.isErr(err, "Redis PUBLISH(%s, 1)", s.queueName) {
		return
	}
}

type metaType struct {
	Id         string
	Method     string
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
