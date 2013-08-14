package qutimer

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"io/ioutil"
	"logger"
	"math/rand"
	"time"
)

type Client struct {
	pool     *redis.Pool
	prefix   string
	timeout  int
	push     *redis.Script
	get      *redis.Script
	schedule *redis.Script

	rand *rand.Rand
}

func New(pool *redis.Pool, prefix string, timeout int) (*Client, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("invalid timeout: %d", timeout)
	}
	return &Client{
		pool:     pool,
		prefix:   prefix,
		timeout:  timeout,
		push:     redis.NewScript(2, pushScript),
		get:      redis.NewScript(2, getScript),
		schedule: redis.NewScript(1, scheduleScript),
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

const pushScript = `
local prefix = KEYS[1]
local key    = KEYS[2]
local task      = ARGV[1]
local ontime    = ARGV[2]
local overwrite = ARGV[3]
local timer        = prefix..":timer"
local push_channel = prefix..":action:push"
local queue        = prefix..":queue:"..key
local queue_overwrite = queue..":overwrite"
local queue_data      = queue..":data"

if overwrite == "always" or redis.call("LLEN", queue) == 0 or redis.call("EXISTS", queue_overwrite) ~= 0 then
	redis.call("ZADD", timer, ontime, key)
end
redis.call("DEL", queue_overwrite)

redis.call("RPUSH", queue, task)
for i=4, #ARGV, 2 do
	local k = ARGV[i]
	local v = ARGV[i+1]
	if k and v then
		redis.call("HSET", queue_data, k, v)
	end
end

redis.call("PUBLISH", push_channel, "insert")`

func (c *Client) Push(key string, task io.Reader, data map[string]string, ontime int64, overwrite string) error {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := ioutil.ReadAll(task)
	if err != nil {
		return err
	}
	if ontime <= 0 {
		ontime = time.Now().Unix()
	}
	argv := []interface{}{c.prefix, key, b, ontime, overwrite}
	for k, v := range data {
		argv = append(argv, k)
		argv = append(argv, v)
	}
	if _, err := c.push.Do(conn, argv...); err != nil {
		return err
	}
	return nil
}

const getScript = `
local prefix = KEYS[1]
local key    = KEYS[2]
local number  = tonumber(ARGV[1])
local timeout = tonumber(ARGV[2])
local random  = ARGV[3]
local now     = tonumber(ARGV[4])
local timer = prefix..":timer"
local queue = prefix..":queue:"..key
local queue_overwrite = queue..":overwrite"
local queue_locker    = queue..":locker"
local queue_start     = queue..":start"

redis.call("SET", queue_locker, random, "EX", timeout, "NX")
local v = redis.call("GET", queue_locker)
if v ~= random then
	return {err="locked"}
end

redis.call("DEL", queue_start)
redis.call("SET", queue_overwrite, "1")
redis.call("ZADD", timer, now+timeout, key)
if number > 0 then
	number = number - 1
end
return redis.call("LRANGE", queue, 0, number)`

func (c *Client) Get(key string, number int) (*Task, error) {
	conn := c.pool.Get()
	defer conn.Close()

	now := time.Now().Unix()
	random := fmt.Sprintf("%d.%d", now, c.rand.Int())
	reply, err := redis.Values(c.get.Do(conn, c.prefix, key, number, c.timeout, random, now))
	if err != nil {
		return nil, err
	}
	return newTask(c.pool, c.prefix, key, number, reply), nil
}

const scheduleScript = `
local prefix = KEYS[1]
local now = ARGV[1]
local timer = prefix..":timer"

while true do
	local first = redis.call("ZRANGEBYSCORE", timer, "-INF", "+INF", "WITHSCORES", "LIMIT", 0, 1)
	if next(first) == nil then
		return {-1, ""}
	end
	local till = first[2] - now
	if till > 0 then
		return {till, ""}
	end
	local key = first[1]
	local queue = prefix..":queue:"..key
	local queue_start     = queue..":start"
	local queue_overwrite = queue..":overwrite"
	local queue_data      = queue..":data"
	local start = redis.call("GET", queue_start)
	if start and now - start >= 5 then
		redis.call("DEL", queue, queue_start, queue_data)
		redis.call("ZREM", timer, key)
	else
		redis.call("SET", queue_start, now, "NX")
		return {0, first[1]}
	end
end`

func (c *Client) Schedule() (string, time.Duration, error) {
	conn := c.pool.Get()
	defer conn.Close()

	reply, err := redis.Values(c.schedule.Do(conn, c.prefix, time.Now().Unix()))
	if err != nil {
		return "", time.Duration(c.timeout) * time.Second, err
	}
	next, err := redis.Int64(reply[0], nil)
	if err != nil {
		return "", time.Duration(c.timeout) * time.Second, err
	}
	if next != 0 {
		if next < 0 {
			return "", time.Duration(c.timeout) * time.Second, nil
		}
		return "", time.Duration(next) * time.Second, nil
	}

	key, err := redis.String(reply[1], nil)
	if err != nil {
		return "", time.Duration(c.timeout) * time.Second, err
	}
	return key, 0, nil
}

func (c *Client) Wait(timeout time.Duration) error {
	conn, err := c.pool.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch := make(chan int)
	go func() {
		defer close(ch)
		pubsub := redis.PubSubConn{conn}
		pushChannel := fmt.Sprintf("%s:action:push", c.prefix)
		if err = pubsub.Subscribe(pushChannel); err != nil {
			logger.ERROR("can't subscribe push channel: %s", err)
			return
		}
		for {
			n := pubsub.Receive()
			switch n.(type) {
			case redis.Message:
				return
			case redis.PMessage:
				return
			case error:
				return
			}
		}
	}()

	select {
	case <-ch:
	case <-time.After(timeout):
	}
	return err
}
