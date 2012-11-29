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

	redis := broker.NewRedisMultiplexer(&config)
	sender, err := broker.NewSender(&config)
	if err != nil {
		log.Crit("can't create sender: %s", err)
		os.Exit(-1)
		return
	}

	url := fmt.Sprintf("http://%s:%d", config.ExfeService.Addr, config.ExfeService.Port)
	log.Info("start at %s", url)

	bus, err := gobus.NewServer(url, log)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	var count int

	if config.ExfeService.Services.TokenManager {
		tkMng, err := NewTokenManager(&config)
		if err != nil {
			log.Crit("create token manager failed: %s", err)
			os.Exit(-1)
			return
		}

		count, err = bus.Register(tkMng)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register TokenManager %d methods.", count)
	}

	if config.ExfeService.Services.Iom {
		iom := NewIom(&config, redis)

		count, err = bus.Register(iom)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register IOM %d methods.", count)
	}

	if config.ExfeService.Services.Thirdpart {
		thirdpart, err := NewThirdpart(&config)
		if err != nil {
			log.Crit("create thirdpart failed: %s", err)
			os.Exit(-1)
			return
		}

		count, err = bus.Register(thirdpart)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register Thirdpart %d methods.", count)
	}

	if config.ExfeService.Services.Notifier {
		localTemplate, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
		if err != nil {
			log.Crit("load local template failed: %s", err)
			os.Exit(-1)
			return
		}

		conversation := NewConversation(localTemplate, &config, sender)
		count, err = bus.Register(conversation)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register Conversation %d methods.", count)

		cross := NewCross(localTemplate, &config, sender)
		count, err = bus.Register(cross)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register Cross %d methods.", count)

		user := NewUser(localTemplate, &config, sender)
		count, err = bus.Register(user)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register User %d methods.", count)
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
