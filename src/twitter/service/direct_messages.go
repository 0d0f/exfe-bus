package twitter_service

import (
	"net/url"
	"fmt"
	"oauth"
	"log"
	"encoding/json"
)

type DirectMessagesNewArg struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string

	Message      string
	ToUserName   string
	ToUserId     string
}

func (m *DirectMessagesNewArg) String() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) ToUser:%s(%s) Message:%s}",
		m.ClientToken, m.ClientSecret, m.AccessToken, m.AccessSecret, m.ToUserName, m.ToUserId, m.Message)
}

type DirectMessagesNewReply struct {
	Sender TwitterUserInfo
	Recipient TwitterUserInfo
}

type DirectMessagesNew struct {
}

func (m *DirectMessagesNew) Do(arg *DirectMessagesNewArg, reply *DirectMessagesNewReply) error {
	log.Printf("[Info][direct_messages/new]Call by arg: %s", arg)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")
	params := make(url.Values)
	if arg.ToUserId != "" {
		params.Add("user_id", arg.ToUserId)
	} else {
		params.Add("screen_name", arg.ToUserName)
	}
	params.Add("text", arg.Message)
	retReader, err := client.Do("POST", "/direct_messages/new.json", params)
	if err != nil {
		log.Printf("[Error][direct_messages/new]Twitter access error: %s", err)
		return err
	}

	decoder := json.NewDecoder(retReader)
	err = decoder.Decode(reply)
	if err != nil {
		// PUZZLE: if decode to reply directly, it will return "json: cannot unmarshal null into Go value of type string" error with right reply
		// log.Printf("[Error][direct_messages/new]Parse twitter response error: %s", err)
		// return err
	}

	return nil
}
