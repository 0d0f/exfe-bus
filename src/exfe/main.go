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

	cross := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "cross")
	cross.Register(exfe_service.NewCross(c))
	go cross.Serve(c.Cross.Time_out * 1e9)

	user := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "user")
	user.Register(exfe_service.NewUser(c))
	user.Serve(c.User.Time_out * 1e9)

	log.Print("Service stop")
}
