package exfe_service

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

func LoadFile(filename string, config interface{}) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Site_url    string
	Site_api    string
	Site_img    string
	App_url     string
	Iom_url     string
	EmailName   string `json:"email_name"`
	EmailDomain string `json:"email_domain"`

	Redis struct {
		Netaddr  string
		Db       int
		Password string
	}
	Iom struct {
		Port int
	}
	Cross struct {
		Time_out time.Duration
		Delay    map[string]int
	}
	Bot struct {
		Iom_timeout   time.Duration
		Imap_time_out time.Duration
		Imap_host     string
		Imap_user     string
		Imap_password string
	}
	User struct {
		Time_out time.Duration
	}
	Twitter struct {
		Client_token  string
		Client_secret string
		Access_token  string
		Access_secret string
		Screen_name   string
		Time_out      time.Duration
	}
	Apn struct {
		Time_out time.Duration
		Cert     string
		Key      string
		Server   string
		Rootca   string
	}
	C2DM struct {
		Time_out time.Duration
		Email    string
		Password string
		Appid    string
	}
	Email struct {
		Time_out time.Duration
		Host     string
		Port     uint
		User     string
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

	LoadFile(configFile, &c)

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
