package exfe_service

import (
	"time"
	"flag"
	"config"
	"os"
	"fmt"
)

type Config struct {
	Site_url string
	Site_api string

	Redis struct {
		Netaddr string
		Db int
		Password string
	}
	Cross struct {
		Time_out time.Duration
		Delay map[string]int
	}

	Twitter struct {
		Client_token string
		Client_secret string
		Access_token string
		Access_secret string
		Screen_name string
		Time_out time.Duration
	}
	Apn struct {
		Time_out time.Duration
		Cert string
		Key string
		Server string
	}
	C2DM struct {
		Time_out time.Duration
		Email string
		Password string
		Appid string
	}
	Email struct {
		Time_out time.Duration
		Host string
		Port uint
		User string
		Password string
	}
}

func InitConfig() *Config {
	var c Config

	var pidfile string
	var configFile string

	flag.StringVar(&pidfile, "pid", "", "Specify the pid file")
	flag.StringVar(&configFile, "config", "exfe.json", "Specify the configuration file")
	flag.Parse()

	config.LoadFile(configFile, &c)

	flag.Parse()
	if pidfile != "" {
		pid, err := os.Create(pidfile)
		if err != nil {
			panic(fmt.Sprintf("Can't create pid(%s): %s", pidfile, err))
		}
		pid.WriteString(fmt.Sprintf("%d", os.Getpid()))
	}

	return &c
}
