package main

import (
	"exfe/service"
	"config"
	"gobus"
	"log/syslog"
	"flag"
	"fmt"
	"os"
)

func main() {
	log, err := syslog.New(syslog.LOG_INFO, "exfe")
	if err != nil {
		panic(err)
	}
	log.Info("Service start")

	var c exfe_service.Config

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
			log.Crit(fmt.Sprintf("Can't create pid(%s): %s", pidfile, err))
			return
		}
		pid.WriteString(fmt.Sprintf("%d", os.Getpid()))
	}

	twitter := exfe_service.NewCrossTwitter(&c)
	go twitter.Serve()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "cross")
	server.Register(exfe_service.NewCross(&c))
	server.Serve(c.Cross.Time_out * 1e9)

	log.Info("Service stop")
}
