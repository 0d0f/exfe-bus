package twitter_service

import (
	"fmt"
	"oauth"
	"net/url"
	"log"
	"encoding/json"
)

type FriendshipsExistsArg struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string

	UserA string
	UserB string
}

func (f *FriendshipsExistsArg) String() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) UserA:%s UserB:%s}",
		f.ClientToken, f.ClientSecret, f.AccessToken, f.AccessSecret,
		f.UserA, f.UserB)
}

type FriendshipsExistsReply bool

type FriendshipsExists struct {
}

func (f *FriendshipsExists) GetFriendship(arg *FriendshipsExistsArg, reply *FriendshipsExistsReply) error {
	log.Printf("[Info][friendships/exists]Call by arg: %s", arg)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	params.Add("user_a", arg.UserA)
	params.Add("user_b", arg.UserB)
	retReader, err := client.Do("GET", "/friendships/exists.json", params)
	if err != nil {
		log.Printf("[Error][friendships/exists]Twitter access error: %s", err)
		return err
	}

	decoder := json.NewDecoder(retReader)
	err = decoder.Decode(reply)
	if err != nil {
		log.Printf("[Error][friendships/exists]Parse twitter response error: %s", err)
		return err
	}

	return nil
}
