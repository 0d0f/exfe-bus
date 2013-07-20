package main

import (
	"daemon"
	"fmt"
	"github.com/garyburd/redigo/redis"
	logger_ "github.com/googollee/go-logger"
	"gobus"
	"launchpad.net/tomb"
	"logger"
	"model"
	"os"
	"time"
)

func main() {
	var config model.Config
	output, quit := daemon.Init("exfe.json", &config)
	log, err := logger_.New(output, "service bus")
	if err != nil {
		panic(err)
		return
	}
	config.Log = log

	tombs := make([]*tomb.Tomb, 0)

	addr := fmt.Sprintf("%s:%d", config.ExfeQueue.Addr, config.ExfeQueue.Port)
	logger.NOTICE("start at %s", addr)

	bus, err := gobus.NewServer(addr)
	if err != nil {
		logger.ERROR("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 30 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Redis.Netaddr)
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

	q, err := NewQueue(&config, redisPool)
	if err != nil {
		logger.ERROR("launch queue failed: %s", err)
		os.Exit(-1)
		return
	}
	defer q.Quit()
	err = bus.RegisterRestful(q)
	if err != nil {
		logger.ERROR("register queue failed: %s", err)
		os.Exit(-1)
		return
	}
	logger.NOTICE("launch queue")

	go func() {
		<-quit
		for i, _ := range tombs {
			tombs[i].Kill(nil)
			tombs[i].Wait()
		}
		logger.NOTICE("quit")
		os.Exit(-1)
		return
	}()

	err = bus.ListenAndServe()
	if err != nil {
		logger.ERROR("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
}
