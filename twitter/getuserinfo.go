package main

import (
	"config"
	"flag"
	"fmt"
	"gobus"
	"io/ioutil"
	"log"
	"net/url"
	"oauth"
	"time"
)

type UserInfo struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string

	UserId     string
	ScreenName string
}

func (i *UserInfo) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) User %s(%s)}",
		i.ClientToken, i.ClientSecret, i.AccessToken, i.AccessSecret, i.ScreenName, i.UserId)
}

func (i *UserInfo) Do(messages []interface{}) []interface{} {
	message, ok := messages[0].(*UserInfo)
	if !ok {
		log.Printf("Can't convert input into UserInfo: %s", messages)
		return nil
	}

	log.Printf("Try to get %s(%s) userinfo...", message.ScreenName, message.UserId)

	client := oauth.CreateClient(message.ClientToken, message.ClientSecret, message.AccessToken, message.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if message.ScreenName != "" {
		params.Add("screen_name", message.ScreenName)
	} else {
		params.Add("user_id", message.UserId)
	}
	retReader, err := client.Do("GET", "/users/show.json", params)
	if err != nil {
		log.Printf("Twitter access error: %s", err)
		return []interface{}{map[string]string{"error": err.Error()}}
	}

	retBytes, err := ioutil.ReadAll(retReader)
	if err != nil {
		log.Printf("Can't load twitter response: %s", err)
		return []interface{}{map[string]string{"error": err.Error()}}
	}

	// TODO:
	// twitter info update
	return []interface{}{map[string]string{"result": string(retBytes)}}
}

func (i *UserInfo) MaxJobsCount() int {
	return 1
}

func (i *UserInfo) JobGenerator() interface{} {
	return &UserInfo{}
}

const (
	queue = "gobus:queue:twitter:userinfo"
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
		&UserInfo{},
		config.Int("service.limit"))
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Run(time.Duration(config.Int("service.time_out")))
}
