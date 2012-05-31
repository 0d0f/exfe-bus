package main

import (
	"apn/service"
	"exfe/service"
	"gobus"
	"log/syslog"
	"fmt"
)

func main() {
	log, err := syslog.New(syslog.LOG_INFO, "exfe.apn")
	if err != nil {
		panic(err)
	}
	log.Info("Service start")

	c := exfe_service.InitConfig()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "iOSAPN")

	apn, err := apn_service.NewApn(c.Apn.Cert, c.Apn.Key, c.Apn.Server, log)
	if err != nil {
		log.Crit(fmt.Sprintf("Launch Apn service error: %s", err))
		panic(err)
	}
	server.Register(apn)

	server.Serve(c.Apn.Time_out * 1e9)
}
