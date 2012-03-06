package main

import (
	"oauth"
	"gobus"
	"net/url"
	"fmt"
	"io/ioutil"
)

type UserInfo struct {
	ClientToken string
	ClientSecret string
	AccessToken string
	AccessSecret string

	UserId string
	ScreenName string
}

func (i *UserInfo) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) User %s(%s)}",
		i.ClientToken, i.ClientSecret, i.AccessToken, i.AccessSecret, i.ScreenName, i.UserId)
}

func (i *UserInfo) Do(messages []interface{}) []interface{} {
	message := messages[0].(*UserInfo)
	client := oauth.CreateClient(message.ClientToken, message.ClientSecret, message.AccessToken, message.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if message.ScreenName != "" {
		params.Add("screen_name", message.ScreenName)
	} else {
		params.Add("user_id", message.UserId)
	}
	retReader, err := client.Do("GET", "/users/show.json", params)
	if err != nil {
		return []interface{}{map[string]string{"error": err.Error()}}
	}

	retBytes, err := ioutil.ReadAll(retReader)
	if err != nil {
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
	timeOut = 5e9 // 5 seconds
	limit = 10
)

func main() {
	service := gobus.CreateService("", 0, "", queue, &UserInfo{}, limit)
	defer func() {
		service.Close()
		service.Clear()
	}()

	service.Run(timeOut)
}
