package main

import (
	"broker"
	"daemon"
	"fmt"
	"formatter"
	"github.com/googollee/go-logger"
	"gobus"
	"model"
	"net/http"
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

	db := broker.NewDBMultiplexer(&config)
	redis := broker.NewRedisMultiplexer(&config)

	localTemplate, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
	if err != nil {
		log.Crit("load local template failed: %s", err)
		os.Exit(-1)
		return
	}
	gate, err := NewGate(&config)
	if err != nil {
		log.Crit("can't create gate: %s", err)
		os.Exit(-1)
		return
	}
	streaming, err := NewStreaming(&config, gate)
	if err != nil {
		log.Crit("create streaming failed: %s", err)
		os.Exit(-1)
		return
	}
	platform, err := broker.NewPlatform(&config)
	if err != nil {
		log.Crit("can't create platform: %s", err)
		os.Exit(-1)
		return
	}

	gateAddr := fmt.Sprintf("%s:%d", config.ExfeGate.Addr, config.ExfeGate.Port)
	go func() {
		log.Info("launch gate at %s", gateAddr)
		err := http.ListenAndServe(gateAddr, streaming)
		log.Crit("launch gate failed: %s", err)
	}()

	url := fmt.Sprintf("%s:%d", config.ExfeService.Addr, config.ExfeService.Port)
	log.Info("start at %s", url)

	bus, err := gobus.NewServer(url, log)
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
	log.Info("register status")

	err = bus.Register(streaming)
	if err != nil {
		log.Crit("streaming register failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register streaming")

	if config.ExfeService.Services.Token {
		token, err := NewToken(&config, db)
		if err != nil {
			log.Crit("create token failed: %s", err)
			os.Exit(-1)
			return
		}
		err = bus.RegisterRestful(token)
		if err != nil {
			log.Crit("regiest token failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register Token")
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
		thirdpart, err := NewThirdpart(&config, streaming, platform)
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

	if config.ExfeService.Services.Notifier {
		notifier := NewNotifier(localTemplate, &config, platform)
		err = bus.Register(notifier)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register Notifier")
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
