package main

import (
	"config"
	"flag"
	"gobus"
	"io/ioutil"
	"log"
	"net/url"
	"oauth"
	"time"
	"./pkg/twitter"
)

type UserInfoService struct {
}

func (i *UserInfoService) Do(arg twitter.UserInfo, reply *string) error {
	log.Printf("Try to get %s(%s) userinfo...", arg.ScreenName, arg.UserId)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if arg.ScreenName != "" {
		params.Add("screen_name", arg.ScreenName)
	} else {
		params.Add("user_id", arg.UserId)
	}
	retReader, err := client.Do("GET", "/users/show.json", params)
	if err != nil {
		log.Printf("Twitter access error: %s", err)
		return err
	}

	retBytes, err := ioutil.ReadAll(retReader)
	if err != nil {
		log.Printf("Can't load twitter response: %s", err)
		return err
	}

	// TODO:
	// twitter info update
	*reply = string(retBytes)
	return nil
}

const (
	queue = "twitter:userinfo"
)

func main() {
	log.SetPrefix("[UserInfo]")
	log.Printf("Service start, queue: %s", queue)

	var configFile string
	flag.StringVar(&configFile, "config", "twitter_sender.yaml", "Specify the configuration file")
	flag.Parse()

	config := config.LoadFile(configFile)

	service := gobus.CreateService(
		config.String("redis.netaddr"),
		config.Int("redis.db"),
		config.String("redis.password"),
		queue,
		&UserInfoService{})
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Serve(time.Duration(config.Int("service.time_out")))
}
