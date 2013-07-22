package main

import (
	"bot/mail"
	"broker"
	"daemon"
	"database/sql"
	"fmt"
	"formatter"
	_ "github.com/go-sql-driver/mysql"
	"launchpad.net/tomb"
	"logger"
	"model"
)

func main() {
	var config model.Config
	quit := daemon.Init("exfe.json", &config)
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

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&autocommit=true",
		config.DB.Username, config.DB.Password, config.DB.Addr, config.DB.Port, config.DB.DbName))
	if err != nil {
		logger.ERROR("mysql error:", err)
		return
	}
	saver := broker.NewKVSaver(db)

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
