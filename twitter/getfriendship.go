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
	friendship *twitter.Friendship
}

func (f *FriendshipService) Do(messages []interface{}) []interface{} {
	data, ok := messages[0].(*FriendshipService)
	if !ok {
		log.Printf("Can't convert input into Friendship: %s", messages)
		return nil
	}

	log.Printf("Try to get %s/%s friendship...", data.friendship.UserA, data.friendship.UserB)

	client := oauth.CreateClient(data.friendship.ClientToken, data.friendship.ClientSecret, data.friendship.AccessToken, data.friendship.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	params.Add("user_a", data.friendship.UserA)
	params.Add("user_b", data.friendship.UserB)
	retReader, err := client.Do("GET", "/friendships/exists.json", params)
	if err != nil {
		log.Printf("Twitter access error: %s", err)
		return []interface{}{map[string]string{"error": err.Error()}}
	}

	retBytes, err := ioutil.ReadAll(retReader)
	if err != nil {
		log.Printf("Can't load twitter response: %s", err)
		return []interface{}{map[string]string{"error": err.Error()}}
	}

	return []interface{}{twitter.Response{
		Result: string(retBytes),
	}}
}

func (f *FriendshipService) MaxJobsCount() int {
	return 1
}

func (f *FriendshipService) JobGenerator() interface{} {
	return &FriendshipService{}
}

const (
	queue = "gobus:queue:twitter:friendship"
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
		&FriendshipService{},
		config.Int("service.limit"))
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Run(time.Duration(config.Int("service.time_out")))
}
