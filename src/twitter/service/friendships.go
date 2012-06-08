package twitter_service

import (
	"fmt"
	"oauth"
	"net/url"
	"log"
	"os"
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

type Friendships struct {
	log *log.Logger
}

func NewFriendships() *Friendships {
	log := log.New(os.Stderr, "exfe.twitter.friendships", log.LstdFlags)
	return &Friendships{
		log: log,
	}
}

func (f *Friendships) GetFriendship(arg *FriendshipsExistsArg, reply *FriendshipsExistsReply) error {
	f.log.Printf("exists: %s", arg)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	params.Add("user_a", arg.UserA)
	params.Add("user_b", arg.UserB)
	retReader, err := client.Do("GET", "/friendships/exists.json", params)
	if err != nil {
		f.log.Printf("Twitter access error: %s", err)
		return err
	}

	decoder := json.NewDecoder(retReader)
	err = decoder.Decode(reply)
	if err != nil {
		f.log.Printf("Parse twitter response error: %s", err)
		return err
	}

	return nil
}
