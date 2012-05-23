package main

import (
	"apn/service"
	"exfe/service"
	"config"
	"gobus"
	"log/syslog"
	"flag"
	"fmt"
	"os"
)

func main() {
	log, err := syslog.New(syslog.LOG_INFO, "exfe.apn")
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

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "iOSAPN")

	apn, err := apn_service.NewApn(c.Apn.Cert, c.Apn.Key, c.Apn.Server, log)
	if err != nil {
		log.Crit(fmt.Sprintf("Launch Apn service error: %s", err))
		panic(err)
	}
	server.Register(apn)

	server.Serve(c.Apn.Time_out * 1e9)
}
