package main

import (
	"exfe/service"
	"log"
	"old_gobus"
	"twitter/service"
)

const (
	queue_friendship = "twitter:friendship"
	queue_info       = "twitter:userinfo"
	queue_tweet      = "twitter:tweet"
	queue_dm         = "twitter:directmessage"
)

func main() {
	log.SetPrefix("exfe.twitter")
	log.Print("Service start")

	c := exfe_service.InitConfig()

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "twitter")

	err := server.Register(twitter_service.NewFriendships())
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	err = server.Register(twitter_service.NewUsers(c.Site_api))
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	err = server.Register(twitter_service.NewStatuses())
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
	err = server.Register(twitter_service.NewDirectMessages(c.Site_api))
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}

	err = server.Serve(c.Twitter.Time_out * 1e9)
	if err != nil {
		log.Crit("gobus launch failed: %s", err)
		os.Exit(-1)
		return
	}
}
