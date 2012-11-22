package main

import (
	"bot/email"
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
	url := fmt.Sprintf("http://%s:%d/Thirdpart", config.ExfeService.Addr, config.ExfeService.Port)
	client, err := gobus.NewClient(url)
	if err != nil {
		log.Crit("create gobus client failed: %s", err)
		os.Exit(-1)
		return
	}

	tomb := email.Daemon(&config, localTemplate, client)

	<-quit
	tomb.Kill(nil)
	tomb.Wait()

	log.Info("quit")
}
