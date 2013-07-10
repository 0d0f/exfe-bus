package main

import (
	"broker"
	"daemon"
	"database/sql"
	"fmt"
	"formatter"
	_ "github.com/go-sql-driver/mysql"
	l "github.com/googollee/go-logger"
	"gobus"
	"logger"
	"model"
	"os"
	"routex"
	"splitter"
	"token"
)

func main() {
	var config model.Config
	output, quit := daemon.Init("exfe.json", &config)

	if config.Proxy != "" {
		broker.SetProxy(config.Proxy)
	}

	log, err := l.New(output, "service bus")
	if err != nil {
		panic(err)
		return
	}
	config.Log = log

	database, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&autocommit=true",
		config.DB.Username, config.DB.Password, config.DB.Addr, config.DB.Port, config.DB.DbName))
	if err != nil {
		logger.ERROR("mysql error:", err)
		os.Exit(-1)
		return
	}
	defer database.Close()
	err = database.Ping()
	if err != nil {
		logger.ERROR("mysql error:", err)
		os.Exit(-1)
		return
	}

	redis_ := broker.NewRedisMultiplexer(&config)
	redis, err := broker.NewRedisPool(&config)
	if err != nil {
		logger.ERROR("redis connect error: %s", err)
		return
	}

	localTemplate, err := formatter.NewLocalTemplate(config.TemplatePath, config.DefaultLang)
	if err != nil {
		logger.ERROR("load local template failed: %s", err)
		os.Exit(-1)
		return
	}
	platform, err := broker.NewPlatform(&config)
	if err != nil {
		logger.ERROR("can't create platform: %s", err)
		os.Exit(-1)
		return
	}

	addr := fmt.Sprintf("%s:%d", config.ExfeService.Addr, config.ExfeService.Port)
	log.Info("start at %s", addr)

	bus, err := gobus.NewServer(addr)
	if err != nil {
		logger.ERROR("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	status := NewStatus()
	err = bus.Register(status)
	if err != nil {
		logger.ERROR("status register failed: %s", err)
		os.Exit(-1)
		return
	}
	log.Info("register Status")

	register := func(name string, service interface{}, err error) {
		if err != nil {
			logger.ERROR("create %s failed: %s", name, err)
			os.Exit(-1)
			return
		}
		err = bus.RegisterRestful(service)
		if err != nil {
			logger.ERROR("regiest %s failed: %s", name, err)
			os.Exit(-1)
			return
		}
		log.Info("register %s", name)
	}

	if config.ExfeService.Services.Live {
		live, err := NewLive(&config, platform)
		register("live", live, err)
	}

	if config.ExfeService.Services.Token {
		repo, err := NewTokenRepo(&config, database)
		if err != nil {
			logger.ERROR("can't create token repo: %s", err)
			os.Exit(-1)
			return
		}
		token := token.New(repo)
		register("token", token, nil)
	}

	if config.ExfeService.Services.Splitter {
		splitter := splitter.NewSplitter(&config)
		register("splitter", splitter, nil)
	}

	if config.ExfeService.Services.Notifier {
		notifier, err := NewV3Notifier(localTemplate, &config, platform)
		register("notifier", notifier, err)
	}

	if config.ExfeService.Services.Thirdpart {
		poster, err := registerThirdpart(&config, platform)
		register("poster", poster, err)
	}

	if config.ExfeService.Services.Iom {
		iom := NewIom(&config, redis_)

		err = bus.Register(iom)
		if err != nil {
			logger.ERROR("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register IOM")
	}

	if config.ExfeService.Services.Thirdpart {
		thirdpart, err := NewThirdpart(&config, platform)
		if err != nil {
			logger.ERROR("create thirdpart failed: %s", err)
			os.Exit(-1)
			return
		}

		err = bus.Register(thirdpart)
		if err != nil {
			logger.ERROR("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		log.Info("register Thirdpart")
	}

	if config.ExfeService.Services.Routex {
		location := &routex.LocationSaver{redis}
		route := &routex.RouteSaver{database}
		routex := routex.New(location, route, platform, &config)
		register("routex", routex, nil)
	}

	go func() {
		<-quit
		log.Info("quit")
		os.Exit(-1)
		return
	}()
	defer func() {
		re := recover()
		logger.ERROR("crashed: %s", re)
	}()
	err = bus.ListenAndServe()
	if err != nil {
		logger.ERROR("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
}
