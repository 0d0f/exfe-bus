package main

import (
	"bot/mail"
	"broker"
	"daemon"
	"formatter"
	"launchpad.net/tomb"
	"logger"
	"model"
)

func main() {
	var config model.Config
	_, quit := daemon.Init("exfe.json", &config)
	logger.SetDebug(config.Debug)

	logger.NOTICE("bot start")
	platform, err := broker.NewPlatform(&config)
	if err != nil {
		logger.ERROR("create platform failed: %s", err)
		return
	}
	templ, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
	if err != nil {
		logger.ERROR("create local template failed: %s", err)
		return
	}

	db := broker.NewDBMultiplexer(&config)
	saver := NewCrossSaver(db)

	var tombs []*tomb.Tomb
	mail, err := mail.New(&config, templ, platform, saver)
	if err != nil {
		logger.ERROR("create mail bot failed: %s", err)
		return
	}
	tombs = append(tombs, &mail.Tomb)
	go mail.Daemon()

	<-quit
	logger.NOTICE("bot quiting...")

	for _, tomb := range tombs {
		tomb.Kill(nil)
		tomb.Wait()
	}
	logger.NOTICE("quit")
}
