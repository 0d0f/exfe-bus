package main

import (
	"bot/mail"
	"broker"
	"daemon"
	"formatter"
	"github.com/googollee/go-logger"
	"launchpad.net/tomb"
	"model"
)

func main() {
	var config model.Config
	output, quit := daemon.Init("exfe.json", &config)

	log, err := logger.New(output, "bot")
	if err != nil {
		panic(err)
	}
	log.Notice("start")
	config.Log = log
	platform, err := broker.NewPlatform(&config)
	if err != nil {
		log.Crit("create platform failed: %s", err)
		return
	}
	templ, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
	if err != nil {
		log.Crit("create local template failed: %s", err)
		return
	}

	var tombs []*tomb.Tomb
	mail, err := mail.New(&config, templ, platform)
	if err != nil {
		log.Crit("create mail bot failed: %s", err)
		return
	}
	tombs = append(tombs, &mail.Tomb)
	go mail.Daemon()

	<-quit
	log.Notice("quiting...")

	for _, tomb := range tombs {
		tomb.Kill(nil)
		tomb.Wait()
	}
	log.Notice("quit")
}
