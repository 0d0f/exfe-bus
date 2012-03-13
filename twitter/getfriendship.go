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

type FriendshipService struct {
}

func (f *FriendshipService) Do(arg twitter.Friendship, reply *string) error {
	log.Printf("Try to get %s/%s friendship...", arg.UserA, arg.UserB)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	params.Add("user_a", arg.UserA)
	params.Add("user_b", arg.UserB)
	retReader, err := client.Do("GET", "/friendships/exists.json", params)
	if err != nil {
		log.Printf("Twitter access error: %s", err)
		return err
	}

	retBytes, err := ioutil.ReadAll(retReader)
	if err != nil {
		log.Printf("Can't load twitter response: %s", err)
		return err
	}

	*reply = string(retBytes)
	return nil
}

const (
	queue = "twitter:friendship"
)

func main() {
	log.SetPrefix("[Friendship]")
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
		&FriendshipService{})

	service.Serve(time.Duration(config.Int("service.time_out")))
}
