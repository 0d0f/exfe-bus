package gosque

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/simonz05/godis"
	"time"
)

type Client struct {
	redis     *godis.Client
	queueName string
}

func CreateQueue(netaddr string, db int, password, queueName string) *Client {
	return &Client{
		redis:     godis.New(netaddr, db, password),
		queueName: queueName,
	}
}

func (c *Client) Close() {
	c.redis.Quit()
}

func (c *Client) PutJob(name string, v interface{}) error {
	meta := metaType{
		Class: name,
		Args: v,
	}
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.Encode(meta)
	_, err := c.redis.Rpush(c.queueName, buf.String())
	return err
}

func (c *Client) jobLoop(jobRecv chan<- interface{}, name string, generateFunc func() interface{}, timeOut time.Duration) {
	for {
		queueLen, err := c.redis.Llen(c.queueName)
		if err != nil {
			fmt.Printf("Error redis(LLEN %s): %s\n", c.queueName, err)
			continue
		}
		if queueLen == 0 {
			time.Sleep(timeOut)
			continue
		}

		elem, err := c.redis.Lpop(c.queueName)
		if err != nil {
			fmt.Printf("Error redis(LPOP %s): %s\n", c.queueName, err)
			continue
		}

		buffer := bytes.NewBuffer(elem)
		decoder := json.NewDecoder(buffer)

		value := metaType{
			Args: generateFunc(),
		}
		err = decoder.Decode(&value)
		if err != nil {
			fmt.Printf("Error parse value(%s): %s\n", string(elem), err)
			continue
		}

		if value.Class != name {
			c.redis.Rpush(c.queueName, string(elem))
			continue
		}

		go func() {
			jobRecv <- value.Args
		}()
	}
}

func (c *Client) IncomingJob(name string, generateFunc func() interface{}, timeOut time.Duration) <-chan interface{} {
	jobChan := make(chan interface{})
	go c.jobLoop(jobChan, name, generateFunc, timeOut)
	return jobChan
}

type metaType struct {
	Class string
	Args interface{}
	Id   string
}
