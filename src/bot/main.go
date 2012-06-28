package main

import (
	"exfe/service"
	"fmt"
	"gobus"
	"log"
)

var config *exfe_service.Config
var helper string
var client *gobus.Client

func main() {
	config = exfe_service.InitConfig()
	helper = fmt.Sprintf("WRONG SYNTAX. Please enclose the 2-character mark in your reply to indicate mentioning 'X', e.g.:\n@%s Sure, be there or be square! #Z4", config.Twitter.Screen_name)
	client = gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "twitter")
	log.SetPrefix("exfe.twitter_bot")

	InitTwitter(config.Twitter.Screen_name)

	processTwitter(config)
}
