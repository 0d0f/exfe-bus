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

func getCallback(log *logger.SubLogger, config *model.Config) func(string, [][]byte) {
	return func(key string, datas [][]byte) {
		names := strings.SplitN(key, ",", 4)
		if len(names) != 4 {
			log.Crit("can't split service and method from key: %s", key)
			return
		}
		service, method, key := names[0], names[1], names[2]
		url := fmt.Sprintf("http://%s:%d/%s", config.ExfeService.Addr, config.ExfeService.Port, service)
		client, err := gobus.NewClient(url)
		if err != nil {
			log.Crit("can't create gobus client for service %s(%s): %s", service, url, err)
			return
		}

		arg := make([]interface{}, 0)
		for _, data := range datas {
			var d interface{}
			err := json.Unmarshal(data, &d)
			if err != nil {
				log.Err("can't unmarshal(%+v): %s", data, err)
				continue
			}
			if key != "" {
				arg = append(arg, d)
			} else {
				var i int
				err := client.Do(method, d, &i)
				if err != nil {
					log.Err("call %s|%s with %s failed: %s", service, method, string(data), err)
				}
			}
		}
		if key != "" {
			var i int
			err := client.Do(method, arg, &i)
			if err != nil {
				j, _ := json.Marshal(arg)
				log.Err("call %s|%s failed(%s) with %s", service, method, err, string(j))
			}
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

	url := fmt.Sprintf("http://%s:%d", config.ExfeQueue.Addr, config.ExfeQueue.Port)
	log.Info("start at %s", url)

	bus, err := gobus.NewServer(url, log)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	instant := NewInstant(&config)
	var count int
	count, err = bus.Register(instant)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register Instant %d methods.", count)

	for name, delayInSecond := range config.ExfeQueue.Head {
		head, headTomb := NewHead(delayInSecond, &config)
		tombs = append(tombs, headTomb)
		count, err = bus.RegisterName(name, head)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register %s %d methods.", name, count)
	}

	for name, delayInSecond := range config.ExfeQueue.Tail {
		tail, tailTomb := NewTail(delayInSecond, &config)
		tombs = append(tombs, tailTomb)
		count, err = bus.RegisterName(name, tail)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register %s %d methods.", name, count)
	}

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
