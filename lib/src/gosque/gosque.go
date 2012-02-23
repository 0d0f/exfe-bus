package gosque

import (
	"encoding/json"
	"github.com/simonz05/godis"
	"bytes"
	"fmt"
	"time"
)

type Client struct {
	redis *godis.Client
	queueName string
}

func CreateQueue(netaddr string, db int, password, queueName string) *Client {
	return &Client{
		redis: godis.New(netaddr, db, password),
		queueName: queueName,
	}
}

func (c *Client) Close() {
	c.redis.Quit()
}

func (c *Client) PutJob(v interface{}) error {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.Encode(v)
	_, err := c.redis.Rpush(c.queueName, buf.String())
	return err
}

func (c *Client) jobLoop(jobRecv chan<- interface{}, generateFunc func() interface{}, timeOut time.Duration) {
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

		value := generateFunc()
		err = decoder.Decode(value)
		if err != nil {
			fmt.Printf("Error parse value: %s\n", string(elem))
			continue
		}

		jobRecv <- value
	}
}

func (c *Client) IncomingJob(generateFunc func() interface{}, timeOut time.Duration) <-chan interface{} {
	jobChan := make(chan interface{})
	go c.jobLoop(jobChan, generateFunc, timeOut)
	return jobChan
}
