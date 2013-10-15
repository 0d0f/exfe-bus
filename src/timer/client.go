package timer

import (
	"github.com/garyburd/redigo/redis"
	"strings"
)

type Client struct {
	pool       *redis.Pool
	prefix     string
	setName    string
	notifyName string
	dataName   string
}

func NewClient(pool *redis.Pool, prefix string) *Client {
	return &Client{
		pool:       pool,
		prefix:     prefix,
		setName:    sortedSetname(prefix),
		notifyName: notifyName(prefix),
		dataName:   dataName(prefix),
	}
}

func (c *Client) Send(ontime int64, name string, data string) error {
	conn := c.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("HSET", c.dataName, name, data); err != nil {
		return err
	}
	if _, err := conn.Do("ZADD", c.setName, ontime, name); err != nil {
		return err
	}
	if _, err := conn.Do("PUBLISH", c.notifyName, ""); err != nil {
		return err
	}
	return nil
}

func (c *Client) PListen(pattern string, ch chan Event) error {
	conn, err := c.pool.Dial()
	if err != nil {
		return err
	}
	pubsub := redis.PubSubConn{conn}
	if err := pubsub.PSubscribe(c.prefix + pattern); err != nil {
		return err
	}
	for {
		d := pubsub.Receive()
		var event Event
		switch n := d.(type) {
		case redis.Message:
			event = Event{
				Name: n.Channel,
				Data: n.Data,
			}
		case redis.PMessage:
			event = Event{
				Name: n.Channel,
				Data: n.Data,
			}
		case error:
			return n
		default:
			continue
		}
		if !strings.HasPrefix(event.Name, c.prefix) {
			continue
		}
		event.Name = event.Name[len(c.prefix):]
		ch <- event
	}
	return nil
}

func (c *Client) Close() error {
	return nil
}
