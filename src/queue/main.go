package main

import (
	"daemon"
	"fmt"
	"github.com/googollee/go-logger"
	"gobus"
	"model"
	"os"
)

func main() {
	var config model.Config
	output, quit := daemon.Init("exfe.json", &config)

	log, err := logger.New(output, "service bus", logger.Lshortfile)
	if err != nil {
		panic(err)
		return
	}
	config.Log = log

	servicesName := map[string]string{
		"conversation": "Conversation",
	}
	services := make(map[string]*gobus.Client)
	for k, v := range servicesName {
		s, err := gobus.NewClient(fmt.Sprintf("http://%s:%d/%s", config.ExfeService.Addr, config.ExfeService.Port, v))
		if err != nil {
			log.Crit("can't create gobus client for service %s: %s", k, err)
			os.Exit(-1)
		}
		services[k] = s
	}

	instant := NewInstant(services)

	url := fmt.Sprintf("http://%s:%d", config.ExfeQueue.Addr, config.ExfeQueue.Port)
	log.Info("start at %s", url)

	bus, err := gobus.NewServer(url, log)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	bus.Register(instant)

	go func() {
		<-quit
		log.Info("quit")
		os.Exit(-1)
		return
	}()
	bus.ListenAndServe()
}
