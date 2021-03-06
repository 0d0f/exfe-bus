package daemon

import (
	"encoding/json"
	"flag"
	"fmt"
	"logger"
	"model"
	"os"
	"os/signal"
	"syscall"
)

func Init(defaultConfig string, config *model.Config) (quit <-chan os.Signal) {
	var pidfile string
	var configFile string
	var syslog bool

	flag.StringVar(&pidfile, "pid", "", "Specify the pid file")
	flag.StringVar(&configFile, "config", defaultConfig, "Specify the configuration file")
	flag.BoolVar(&syslog, "syslog", false, "Specify using syslog as log output")
	flag.Parse()

	f, err := os.Open(configFile)
	if err != nil {
		panic(fmt.Sprintf("open config %s error: %s", configFile, err))
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		panic(fmt.Sprintf("parse config %s error: %s", configFile, err))
	}

	flag.Parse()
	if pidfile != "" {
		pid, err := os.Create(pidfile)
		if err != nil {
			panic(fmt.Sprintf("create pid %s error: %s", pidfile, err))
		}
		pid.WriteString(fmt.Sprintf("%d", os.Getpid()))
	}

	sigChan := make(chan os.Signal)
	quit = sigChan
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGQUIT)

	logger.SetDebug(config.Debug)

	return
}
