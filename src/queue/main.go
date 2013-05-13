package main

import (
	"broker"
	"daemon"
	"fmt"
	"gobus"
	"launchpad.net/tomb"
	"logger"
	"model"
	"os"
)

func main() {
	var config model.Config
	_, quit := daemon.Init("exfe.json", &config)
	tombs := make([]*tomb.Tomb, 0)

	addr := fmt.Sprintf("%s:%d", config.ExfeQueue.Addr, config.ExfeQueue.Port)
	logger.NOTICE("start at %s", addr)

	bus, err := gobus.NewServer(addr)
	if err != nil {
		logger.ERROR("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	r, err := broker.NewRedisPool(&config)
	if err != nil {
		logger.ERROR("launch redis pool failed: %s", err)
		os.Exit(-1)
		return
	}

	q, err := NewQueue(&config, r)
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
