package main

import (
	"daemon"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"gobus"
	"model"
	"os"
	"strings"
)

type PushArg struct {
	Service string      `json:"service"`
	Method  string      `json:"method"`
	Key     string      `json:"key"`
	Data    interface{} `json:"data"`
}

func getCallback(log *logger.SubLogger, services map[string]*gobus.Client) func(key string, datas [][]byte) {
	return func(key string, datas [][]byte) {
		names := strings.SplitN(key, ",", 3)
		if len(names) != 3 {
			log.Crit("can't split service and method from key: %s", key)
			return
		}
		serviceName, method := names[0], names[1]
		service, ok := services[serviceName]
		if !ok {
			log.Err("can't find service: %s", serviceName)
			return
		}
		d := make([]interface{}, 0)
		for i, _ := range datas {
			var data interface{}
			err := json.Unmarshal(datas[i], &data)
			if err != nil {
				log.Err("can't unmarshal(%+v): %s", datas[i], err)
			}
			d = append(d, data)
		}
		var i int
		err := service.Do("method", d, &i)
		if err != nil {
			log.Crit("call service %s method %s failed: %s", serviceName, method, err)
		}
	}
}

func main() {
	var config model.Config
	output, quit := daemon.Init("exfe.json", &config)
	quits := make([]chan int, 0)

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

	head10mQuit := make(chan int)
	head10m := NewHead10m(services, &config, head10mQuit)
	quits = append(quits, head10mQuit)

	url := fmt.Sprintf("http://%s:%d", config.ExfeQueue.Addr, config.ExfeQueue.Port)
	log.Info("start at %s", url)

	bus, err := gobus.NewServer(url, log)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	bus.Register(instant)
	bus.Register(head10m)

	go func() {
		<-quit
		for i, _ := range quits {
			quits[i] <- 1
		}
		log.Info("quit")
		os.Exit(-1)
		return
	}()
	bus.ListenAndServe()
}
