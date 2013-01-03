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
	dispatcher := gobus.NewDispatcher(gobus.NewTable(config.Dispatcher))
	sender, err := broker.NewSender(&config, dispatcher)
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
		tkMng, err := NewTokenManager(&config, db)
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

	if config.ExfeService.Services.ShortToken {
		shorttoken, err := NewShortToken(&config, db)
		if err != nil {
			log.Crit("shorttoken can't created: %s", err)
			os.Exit(-1)
		}

		count, err = bus.RegisterPath("/shorttoken", shorttoken)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register shorttoken %d methods.", count)

		count, err = bus.RegisterPath("/shorttoken/{key}", shorttoken)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register shorttoken/key %d methods.", count)
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

	if config.ExfeService.Services.Conversation {
		conversation, err := NewConversation_(&config, db, redis, dispatcher)
		if err != nil {
			log.Crit("conversation can't created: %s", err)
			os.Exit(-1)
		}

		count, err = bus.RegisterPath("/cross/{cross_id}/Conversation", conversation)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register conversation %d methods.", count)

		count, err = bus.RegisterPath("/cross/{cross_id}/Conversation/{post_id}", conversation)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register conversation/post %d methods.", count)

		count, err = bus.RegisterPath("/cross/{cross_id}/user/{user_id}/unread_count", conversation)
		if err != nil {
			log.Crit("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register conversation/unread %d methods.", count)
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
