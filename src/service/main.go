package main

import (
	"broker"
	"daemon"
	"database/sql"
	"fmt"
	"formatter"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"gobus"
	"logger"
	"model"
	"notifier"
	"os"
	"routex"
	"splitter"
	"time"
	"token"
)

func main() {
	var config model.Config
	quit := daemon.Init("exfe.json", &config)

	if config.Proxy != "" {
		broker.SetProxy(config.Proxy)
	}

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

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 30 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Redis.Netaddr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	cachePool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 30 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisCache.Netaddr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
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
	logger.NOTICE("start at %s", addr)

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
	logger.NOTICE("register Status")

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
		logger.NOTICE("register %s", name)
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

	if config.ExfeService.Services.Thirdpart {
		poster, err := registerThirdpart(&config, platform)
		register("poster", poster, err)
	}

	if config.ExfeService.Services.Notifier {
		err := notifier.SetupResponse(&config, notifier.NewResponseSaver(cachePool))
		if err != nil {
			logger.ERROR("can't setup response")
			return
		}
		user := notifier.NewUser(localTemplate, &config, platform)
		register("notifier/user", user, nil)
		cross := notifier.NewCross(localTemplate, &config, platform)
		register("notifier/cross", cross, nil)
		routex := notifier.NewRoutex(localTemplate, &config, platform)
		register("notifier/routex", routex, nil)
	}

	if config.ExfeService.Services.Iom {
		iom := NewIom(&config, redisPool)

		err = bus.Register(iom)
		if err != nil {
			logger.ERROR("gobus launch failed: %s", err)
			os.Exit(-1)
			return
		}
		logger.NOTICE("register IOM")
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
		logger.NOTICE("register Thirdpart")
	}

	if config.ExfeService.Services.Routex {
		routexSaver := routex.NewRoutexSaver(database)
		breadcrumbsSaver := routex.NewBreadcrumbsSaver(database)
		breadcrumbsCache := routex.NewBreadcrumbCacheSaver(cachePool)
		geomarksSaver := &routex.GeomarksSaver{database}
		c := routex.NewGeoConversion(database)
		rx, err := routex.New(routexSaver, breadcrumbsCache, breadcrumbsSaver, geomarksSaver, c, platform, &config)
		register("routex", rx, err)
	}

	go func() {
		<-quit
		logger.NOTICE("quit")
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
