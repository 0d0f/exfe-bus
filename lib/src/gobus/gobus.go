package gobus

import (
	"github.com/simonz05/godis"
	"time"
	"fmt"
	"bytes"
	"encoding/json"
)

type Worker interface {
	Do(jobs []interface{}) []interface{}
	MaxJobsCount() int
	JobGenerator() interface{}
}

type innerValue struct {
	ResponseKey string
	Data interface{}
}

//////////////////////////////////////////

const (
	Running int = iota
	Stopped
)

type Service struct {
	redis *godis.Client
	queueName string
	worker Worker
	status int
	quitChan chan int
	isQuitChan chan int
}

func CreateService(netaddr string, db int, password, queueName string, worker Worker) *Service {
	return &Service{
		redis: godis.New(netaddr, db, password),
		queueName: queueName,
		worker: worker,
		status: Stopped,
		quitChan: make(chan int),
		isQuitChan: make(chan int),
	}
}

func (s *Service) Close() error {
	if s.IsRunning() { s.Stop() }
	return s.redis.Quit()
}

func (s *Service) Run(timeOut time.Duration) error {
	s.status = Running
	go func() {
Loop:
		for {
			select {
			case <-s.quitChan:
				break Loop;
			case <-time.After(timeOut):
				for {
					jobsCount, err := s.redis.Llen(s.queueName)
					if err != nil {
						fmt.Printf("Redis LLEN(%s) error: %s\n", s.queueName, err)
						continue
					}
					if jobsCount == 0 {
						break
					}

					jobs, meta := s.getJobs(s.worker.MaxJobsCount())
					if len(meta) == 0 {
						fmt.Println("Can't get jobs, unknown error")
						continue
					}
					rets := s.worker.Do(jobs)
					for i, ret := range rets {
						s.sendBack(ret, meta[i])
					}
				}
			}
		}
		s.status = Stopped
		s.isQuitChan <- 1
	}()
	return nil
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

func (s *Service) Empty() {
	idCountName := fmt.Sprintf("%s:idcount", s.queueName)
	s.redis.Del(idCountName)
	s.redis.Del(s.queueName)
}

func jsonToValue(str []byte, value interface{}) error {
	buf := bytes.NewBuffer(str)
	decoder := json.NewDecoder(buf)
	return decoder.Decode(&value)
}

func valueToJson(value interface{}) (string, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(value)
	return buf.String(), err
}

func (s *Service) getJobs(max int) (jobs []interface{}, meta []innerValue) {
	for i:=s.worker.MaxJobsCount(); ; {
		if i > 0 { i-- }
		retBytes, err := s.redis.Lpop(s.queueName)
		if err != nil {
			return
		}

		var value innerValue
		value.Data = s.worker.JobGenerator()
		err = jsonToValue(retBytes, &value)
		if err != nil {
			fmt.Printf("JSON decode error: %s, text: %s\n", err, string(retBytes))
			continue
		}

		jobs = append(jobs, value.Data)
		meta = append(meta, value)
		if i == 0 { return }
	}
	return
}

func (s *Service) sendBack(ret interface{}, meta innerValue) {
	key := meta.ResponseKey
	str, err := valueToJson(ret)
	if err != nil {
		fmt.Printf("JSON Encode error: %s, value: %v\n", err, ret)
		return
	}

	err = s.redis.Set(key, str)
	if err != nil {
		fmt.Printf("Redis SET(%s %s) error: %s\n", key, str, err)
		return
	}

	_, err = s.redis.Publish(key, 1)
	if err != nil {
		fmt.Printf("Redis PUBLISH(%s, 1) error: %s\n", s.queueName, err)
		return
	}
}

/////////////////////////////////////////

type valueGeneratorFunc func()interface{}

type Client struct {
	netaddr string
	db int
	password string
	redis *godis.Client
	queueName string
	valueGenerator valueGeneratorFunc
}

func CreateClient(netaddr string, db int, password, queueName string, valueGenerator valueGeneratorFunc) *Client {
	redis := godis.New(netaddr, db, password)

	return &Client{
		netaddr: netaddr,
		db: db,
		password: password,
		redis: redis,
		queueName: queueName,
		valueGenerator: valueGenerator,
	}
}

func (c *Client) Close() error {
	return c.redis.Quit()
}

func (c *Client) Do(v interface{}) (interface{}, error) {
	idCountName := c.getIdCountName()
	idCount, err := c.redis.Incr(idCountName)
	if err != nil {
		return nil, err
	}
	key := c.getResponseKey(idCount)

	value := innerValue{
		ResponseKey: key,
		Data: v,
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(value)
	if err != nil {
		return nil, err
	}

	_, err = c.redis.Rpush(c.queueName, buf.String())
	if err != nil {
		return nil, err
	}
	_, err = c.redis.Publish(c.queueName, 0)
	if err != nil {
		return nil, err
	}

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

	buf = bytes.NewBuffer(retBytes)
	decoder := json.NewDecoder(buf)
	ret := c.valueGenerator()
	err = decoder.Decode(&ret)
	return ret, err
}

func (c *Client) getResponseKey(id int64) string {
	return fmt.Sprintf("%s:%d", c.queueName, id)
}

func (c *Client) getIdCountName() string {
	return fmt.Sprintf("%s:idcount", c.queueName)
}
