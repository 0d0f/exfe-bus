package main

import (
	"twitter/service"
	"twitter/job"
	"config"
	"gobus"
	"gosque"
	"log"
	"flag"
	"fmt"
	"os"
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

	var pidfile string
	var configFile string

	flag.StringVar(&pidfile, "pid", "", "Specify the pid file")
	flag.StringVar(&configFile, "config", "twitter.json", "Specify the configuration file")
	flag.Parse()

	config.LoadFile(configFile, &c)

	flag.Parse()
	if pidfile != "" {
		pid, err := os.Create(pidfile)
		if err != nil {
			log.Fatal("Can't create pid(%s): %s", pidfile, err)
			return
		}
		pid.WriteString(fmt.Sprintf("%d", os.Getpid()))
	}

	runService(&c)

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

	job := twitter_job.Twitter_job{
		Config: &c,
		Sendtweet: sendtweet,
		Senddm: senddm,
		Getinfo: getinfo,
		Getfriendship: getfriendship,
	}

	queue := gosque.CreateQueue("", 0, "", "twitter")
	err := queue.Register(&job)
	if err != nil {
		log.Fatal(err)
	}
	queue.Serve(5e9)

	log.Printf("Service stop")
}
