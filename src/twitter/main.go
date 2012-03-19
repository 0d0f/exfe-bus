package main

import (
	"twitter/service"
	"twitter/job"
	"config"
	"gobus"
	"gosque"
	"log"
)

const (
	queue_friendship = "twitter:friendship"
	queue_info = "twitter:userinfo"
	queue_tweet = "twitter:tweet"
	queue_dm = "twitter:directmessage"
)

func runService(c *twitter_job.Config) {
	friendship, err := gobus.CreateService(
		c.Redis.Netaddr,
		c.Redis.Db,
		c.Redis.Password,
		queue_friendship,
		&twitter_service.FriendshipsExists{})
	if err != nil {
		log.Fatal("FriendshipsExists service launch failed:", err)
	}

	go friendship.Serve(c.Service.Time_out)

	user := new(twitter_service.UsersShow)
	user.SiteUrl = c.Site_url
	info, err := gobus.CreateService(
		c.Redis.Netaddr,
		c.Redis.Db,
		c.Redis.Password,
		queue_info,
		user)
	if err != nil {
		log.Fatal("UsersShow service launch failed:", err)
	}

	go info.Serve(c.Service.Time_out)

	tweet, err := gobus.CreateService(
		c.Redis.Netaddr,
		c.Redis.Db,
		c.Redis.Password,
		queue_tweet,
		&twitter_service.StatusesUpdate{})
	if err != nil {
		log.Fatal("StatusesUpdate service launch failed:", err)
	}

	go tweet.Serve(c.Service.Time_out)

	d := new(twitter_service.DirectMessagesNew)
	d.SiteUrl = c.Site_url
	dm, err := gobus.CreateService(
		c.Redis.Netaddr,
		c.Redis.Db,
		c.Redis.Password,
		queue_dm,
		d)
	if err != nil {
		log.Fatal("DirectMessagesNew service launch failed:", err)
	}

	go dm.Serve(c.Service.Time_out)
}

func main() {
	log.SetPrefix("[TwitterSender]")
	log.Printf("Service start")

	var c twitter_job.Config
	config.LoadHelper("twitter.json", &c)

	runService(&c)

	client := gosque.CreateQueue(
		c.Redis.Netaddr,
		c.Redis.Db,
		c.Redis.Password,
		"resque:queue:twitter")

	sendtweet := gobus.CreateClient(
		c.Redis.Netaddr,
		c.Redis.Db,
		c.Redis.Password,
		queue_tweet)

	senddm := gobus.CreateClient(
		c.Redis.Netaddr,
		c.Redis.Db,
		c.Redis.Password,
		queue_dm)

	getinfo := gobus.CreateClient(
		c.Redis.Netaddr,
		c.Redis.Db,
		c.Redis.Password,
		queue_info)

	getfriendship := gobus.CreateClient(
		c.Redis.Netaddr,
		c.Redis.Db,
		c.Redis.Password,
		queue_friendship)

	recv := client.IncomingJob("twitter_job", twitter_job.TwitterSenderGenerator, 5e9)
	for {
		select {
		case job := <-recv:
			twitterSender := job.(*twitter_job.TwitterSender)
			twitterSender.Config = &c
			twitterSender.Sendtweet = sendtweet
			twitterSender.Senddm = senddm
			twitterSender.Getinfo = getinfo
			twitterSender.Getfriendship = getfriendship
			go func() {
				twitterSender.Do()
			}()
		}
	}
	log.Printf("Service stop")
}
