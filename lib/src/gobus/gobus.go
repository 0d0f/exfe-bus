package gobus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/simonz05/godis"
	"time"
)

type Worker interface {
	Do(jobs []interface{}) []interface{}
	MaxJobsCount() int
	JobGenerator() interface{}
}

//////////////////////////////////////////

const (
	Running int = iota
	Stopped
)

type Service struct {
	redis      *godis.Client
	queueName  string
	worker     Worker
	status     int
	quitChan   chan int
	isQuitChan chan int
}

func CreateService(netaddr string, db int, password, queueName string, worker Worker) *Service {
	return &Service{
		redis:      godis.New(netaddr, db, password),
		queueName:  queueName,
		worker:     worker,
		status:     Stopped,
		quitChan:   make(chan int),
		isQuitChan: make(chan int),
	}
}

func (s *Service) Close() error {
	if s.IsRunning() {
		s.Stop()
	}
	return s.redis.Quit()
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

func (s *Service) IsRunning() bool {
	return s.status == Running
}

func (s *Service) Stop() error {
	if !s.IsRunning() {
		return fmt.Errorf("Service has stopped")
	}

	s.quitChan <- 1
	<-s.isQuitChan
	return nil
}

func (s *Service) Clear() {
	s.redis.Del(s.queueName)
	keys, _ := s.redis.Keys(s.queueName + ":*")
	for _, k := range keys {
		s.redis.Del(k)
	}
}

func (s *Service) handleQueue() {
	for {
		jobsCount, err := s.redis.Llen(s.queueName)
		if s.isErr(err, "Redis LLEN(%s)", s.queueName) || jobsCount == 0 {
			break
		}

		jobs, meta := s.getJobs(s.worker.MaxJobsCount())
		if len(meta) == 0 {
			fmt.Println("Can't get jobs, unknown error")
			continue
		}
		rets := s.worker.Do(jobs)
		if rets != nil {
			for i, ret := range rets {
				s.sendBack(ret, meta[i])
			}
		}
	}
}

func (s *Service) getJobs(max int) (jobs []interface{}, metas []metaType) {
	for i := s.worker.MaxJobsCount(); i > 0; i-- {
		retBytes, err := s.redis.Lpop(s.queueName)
		if err != nil {
			return
		}

		var meta metaType
		meta.Data = s.worker.JobGenerator()
		err = jsonToValue(retBytes, &meta)
		if s.isErr(err, "JSON decode text(%s)", string(retBytes)) {
			continue
		}

		jobs, metas = append(jobs, meta.Data), append(metas, meta)
	}
	return
}

func (s *Service) sendBack(ret interface{}, meta metaType) {
	key := meta.ResponseKey
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

func (s *Service) isErr(err error, format string, a ...interface{}) bool {
	if err != nil {
		fmt.Printf("gobus server(%s) ", s.queueName)
		fmt.Printf(format, a)
		fmt.Println(" error:", err)
	}
	return err != nil
}

/////////////////////////////////////////

type Client struct {
	netaddr        string
	db             int
	password       string
	redis          *godis.Client
	queueName      string
	valueGenerator valueGeneratorFunc
}

func CreateClient(netaddr string, db int, password, queueName string, valueGenerator valueGeneratorFunc) *Client {
	redis := godis.New(netaddr, db, password)

	return &Client{
		netaddr:        netaddr,
		db:             db,
		password:       password,
		redis:          redis,
		queueName:      queueName,
		valueGenerator: valueGenerator,
	}
}

func (c *Client) Close() error {
	return c.redis.Quit()
}

func (c *Client) Send(v interface{}) (string, error) {
	idCountName := c.getIdCountName()
	idCount, err := c.redis.Incr(idCountName)
	if err != nil {
		return "", err
	}
	key := c.getResponseKey(idCount)

	value := metaType{
		ResponseKey: key,
		Data:        v,
	}

	j, err := valueToJson(value)
	if err != nil {
		return "", err
	}

	_, err = c.redis.Rpush(c.queueName, j)
	if err != nil {
		return "", err
	}
	_, err = c.redis.Publish(c.queueName, 0)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (c *Client) WaitResponse(key string) (interface{}, error) {
	sub := godis.NewSub(c.netaddr, c.db, c.password)
	sub.Subscribe(key)
	_ = <-sub.Messages
	sub.Close()

	retBytes, err := c.redis.Get(key)
	if err != nil {
		return nil, err
	}

	_, err = c.redis.Del(key)
	if err != nil {
		return nil, err
	}

	ret := c.valueGenerator()
	err = jsonToValue(retBytes, &ret)
	return ret, err
}

func (c *Client) Do(v interface{}) (interface{}, error) {
	key, err := c.Send(v)
	if err != nil {
		return nil, err
	}
	return c.WaitResponse(key)
}

func (c *Client) getResponseKey(id int64) string {
	return fmt.Sprintf("%s:%d", c.queueName, id)
}

func (c *Client) getIdCountName() string {
	return fmt.Sprintf("%s:idcount", c.queueName)
}

func (c *Client) IsErr(err error, format string, a ...interface{}) bool {
	if err != nil {
		fmt.Printf("gobus client(%s) ", c.queueName)
		fmt.Printf(format, a)
		fmt.Println(" error:", err)
	}
	return err != nil
}

//////////////////////////////////////////

type valueGeneratorFunc func() interface{}

type metaType struct {
	ResponseKey string
	Data        interface{}
}

func jsonToValue(input []byte, value interface{}) error {
	buf := bytes.NewBuffer(input)
	decoder := json.NewDecoder(buf)
	return decoder.Decode(&value)
}

func valueToJson(value interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(value)
	return buf.String(), err
}
