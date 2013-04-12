package main

import (
	"broker"
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
		url := fmt.Sprintf("http://%s:%d/%s?method=%s", config.ExfeService.Addr, config.ExfeService.Port, service, method)
		client := gobus.NewClient(new(gobus.JSON))

		arg := make([]interface{}, 0)
		for _, data := range datas {
			var d interface{}
			err := json.Unmarshal(data, &d)
			if err != nil {
				log.Err("can't unmarshal %s(%+v)", err, data)
				continue
			}
			if key != "" {
				arg = append(arg, d)
			} else {
				var i int
				err := client.Do(url, "POST", d, &i)
				if err != nil {
					log.Err("call %s|%s failed(%s) with %s", service, method, err, string(data))
				}
			}
		}
		if key != "" {
			var i int
			err := client.Do(url, "POST", arg, &i)
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

	// table, err := gobus.NewTable(config.Dispatcher)
	// if err != nil {
	// 	log.Crit("create gobus table failed: %s", err)
	// 	os.Exit(-1)
	// 	return
	// }
	// dispatcher := gobus.NewDispatcher(table)

	addr := fmt.Sprintf("%s:%d", config.ExfeQueue.Addr, config.ExfeQueue.Port)
	log.Info("start at %s", addr)

	bus, err := gobus.NewServer(addr, log)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	// r, err := broker.NewRedisPool(&config)
	// if err != nil {
	// 	log.Crit("launch redis pool failed: %s", err)
	// 	os.Exit(-1)
	// 	return
	// }

	// q, err := NewQueue(&config, r, dispatcher)
	// if err != nil {
	// 	log.Crit("launch queue failed: %s", err)
	// 	os.Exit(-1)
	// 	return
	// }
	// defer q.Quit()
	// err = bus.RegisterRestful(q)
	// if err != nil {
	// 	log.Crit("register queue failed: %s", err)
	// 	os.Exit(-1)
	// 	return
	// }
	// log.Notice("launch queue")

	instant := NewInstant(&config)
	var count int
	err = bus.Register(instant)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register Instant %d methods.", count)

	for name, delayInSecond := range config.ExfeQueue.Head {
		head, headTomb := NewHead(delayInSecond, name, &config)
		tombs = append(tombs, headTomb)
		err = bus.Register(head)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register %s %d methods.", name, count)
	}

	for name, delayInSecond := range config.ExfeQueue.Tail {
		tail, tailTomb := NewTail(delayInSecond, name, &config)
		tombs = append(tombs, tailTomb)
		err = bus.Register(tail)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register %s %d methods.", name, count)
	}

	redis := broker.NewRedisImp(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)
	queue, err := NewQueue_(&config, redis)
	if err != nil {
		log.Crit("queue launch failed: %s", err)
		os.Exit(-1)
		return
	}
	err = bus.Register(queue)
	if err != nil {
		log.Crit("register queue failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register queue methods.")

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
	defer func() {
		re := recover()
		if re != nil {
			log.Crit("crash: %s", re)
		}
	}()
	err = bus.ListenAndServe()
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
}
