package main

import (
	"oauth"
	"gosque"
	"net/url"
	"fmt"
)

type Tweet struct {
	ClientToken string
	ClientSecret string
	AccessToken string
	AccessSecret string
	Tweet string
}

func (t *Tweet) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) Tweet:%s}", t.ClientToken, t.ClientSecret, t.AccessToken, t.AccessSecret, t.Tweet)
}

const (
	queue = "resque:twitter:tweet"
	maxJobs = 10
	timeOut = 5e9 // 5 seconds
)

func generateEmptyTweet() interface{} {
	return &Tweet{}
}

func sendTweet(tweet *Tweet) error {
	client := oauth.CreateClient(tweet.ClientToken, tweet.ClientSecret, tweet.AccessToken, tweet.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	params.Add("status", tweet.Tweet)
	_, err := client.Do("POST", "/statuses/update.json", params)
	return err
}

func main() {
	gosque := gosque.CreateQueue("", 0, "", queue)
	defer func() { gosque.Close() }()

	jobRecv := gosque.IncomingJob(generateEmptyTweet, timeOut)
	for {
		job := (<-jobRecv).(*Tweet)
		fmt.Println("Process job: ", job.Tweet)

		err := sendTweet(job)
		if err != nil {
			fmt.Printf("Send tweet (%s) failed: %s\n", job, err)
		}
	}
}
