package main

import (
	"c2dm/service"
	"exfe/service"
	"gobus"
	"log"
)

func main() {
	log.SetPrefix("exfe.gcm")
	log.Print("Service start")

	c := exfe_service.InitConfig()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "Android")

	gcm := gcm_service.NewGCM(c.GCM.Key)
	server.Register(gcm)

	server.Serve(c.GCM.Time_out * 1e9)
}
