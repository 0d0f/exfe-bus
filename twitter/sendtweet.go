package main

import (
	"./pkg/twitter"
	"config"
	"flag"
	"gobus"
	"io/ioutil"
	"log"
	"net/url"
	"oauth"
	"time"
)

type TweetService struct {
}

func (t *TweetService) Do(arg twitter.Tweet, reply *string) error {
	log.Printf("Try to send tweet(%s)...", arg.Tweet)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	params.Add("status", arg.Tweet)

	retReader, err := client.Do("POST", "/statuses/update.json", params)
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
	queue = "twitter:tweet"
)

func main() {
	log.SetPrefix("[Tweet]")
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
		&TweetService{})
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Serve(time.Duration(config.Int("service.time_out")))
}
