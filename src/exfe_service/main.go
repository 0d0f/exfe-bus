package main

import (
	"daemon"
	"fmt"
	"github.com/googollee/go-log"
	"github.com/googollee/go-mysql"
	"gobus"
	"os"
)

type Config struct {
	DB struct {
		Addr     string `json:"addr"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		DbName   string `json:"db_name"`
	} `json:"db"`
	ExfeService struct {
		Addr string `json:"addr"`
		Port int    `json:"port"`
	} `json:"exfe_service"`
	TokenManager struct {
		TableName string `json:"table_name"`
	} `json:"token_manager"`

	loggerOutput log.OutType
	loggerFlags  int
}

func main() {
	var config Config
	var quit <-chan os.Signal
	config.loggerOutput, quit = daemon.Init("exfe.json", &config)
	config.loggerFlags = log.Lshortfile

	l, err := log.New(config.loggerOutput, "service bus", config.loggerFlags)
	if err != nil {
		panic(err)
	}

	dbAddr := fmt.Sprintf("%s:%d", config.DB.Addr, config.DB.Port)
	db, err := mysql.DialTCP(dbAddr, config.DB.Username, config.DB.Password, config.DB.DbName)
	if err != nil {
		l.Crit("db connect failed: %s", err)
		os.Exit(-1)
	}

	tkMng, err := NewTokenManager(&config, db)
	if err != nil {
		l.Crit("create token manager failed: %s", err)
		os.Exit(-1)
	}

	url := fmt.Sprintf("http://%s:%d", config.ExfeService.Addr, config.ExfeService.Port)
	l.Info("start at %s", url)
	bus, err := gobus.NewGobusServer(url, l)
	if err != nil {
		l.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
	}
	bus.Register(tkMng)

	go func() {
		<-quit
		l.Info("quit")
		os.Exit(-1)
	}()
	bus.ListenAndServe()
}
