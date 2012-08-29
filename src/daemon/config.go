package daemon

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/googollee/go-log"
	"os"
	"os/signal"
	"syscall"
)

func Init(defaultConfig string, config interface{}) (loggerOutput log.OutType, quit <-chan os.Signal) {
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

	if syslog {
		loggerOutput = log.Syslog
	} else {
		loggerOutput = log.Stderr
	}

	sigChan := make(chan os.Signal)
	quit = sigChan
	signal.Notify(sigChan, syscall.SIGTERM)

	return
}
