package main

import (
	"exfe/service"
	"log"
	"bot/email"
)

func main() {
	config := exfe_service.InitConfig()
	log.SetPrefix("exfe.bot")

	log.Printf("service start")

	InitTwitter(config)

	quit := make(chan int)

	go processTwitter(config, quit)
	go email.Daemon(config, quit)

	<-quit
	for i := 0; i < 1; i++ {
		quit <- 1
		<-quit
	}
}
