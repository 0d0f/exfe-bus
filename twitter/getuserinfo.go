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
	info *twitter.UserInfo
}

func (i *UserInfoService) Do(messages []interface{}) []interface{} {
	data, ok := messages[0].(*UserInfoService)
	if !ok {
		log.Printf("Can't convert input into UserInfo: %s", messages)
		return nil
	}

	log.Printf("Try to get %s(%s) userinfo...", data.info.ScreenName, data.info.UserId)

	client := oauth.CreateClient(data.info.ClientToken, data.info.ClientSecret, data.info.AccessToken, data.info.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if data.info.ScreenName != "" {
		params.Add("screen_name", data.info.ScreenName)
	} else {
		params.Add("user_id", data.info.UserId)
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

func (i *UserInfoService) MaxJobsCount() int {
	return 1
}

func (i *UserInfoService) JobGenerator() interface{} {
	return &UserInfoService{}
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
		&UserInfoService{},
		config.Int("service.limit"))
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Run(time.Duration(config.Int("service.time_out")))
}
