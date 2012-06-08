package main

import (
	"exfe/service"
	"gobus"
	"log"
)

func main() {
	log.SetPrefix("exfe")
	log.Print("Service start")

	c := exfe_service.InitConfig()

	twitter := exfe_service.NewCrossTwitter(c)
	go twitter.Serve()

	push := exfe_service.NewCrossPush(c)
	go push.Serve()

	email := exfe_service.NewCrossEmail(c)
	go email.Serve()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "cross")
	server.Register(exfe_service.NewCross(c))
	server.Register(exfe_service.NewAuthentication(c))
	server.Serve(c.Cross.Time_out * 1e9)

	log.Print("Service stop")
}
