package main

import (
	"fmt"
	"gobus"
	"io/ioutil"
	"log"
	"net/url"
	"oauth"
)

type Message struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string
	Message      string
	ToUserName   string
	ToUserId     string
}

func (m *Message) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) ToUser:%s(%s) Message:%s}",
		m.ClientToken, m.ClientSecret, m.AccessToken, m.AccessSecret, m.ToUserName, m.ToUserId, m.Message)
}

func (m *Message) Do(messages []interface{}) []interface{} {
	message, ok := messages[0].(*Message)
	if !ok {
		log.Printf("Can't convert input into Message: %s", messages)
	}

	log.Printf("Try to send dm(%s) to user(%s/%s)...", message.Message, message.ToUserName, message.ToUserId)

	client := oauth.CreateClient(message.ClientToken, message.ClientSecret, message.AccessToken, message.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if message.ToUserId != "" {
		params.Add("user_id", message.ToUserId)
	} else {
		params.Add("screen_name", message.ToUserName)
	}
	params.Add("text", message.Message)
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

func (m *Message) MaxJobsCount() int {
	return 1
}

func (m *Message) JobGenerator() interface{} {
	return &Message{}
}

const (
	queue   = "gobus:queue:twitter:directmessage"
	timeOut = 5e9 // 5 seconds
	limit   = 10
)

func main() {
	log.SetPrefix("[DirectMessage]")
	log.Printf("Service start, queue: %s", queue)
	service := gobus.CreateService("", 0, "", queue, &Message{}, limit)
	defer func() {
		log.Printf("Service stop, queue: %s", queue)
		service.Close()
		service.Clear()
	}()

	service.Run(timeOut)
}
