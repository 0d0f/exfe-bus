package main

import (
	"c2dm/service"
	"exfe/service"
	"gobus"
	"log"
)

func main() {
	log.SetPrefix("exfe.c2dm")
	log.Print("Service start")

	c := exfe_service.InitConfig()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "Android")

	c2dm, err := c2dm_service.NewC2DM(c.C2DM.Email, c.C2DM.Password, c.C2DM.Appid)
	if err != nil {
		log.Fatal("Launch service error: %s", err)
	}
	server.Register(c2dm)

	server.Serve(c.C2DM.Time_out * 1e9)
}
