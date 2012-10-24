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

	log, err := logger.New(output, "service bus", logger.Lshortfile)
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
	bus.Register(tkMng)
	bus.Register(iom)
	bus.Register(thirdpart)
	bus.Register(conversation)

	go func() {
		<-quit
		log.Info("quit")
		os.Exit(-1)
		return
	}()
	bus.ListenAndServe()
}
