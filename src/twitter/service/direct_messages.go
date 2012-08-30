package twitter_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"oauth"
	"os"
)

type DirectMessagesNewArg struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string

	Message    string
	Urls       []string
	ToUserName *string
	ToUserId   *string

	IdentityId *int64
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
	Sender    UserInfo
	Recipient UserInfo
}

type DirectMessages struct {
	UpdateInfoService
	log *log.Logger
}

func NewDirectMessages(site_api string) *DirectMessages {
	log := log.New(os.Stderr, "exfe.twitter.directmessages", log.LstdFlags)
	return &DirectMessages{
		UpdateInfoService: UpdateInfoService{
			SiteApi: site_api,
		},
		log: log,
	}
}

func (m *DirectMessages) SendDM(arg *DirectMessagesNewArg, reply *DirectMessagesNewReply) error {
	m.log.Printf("new: %s", arg)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")

	arg.Message = makeText(arg.Message, arg.Urls)

	params, err := arg.getValues()
	if err != nil {
		m.log.Printf("Can't get arg's value: %s", err)
		return err
	}

	retReader, err := client.Do("POST", "/direct_messages/new.json", params)
	if err != nil {
		m.log.Printf("Twitter access error: %s", err)
		return err
	}

	decoder := json.NewDecoder(retReader)
	err = decoder.Decode(reply)
	if err != nil {
		m.log.Printf("Parse twitter response error: %s", err)
		return err
	}

	if arg.IdentityId != nil {
		go func() {
			id := *arg.IdentityId
			err := m.UpdateUserInfo(id, &reply.Recipient, 1)
			if err != nil {
				m.log.Printf("Update identity(%d) info fail: %s", id, err)
			} else {
				m.log.Printf("Update identity(%d) info success", id)
			}
		}()
	}

	return nil
}
