package main

import (
	"fmt"
	"gobus"
	"io/ioutil"
	"log"
	"net/url"
	"oauth"
)

type Friendship struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string

	UserA string
	UserB string
}

func (f *Friendship) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) UserA %s UserB %s}",
		f.ClientToken, f.ClientSecret, f.AccessToken, f.AccessSecret,
		f.UserA, f.UserB)
}

func (f *Friendship) Do(messages []interface{}) []interface{} {
	message, ok := messages[0].(*Friendship)
	if !ok {
		log.Printf("Can't convert input into Friendship: %s", messages)
		return nil
	}

	log.Printf("Try to get %s/%s friendship...", message.UserA, message.UserB)

	client := oauth.CreateClient(message.ClientToken, message.ClientSecret, message.AccessToken, message.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	params.Add("user_a", message.UserA)
	params.Add("user_b", message.UserB)
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

	return []interface{}{map[string]string{"result": string(retBytes)}}
}

func (f *Friendship) MaxJobsCount() int {
	return 1
}

func (f *Friendship) JobGenerator() interface{} {
	return &Friendship{}
}

const (
	queue   = "gobus:queue:twitter:friendship"
	timeOut = 5e9 // 5 seconds
	limit   = 10
)

func main() {
	log.SetPrefix("[Friendship]")
	log.Printf("Service start, queue: %s", queue)

	service := gobus.CreateService("", 0, "", queue, &Friendship{}, limit)
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Run(timeOut)
}
