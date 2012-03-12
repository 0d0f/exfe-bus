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

type TweetService struct {
	tweet twitter.Tweet
}

func (t *TweetService) Do(jobs []interface{}) []interface{} {
	data, ok := jobs[0].(*TweetService)
	if !ok {
		log.Println("Can't convert input into TweetService: %s", jobs)
		return nil
	}

	log.Printf("Try to send tweet(%s)...", data.tweet.Tweet)

	client := oauth.CreateClient(data.tweet.ClientToken, data.tweet.ClientSecret, data.tweet.AccessToken, data.tweet.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	params.Add("status", data.tweet.Tweet)

	retReader, err := client.Do("POST", "/statuses/update.json", params)
	if err != nil {
		log.Printf("Twitter access error: %s", err)
		return []interface{}{map[string]string{"result": "", "error": err.Error()}}
	}

	retBytes, err := ioutil.ReadAll(retReader)
	if err != nil {
		log.Printf("Can't load twitter response: %s", err)
		return []interface{}{map[string]string{"result": "", "error": err.Error()}}
	}

	return []interface{}{twitter.Response{
		Error: "",
		Result: string(retBytes),
	}}
}

func (t *TweetService) MaxJobsCount() int {
	return 1
}

func (t *TweetService) JobGenerator() interface{} {
	return &TweetService{}
}

const (
	queue = "gobus:queue:twitter:tweet"
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
		&TweetService{},
		config.Int("service.limit"))
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Run(time.Duration(config.Int("service.time_out")))
}
