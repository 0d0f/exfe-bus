package qutimer

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"time"
)

var redisPool = &redis.Pool{
	MaxIdle:     3,
	IdleTimeout: 30 * time.Minute,
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			return nil, err
		}
		return c, nil
	},
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	},
}

var prefix = "task:test"
var pushChannel = fmt.Sprintf("%s:action:push", prefix)
var key = "key"
var number = 10
var timer = fmt.Sprintf("%s:timer", prefix)
var queue = fmt.Sprintf("%s:queue:%s", prefix, key)
var queueOverwrite = fmt.Sprintf("%s:overwrite", queue)
var queueLocker = fmt.Sprintf("%s:locker", queue)
var queueStart = fmt.Sprintf("%s:start", queue)
var queueData = fmt.Sprintf("%s:data", queue)

func init() {
	rand.Seed(time.Now().UnixNano())
	prefix = fmt.Sprintf("task:test:%d.%d", time.Now().Unix(), rand.Intn(10000))
	pushChannel = fmt.Sprintf("%s:action:push", prefix)
	key = "key"
	number = 10
	timer = fmt.Sprintf("%s:timer", prefix)
	queue = fmt.Sprintf("%s:queue:%s", prefix, key)
	queueOverwrite = fmt.Sprintf("%s:overwrite", queue)
	queueLocker = fmt.Sprintf("%s:locker", queue)
	queueStart = fmt.Sprintf("%s:start", queue)
	queueData = fmt.Sprintf("%s:data", queue)

}

func clearQueue() {
	conn := redisPool.Get()
	defer conn.Close()

	conn.Do("DEL", timer)
	conn.Do("DEL", queue)
	conn.Do("DEL", queueOverwrite)
	conn.Do("DEL", queueLocker)
	conn.Do("DEL", queueStart)
	conn.Do("DEL", queueData)
}
