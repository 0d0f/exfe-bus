package main

import (
	"daemon"
	"fmt"
	"formatter"
	"github.com/googollee/go-logger"
	"github.com/googollee/godis"
	"gobus"
	"model"
	"os"
)

func main() {
	var config model.Config
	output, quit := daemon.Init("exfe.json", &config)

	log, err := logger.New(output, "service bus")
	if err != nil {
		panic(err)
		return
	}
	config.Log = log

	redis := godis.New(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)

	tkMng, err := NewTokenManager(&config)
	if err != nil {
		log.Crit("create token manager failed: %s", err)
		os.Exit(-1)
		return
	}

	iom := NewIom(&config, redis)

	thirdpart, err := NewThirdpart(&config)
	if err != nil {
		log.Crit("create thirdpart failed: %s", err)
		os.Exit(-1)
		return
	}

	localTemplate, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
	if err != nil {
		log.Crit("load local template failed: %s", err)
		os.Exit(-1)
		return
	}
	conversation := NewConversation(localTemplate, &config)

	url := fmt.Sprintf("http://%s:%d", config.ExfeService.Addr, config.ExfeService.Port)
	log.Info("start at %s", url)

	bus, err := gobus.NewServer(url, log)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	var count int
	count, err = bus.Register(tkMng)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register TokenManager %d methods.", count)
	count, err = bus.Register(iom)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register IOM %d methods.", count)
	count, err = bus.Register(thirdpart)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register Thirdpart %d methods.", count)
	count, err = bus.Register(conversation)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register Conversation %d methods.", count)

	go func() {
		<-quit
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
