package main

import (
	"twitter"
	"config"
	"flag"
	"gobus"
	"io/ioutil"
	"log"
	"net/url"
	"oauth"
	"time"
)

type MessageService struct {
}

func (m *MessageService) Do(arg *twitter.DirectMessage, reply *string) error {
	log.Printf("Try to send dm(%s) to user(%s/%s)...", arg.Message, arg.ToUserName, arg.ToUserId)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if arg.ToUserId != "" {
		params.Add("user_id", arg.ToUserId)
	} else {
		params.Add("screen_name", arg.ToUserName)
	}
	params.Add("text", arg.Message)
	retReader, err := client.Do("POST", "/direct_messages/new.json", params)
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
	queue = "twitter:directmessage"
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
		&MessageService{})
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Serve(time.Duration(config.Int("service.time_out")))
}
