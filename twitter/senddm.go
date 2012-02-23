package main

import (
	"oauth"
	"gosque"
	"net/url"
	"fmt"
)

type Message struct {
	ClientToken string
	ClientSecret string
	AccessToken string
	AccessSecret string
	Message string
	ToUserName string
	ToUserId string
}

func (m *Message) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) ToUser:%s(%s) Message:%s}",
		m.ClientToken, m.ClientSecret, m.AccessToken, m.AccessSecret, m.ToUserName, m.ToUserId, m.Message)
}

const (
	queue = "resque:twitter:directmessage"
	maxJobs = 10
	timeOut = 5e9 // 5 seconds
)

func generateEmptyMessage() interface{} {
	return &Message{}
}

func sendDirectMessage(message *Message) error {
	client := oauth.CreateClient(message.ClientToken, message.ClientSecret, message.AccessToken, message.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if message.ToUserId != "" {
		params.Add("user_id", message.ToUserId)
	} else {
		params.Add("screen_name", message.ToUserName)
	}
	params.Add("text", message.Message)
	_, err := client.Do("POST", "/direct_messages/new.json", params)
	return err
}

func main() {
	gosque := gosque.CreateQueue("", 0, "", queue)
	defer func() { gosque.Close() }()

	jobRecv := gosque.IncomingJob(generateEmptyMessage, timeOut)
	for {
		job := (<-jobRecv).(*Message)
		fmt.Println("Process job: ", job.Message)

		err := sendDirectMessage(job)
		if err != nil {
			fmt.Printf("Send message (%s) failed: %s\n", job, err)
		}
	}
}
