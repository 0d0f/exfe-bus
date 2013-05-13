package main

import (
	"broker"
	"daemon"
	"fmt"
	"formatter"
	"github.com/googollee/go-logger"
	"gobus"
	"model"
	"os"
	"splitter"
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

	db := broker.NewDBMultiplexer(&config)
	redis := broker.NewRedisMultiplexer(&config)

	localTemplate, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
	if err != nil {
		log.Crit("load local template failed: %s", err)
		os.Exit(-1)
		return
	}
	platform, err := broker.NewPlatform(&config)
	if err != nil {
		log.Crit("can't create platform: %s", err)
		os.Exit(-1)
		return
	}

	addr := fmt.Sprintf("%s:%d", config.ExfeService.Addr, config.ExfeService.Port)
	log.Info("start at %s", addr)

	bus, err := gobus.NewServer(addr)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	status := NewStatus()
	err = bus.Register(status)
	if err != nil {
		log.Crit("status register failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register Status")

	register := func(name string, service interface{}, err error) {
		if err != nil {
			log.Crit("create %s failed: %s", name, err)
			os.Exit(-1)
			return
		}
		err = bus.RegisterRestful(service)
		if err != nil {
			log.Crit("regiest %s failed: %s", name, err)
			os.Exit(-1)
			return
		}
		log.Info("register %s", name)
	}

	if config.ExfeService.Services.Live {
		live, err := NewLive(&config, platform)
		register("live", live, err)
	}

	if config.ExfeService.Services.Token {
		token, err := NewToken(&config, db)
		register("token", token, err)
	}

	if config.ExfeService.Services.Splitter {
		splitter := splitter.NewSplitter(&config)
		register("splitter", splitter, nil)
	}

	if config.ExfeService.Services.Notifier {
		notifier, err := NewV3Notifier(localTemplate, &config, platform)
		register("notifier", notifier, err)
	}

	if config.ExfeService.Services.Iom {
		iom := NewIom(&config, redis)

		err = bus.Register(iom)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register IOM")
	}

	if config.ExfeService.Services.Thirdpart {
		thirdpart, err := NewThirdpart(&config, platform)
		if err != nil {
			log.Crit("create thirdpart failed: %s", err)
			os.Exit(-1)
			return
		}

		err = bus.Register(thirdpart)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register Thirdpart")
	}

	go func() {
		<-quit
		log.Info("quit")
		os.Exit(-1)
		return
	}()
	defer func() {
		re := recover()
		log.Crit("crashed: %s", re)
	}()
	err = bus.ListenAndServe()
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
}
