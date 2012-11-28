package main

import (
	"bot/email"
	"broker"
	"daemon"
	"formatter"
	"github.com/googollee/go-logger"
	"model"
	"os"
)

func main() {
	var config model.Config
	output, quit := daemon.Init("exfe.json", &config)

	log, err := logger.New(output, "bot")
	if err != nil {
		panic(err)
		return
	}
	config.Log = log

	localTemplate, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
	if err != nil {
		log.Crit("load local template failed: %s", err)
		os.Exit(-1)
		return
	}
	sender, err := broker.NewSender(&config)
	if err != nil {
		log.Crit("create gobus client failed: %s", err)
		os.Exit(-1)
		return
	}

	log.Info("start")

	tomb := email.Daemon(&config, localTemplate, sender)

	<-quit
	tomb.Kill(nil)
	tomb.Wait()

	log.Info("quit")
}
