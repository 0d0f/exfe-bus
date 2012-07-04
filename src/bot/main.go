package main

import (
	"exfe/service"
	"log"
)

func main() {
	config := exfe_service.InitConfig()
	log.SetPrefix("exfe.bot")

	log.Printf("service start")

	InitTwitter(config)
	InitEmail(config)

	quit := make(chan int)

	go processTwitter(config, quit)
	go processEmail(quit)

	<-quit
}
