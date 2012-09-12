package main

import (
	"configure"
	"daemon"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/googollee/godis"
	"gobus"
	"os"
)

func main() {
	var config configure.Config
	output, quit := daemon.Init("exfe.json", &config)

	log, err := logger.New(output, "service bus", logger.Lshortfile)
	if err != nil {
		panic(err)
	}
	config.Log = log

	redis := godis.New(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password)

	tkMng, err := NewTokenManager(&config)
	if err != nil {
		log.Crit("create token manager failed: %s", err)
		os.Exit(-1)
	}

	iom := NewIom(&config, redis)

	url := fmt.Sprintf("http://%s:%d", config.ExfeService.Addr, config.ExfeService.Port)
	log.Info("start at %s", url)

	bus, err := gobus.NewServer(url, log)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
	}
	bus.Register(tkMng)
	bus.Register(iom)

	go func() {
		<-quit
		log.Info("quit")
		os.Exit(-1)
	}()
	bus.ListenAndServe()
}
