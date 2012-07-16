package main

import (
	"exfe/service"
	"log"
	"bot/email"
	"bot/twitter"
)

func main() {
	config := exfe_service.InitConfig()
	log.SetPrefix("exfe.bot")
	log.Printf("service start")

	quit := make(chan int)

	go twitter.Daemon(config, quit)
	go email.Daemon(config, quit)

	<-quit
	for i := 0; i < 1; i++ {
		quit <- 1
		<-quit
	}
}
