package twitter_service

import (
	"fmt"
	"oauth"
	"net/url"
	"log"
	"encoding/json"
)

type UsersShowArg struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string

	UserId     string
	ScreenName string
}

func (i *UsersShowArg) String() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) User %s(%s)}",
		i.ClientToken, i.ClientSecret, i.AccessToken, i.AccessSecret, i.ScreenName, i.UserId)
}

type UsersShowReply TwitterUserInfo

type UsersShow struct {
}

func (i *UsersShow) Do(arg *UsersShowArg, reply *UsersShowReply) error {
	log.Printf("[Info][users/show]Call by arg %s", arg)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if arg.ScreenName != "" {
		params.Add("screen_name", arg.ScreenName)
	} else {
		params.Add("user_id", arg.UserId)
	}
	retReader, err := client.Do("GET", "/users/show.json", params)
	if err != nil {
		log.Printf("[Error][users/show]Twitter access error: %s", err)
		return err
	}

	// TODO:
	// twitter info update

	decoder := json.NewDecoder(retReader)
	err = decoder.Decode(reply)
	if err != nil {
		// some user will not fill all field, and twitter responses of these fields  are null,
		// which will cause decode error
		// log.Printf("[Error][users/show]Can't parse twitter reply: %s", err)
		// return err
	}
	return nil
}
