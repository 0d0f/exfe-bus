package main

import (
	"twitter/service"
	"exfe/service"
	"config"
	"gobus"
	"log/syslog"
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

func main() {
	log, err := syslog.New(syslog.LOG_INFO, "exfe.twitter")
	if err != nil {
		panic(err)
	}
	log.Info("Service start")

	var c exfe_service.Config

	var pidfile string
	var configFile string

	flag.StringVar(&pidfile, "pid", "", "Specify the pid file")
	flag.StringVar(&configFile, "config", "exfe.json", "Specify the configuration file")
	flag.Parse()

	config.LoadFile(configFile, &c)

	flag.Parse()
	if pidfile != "" {
		pid, err := os.Create(pidfile)
		if err != nil {
			log.Crit(fmt.Sprintf("Can't create pid(%s): %s", pidfile, err))
			return
		}
		pid.WriteString(fmt.Sprintf("%d", os.Getpid()))
	}

	server := gobus.CreateServer(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password, "twitter")

	server.Register(new(twitter_service.FriendshipsExists))

	user := new(twitter_service.Users)
	user.SiteApi = c.Site_api
	server.Register(user)

	server.Register(new(twitter_service.Statuses))

	d := new(twitter_service.DirectMessages)
	d.SiteApi = c.Site_api
	server.Register(d)

	server.Serve(c.Twitter.Time_out * 1e9)
}
