package main

import (
	"oauth"
	"gobus"
	"net/url"
	"fmt"
	"io/ioutil"
)

type TweetService struct {
	ClientToken string
	ClientSecret string
	AccessToken string
	AccessSecret string
	Tweet string
}

func (t *TweetService) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) Tweet:%s}", t.ClientToken, t.ClientSecret, t.AccessToken, t.AccessSecret, t.Tweet)
}

func (t *TweetService) Do(jobs []interface{}) []interface{} {
	tweet := jobs[0].(*TweetService)
	client := oauth.CreateClient(tweet.ClientToken, tweet.ClientSecret, tweet.AccessToken, tweet.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	params.Add("status", tweet.Tweet)

	retReader, err := client.Do("POST", "/statuses/update.json", params)
	if err != nil {
		return []interface{}{map[string]string{"result": "", "error": err.Error()}, }
	}

	retBytes, err := ioutil.ReadAll(retReader)
	if err != nil {
		return []interface{}{map[string]string{"result": "", "error": err.Error()}, }
	}

	return []interface{}{map[string]string{"result": string(retBytes), "error": ""}, }
}

func (t *TweetService) MaxJobsCount() int {
	return 1
}

func (t *TweetService) JobGenerator() interface{} {
	return &TweetService{}
}

const (
	queue = "gobus:queue:twitter:tweet"
	timeOut = 5e9 // 5 seconds
)

func main() {
	service := gobus.CreateService("", 0, "", queue, &TweetService{})
	defer service.Close()

	service.Run(timeOut)
}
