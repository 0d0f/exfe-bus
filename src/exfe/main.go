package main

import (
	"exfe/service"
	"gobus"
	"log/syslog"
)

func main() {
	log, err := syslog.New(syslog.LOG_INFO, "exfe")
	if err != nil {
		panic(err)
	}
	log.Info("Service start")

	c := exfe_service.InitConfig()

	twitter := exfe_service.NewCrossTwitter(c)
	go twitter.Serve()

	apn := exfe_service.NewCrossApn(c)
	go apn.Serve()

	email := exfe_service.NewCrossEmail(c)
	go email.Serve()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "cross")
	server.Register(exfe_service.NewCross(c))
	server.Register(exfe_service.NewAuthentication(c))
	server.Serve(c.Cross.Time_out * 1e9)

	log.Info("Service stop")
}
