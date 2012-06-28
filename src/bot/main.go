package main

import (
	"exfe/service"
	"log"
)

func main() {
	config := exfe_service.InitConfig()
	log.SetPrefix("exfe.twitter_bot")

	log.Printf("service start")

	InitTwitter(config)

	processTwitter(config)
}
