package main

import (
	"bot/email"
	"exfe/service"
	"log"
)

func main() {
	config := exfe_service.InitConfig()
	log.SetPrefix("exfe.bot")
	log.Printf("service start")

	quit := make(chan int)

	go email.Daemon(config, quit)

	<-quit
}
