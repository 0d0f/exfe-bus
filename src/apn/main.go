package main

import (
	"apn/service"
	"exfe/service"
	"gobus"
	"log"
)

func main() {
	log.SetPrefix("exfe.apn")
	log.Print("Service start")

	c := exfe_service.InitConfig()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "iOS")

	apn, err := apn_service.NewApn(c.Apn.Cert, c.Apn.Key, c.Apn.Server, c.Apn.Rootca)
	if err != nil {
		log.Fatalf("Launch Apn service error: %s", err)
	}
	server.Register(apn)

	server.Serve(c.Apn.Time_out * 1e9)
}
