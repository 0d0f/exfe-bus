package gobus

import (
	"github.com/simonz05/godis"
	"fmt"
	"bytes"
	"encoding/json"
	"time"
)

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

func (c *Client) Do(method string, arg interface{}, reply interface{}, timeOut time.Duration) error {
	c.connRedis()
	defer c.closeRedis()

	meta, err := c.makeMeta(arg)
	if err != nil {
		return err
	}
	meta.Method = method
	meta.NeedReply = true

	err = c.send(meta)
	if err != nil {
		return err
	}
	return c.waitReply(meta.Id, reply, timeOut)
}

func (c *Client) Send(method string, arg interface{}, maxRetry int) error {
	c.connRedis()
	defer c.closeRedis()

	meta, err := c.makeMeta(arg)
	if err != nil {
		return err
	}
	meta.Method = method
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
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(m.Method)
	if err != nil {
		return err
	}

	err = encoder.Encode(m)
	if err != nil {
		return err
	}

	_, err = c.redis.Rpush(c.queueName, buf.String())
	if err != nil {
		return err
	}
	_, err = c.redis.Publish(c.queueName, 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) waitReply(id string, reply interface{}, timeOut time.Duration) error {
	sub := godis.NewSub(c.netaddr, c.db, c.password)
	sub.Subscribe(id)
	select {
	case _ = <-sub.Messages:
	case <-time.After(timeOut * time.Second):
		return fmt.Errorf("Timeout")
	}
	sub.Close()

	retBytes, err := c.redis.Get(id)
	if err != nil {
		return err
	}

	_, err = c.redis.Del(id)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(retBytes)
	decoder := json.NewDecoder(buf)
	var returnType string
	err = decoder.Decode(&returnType)
	if err != nil {
		return err
	}

	switch returnType {
	case "OK":
		err = decoder.Decode(reply)
		if err != nil {
			return err
		}
	case "Panic":
		var p interface{}
		err = decoder.Decode(&p)
		if err != nil {
			return err
		}
		panic(p)
	case "Error":
		var e string
		err = decoder.Decode(&e)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s", e)
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
