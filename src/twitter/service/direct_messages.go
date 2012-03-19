package twitter_service

import (
	"net/url"
	"fmt"
	"oauth"
	"log"
	"bytes"
	"encoding/json"
)

type DirectMessagesNewArg struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string

	Message      string
	ToUserName   *string
	ToUserId     *string

	IdentityId *uint64
}

func (m *DirectMessagesNewArg) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("{User:")
	if m.ToUserName != nil {
		buf.WriteString(fmt.Sprintf("%s", *m.ToUserName))
	}
	if m.ToUserId != nil {
		buf.WriteString(fmt.Sprintf("(%s)", *m.ToUserId))
	}

	buf.WriteString(fmt.Sprintf(" Client:(%s %s) Access:(%s %s) Msg:%s}",
		m.ClientToken, m.ClientSecret, m.AccessToken, m.AccessSecret, m.Message))
	return buf.String()
}

func (m *DirectMessagesNewArg) getValues() (v url.Values, err error) {
	if (m.ToUserName == nil) && (m.ToUserId == nil) {
		return nil, fmt.Errorf("ScreenName and UserId in arg should not both be empty.")
	}

	v = make(url.Values)
	if m.ToUserId != nil {
		v.Add("user_id", *m.ToUserId)
	} else {
		v.Add("screen_name", *m.ToUserName)
	}
	v.Add("text", m.Message)

	return v, nil
}

type DirectMessagesNewReply struct {
	Sender UserInfo
	Recipient UserInfo
}

type DirectMessagesNew struct {
	UpdateInfoService
}

func (m *DirectMessagesNew) Do(arg *DirectMessagesNewArg, reply *DirectMessagesNewReply) error {
	log.Printf("[Info][direct_messages/new]Call by arg: %s", arg)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")

	params, err := arg.getValues()
	if err != nil {
		log.Printf("[Error][users/shwo]Can't get arg's value: %s", err)
		return err
	}

	retReader, err := client.Do("POST", "/direct_messages/new.json", params)
	if err != nil {
		log.Printf("[Error][direct_messages/new]Twitter access error: %s", err)
		return err
	}

	decoder := json.NewDecoder(retReader)
	err = decoder.Decode(reply)
	if err != nil {
		log.Printf("[Error][direct_messages/new]Parse twitter response error: %s", err)
		return err
	}

	if arg.IdentityId != nil {
		go func() {
			id := *arg.IdentityId
			err := m.UpdateUserInfo(id, &reply.Recipient)
			if err != nil {
				log.Printf("[Error][direct_messages/new]Update identity(%d) info fail: %s", id, err)
			} else {
				log.Printf("[Info][direct_messages/new]Update identity(%d) info success", id)
			}
		}()
	}

	return nil
}
