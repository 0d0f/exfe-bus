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
	table, err := gobus.NewTable(config.Dispatcher)
	if err != nil {
		panic(err)
		return
	}
	dispatcher := gobus.NewDispatcher(table)
	sender, err := broker.NewSender(&config, dispatcher)
	if err != nil {
		log.Crit("can't create sender: %s", err)
		os.Exit(-1)
		return
	}
	// platform, err := NewPlatform(&config)
	// if err != nil {
	// 	log.Crit("can't create platform: %s", err)
	// 	os.Exit(-1)
	// 	return
	// }

	url := fmt.Sprintf("%s:%d", config.ExfeService.Addr, config.ExfeService.Port)
	log.Info("start at %s", url)

	bus, err := gobus.NewServer(url, log)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	if config.ExfeService.Services.TokenManager {
		tkMng, err := NewTokenManager(&config, db)
		if err != nil {
			log.Crit("create token manager failed: %s", err)
			os.Exit(-1)
			return
		}

		err = bus.Register(tkMng)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register TokenManager")
	}

	if config.ExfeService.Services.ShortToken {
		shorttoken, err := NewShortToken(&config, db)
		if err != nil {
			log.Crit("shorttoken can't created: %s", err)
			os.Exit(-1)
		}

		err = bus.Register(shorttoken)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register shorttoken")
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
		thirdpart, err := NewThirdpart(&config)
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
		localTemplate, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
		if err != nil {
			log.Crit("load local template failed: %s", err)
			os.Exit(-1)
			return
		}

		notifier := NewNotifier(localTemplate, &config, sender)
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
	err = bus.ListenAndServe()
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
}
