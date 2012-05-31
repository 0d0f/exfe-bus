package main

import (
	"c2dm/service"
	"exfe/service"
	"gobus"
	"log/syslog"
	"fmt"
)

func main() {
	log, err := syslog.New(syslog.LOG_INFO, "exfe.c2dm")
	if err != nil {
		panic(err)
	}
	log.Info("Service start")

	c := exfe_service.InitConfig()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "Android")

	c2dm, err := c2dm_service.NewC2DM(c.C2DM.Email, c.C2DM.Password, c.C2DM.Appid, log)
	if err != nil {
		log.Crit(fmt.Sprintf("Launch service error: %s", err))
		panic(err)
	}
	server.Register(c2dm)

	server.Serve(c.C2DM.Time_out * 1e9)
}
