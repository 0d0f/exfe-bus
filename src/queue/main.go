package main

import (
	"daemon"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"gobus"
	"launchpad.net/tomb"
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

func (a PushArg) String() string {
	return fmt.Sprintf("{service:%s method:%s key:%s}", a.Service, a.Method, a.Key)
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
	tombs := make([]*tomb.Tomb, 0)

	log, err := logger.New(output, "service bus")
	if err != nil {
		panic(err)
		return
	}
	config.Log = log

	servicesName := []string{"Conversation"}
	services := make(map[string]*gobus.Client)
	for _, serviceName := range servicesName {
		s, err := gobus.NewClient(fmt.Sprintf("http://%s:%d/%s", config.ExfeService.Addr, config.ExfeService.Port, serviceName))
		if err != nil {
			log.Crit("can't create gobus client for service %s: %s", serviceName, err)
			os.Exit(-1)
		}
		services[serviceName] = s
	}

	instant := NewInstant(services)

	head10m, head10mTomb := NewHead10m(services, &config)
	tombs = append(tombs, head10mTomb)

	url := fmt.Sprintf("http://%s:%d", config.ExfeQueue.Addr, config.ExfeQueue.Port)
	log.Info("start at %s", url)

	bus, err := gobus.NewServer(url, log)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	var count int
	count, err = bus.Register(instant)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register Instant %d methods.", count)
	count, err = bus.Register(head10m)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register Head10m %d methods.", count)

	go func() {
		<-quit
		for i, _ := range tombs {
			tombs[i].Kill(nil)
			tombs[i].Wait()
		}
		log.Info("quit")
		os.Exit(-1)
		return
	}()
	err = bus.ListenAndServe()
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
}
