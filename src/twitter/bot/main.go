package main

import (
	"exfe/service"
	"twitter/service"
	"fmt"
	"gobus"
	"log/syslog"
)

var log *syslog.Writer
var config *exfe_service.Config
var helper string
var client *gobus.Client

func sendHelp(screen_name string) {
	f := &twitter_service.FriendshipsExistsArg{
		ClientToken:  config.Twitter.Client_token,
		ClientSecret: config.Twitter.Client_secret,
		AccessToken:  config.Twitter.Access_token,
		AccessSecret: config.Twitter.Access_secret,
		UserA:        screen_name,
		UserB:        config.Twitter.Screen_name,
	}
	var isFriend bool
	err := client.Do("GetFriendship", f, &isFriend, 10)
	if err != nil {
		log.Err(fmt.Sprintf("Can't require user %s friendship: %s", screen_name, err))
		isFriend = false
	}

	if isFriend {
		dm := &twitter_service.DirectMessagesNewArg{
			ClientToken:  config.Twitter.Client_token,
			ClientSecret: config.Twitter.Client_secret,
			AccessToken:  config.Twitter.Access_token,
			AccessSecret: config.Twitter.Access_secret,
			Message:      helper,
			ToUserName:   &screen_name,
		}
		client.Send("SendDM", dm, 5)
	} else {
		tweet := &twitter_service.StatusesUpdateArg{
			ClientToken:  config.Twitter.Client_token,
			ClientSecret: config.Twitter.Client_secret,
			AccessToken:  config.Twitter.Access_token,
			AccessSecret: config.Twitter.Access_secret,
			Tweet:        fmt.Sprintf("@%s %s", screen_name, helper),
		}
		client.Send("SendTweet", tweet, 5)
	}
}

func main() {
	config = exfe_service.InitConfig()
	helper = fmt.Sprintf("WRONG SYNTAX. Please enclose the 2-character hash in your reply to indicate mentioning X, e.g.:\n @%s Sure, be there or be square! #Z4", config.Twitter.Screen_name)
	client = gobus.CreateClient(config.Redis.Netaddr, config.Redis.Db, config.Redis.Password, "twitter")
	var err error
	log, err = syslog.New(syslog.LOG_INFO, "exfe.twitter_bot")
	if err != nil {
		panic(err)
	}

	Init(config.Twitter.Screen_name)

	c, _ := connStreaming(config.Twitter.Client_token, config.Twitter.Client_secret, config.Twitter.Access_token, config.Twitter.Access_secret)

	for t := range c {
		hash, post := t.parse()
		time := t.created_at()
		external_id := t.external_id()
		screen_name := t.screen_name()

		fmt.Println(hash, time, external_id, screen_name, post)

		if screen_name == "" {
			continue
		}

		if hash == "" && post != "" {
			sendHelp(screen_name)
			continue
		}
	}
}
