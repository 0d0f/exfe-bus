package main

import (
	"bot/email"
	"daemon"
	"github.com/googollee/go-logger"
	"model"
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

	log.Info("service start")

	tomb := email.Daemon(&config)

	<-quit
	tomb.Kill(nil)
	tomb.Wait()

	log.Info("quit")
}
