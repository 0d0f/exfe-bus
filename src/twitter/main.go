package main

import (
	"twitter/service"
	"twitter/job"
	"flag"
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

func runService(config *twitter_job.Config) {
	friendship := gobus.CreateService(
		config.Redis.Netaddr,
		config.Redis.Db,
		config.Redis.Password,
		queue_friendship,
		&twitter_service.FriendshipsExists{})

	go friendship.Serve(config.Service.Time_out)

	user := new(twitter_service.UsersShow)
	user.SiteUrl = config.Site_url
	info := gobus.CreateService(
		config.Redis.Netaddr,
		config.Redis.Db,
		config.Redis.Password,
		queue_info,
		user)

	go info.Serve(config.Service.Time_out)

	tweet := gobus.CreateService(
		config.Redis.Netaddr,
		config.Redis.Db,
		config.Redis.Password,
		queue_tweet,
		&twitter_service.StatusesUpdate{})

	go tweet.Serve(config.Service.Time_out)

	d := new(twitter_service.DirectMessagesNew)
	d.SiteUrl = config.Site_url
	dm := gobus.CreateService(
		config.Redis.Netaddr,
		config.Redis.Db,
		config.Redis.Password,
		queue_dm,
		d)

	go dm.Serve(config.Service.Time_out)
}

func main() {
	log.SetPrefix("[TwitterSender]")
	log.Printf("Service start")

	var configFile string
	flag.StringVar(&configFile, "config", "twitter.yaml", "Specify the configuration file")
	flag.Parse()
	config := twitter_job.Load(configFile)

	runService(config)

	client := gosque.CreateQueue(
		config.Redis.Netaddr,
		config.Redis.Db,
		config.Redis.Password,
		"resque:queue:twitter")

	sendtweet := gobus.CreateClient(
		config.Redis.Netaddr,
		config.Redis.Db,
		config.Redis.Password,
		queue_tweet)

	senddm := gobus.CreateClient(
		config.Redis.Netaddr,
		config.Redis.Db,
		config.Redis.Password,
		queue_dm)

	getinfo := gobus.CreateClient(
		config.Redis.Netaddr,
		config.Redis.Db,
		config.Redis.Password,
		queue_info)

	getfriendship := gobus.CreateClient(
		config.Redis.Netaddr,
		config.Redis.Db,
		config.Redis.Password,
		queue_friendship)

	recv := client.IncomingJob("twitter_job", twitter_job.TwitterSenderGenerator, 5e9)
	for {
		select {
		case job := <-recv:
			twitterSender := job.(*twitter_job.TwitterSender)
			twitterSender.Config = config
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
