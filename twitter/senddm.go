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

type MessageService struct {
	message *twitter.DirectMessage
}

func (m *MessageService) Do(messages []interface{}) []interface{} {
	data, ok := messages[0].(*MessageService)
	if !ok {
		log.Printf("Can't convert input into Message: %s", messages)
	}

	log.Printf("Try to send dm(%s) to user(%s/%s)...", data.message.Message, data.message.ToUserName, data.message.ToUserId)

	client := oauth.CreateClient(data.message.ClientToken, data.message.ClientSecret, data.message.AccessToken, data.message.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if data.message.ToUserId != "" {
		params.Add("user_id", data.message.ToUserId)
	} else {
		params.Add("screen_name", data.message.ToUserName)
	}
	params.Add("text", data.message.Message)
	retReader, err := client.Do("POST", "/direct_messages/new.json", params)
	if err != nil {
		log.Printf("Twitter access error: %s", err)
		return []interface{}{map[string]string{"error": err.Error()}}
	}

	retBytes, err := ioutil.ReadAll(retReader)
	if err != nil {
		log.Printf("Can't load twitter response: %s", err)
		return []interface{}{map[string]string{"error": err.Error()}}
	}

	return []interface{}{map[string]string{"result": string(retBytes)}}
}

func (m *MessageService) MaxJobsCount() int {
	return 1
}

func (m *MessageService) JobGenerator() interface{} {
	return &MessageService{}
}

const (
	queue = "gobus:queue:twitter:directmessage"
)

func main() {
	log.SetPrefix("[DirectMessage]")
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
		&MessageService{},
		config.Int("service.limit"))
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Run(time.Duration(config.Int("service.time_out")))
}
