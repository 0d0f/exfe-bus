package twitter_service

import (
	"fmt"
	"oauth"
	"net/url"
	"log"
	"os"
	"bytes"
	"encoding/json"
)

type UsersShowArg struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string
	UpdateId     int64

	UserId     *string
	ScreenName *string

	IdentityId *uint64
}

func (i *UsersShowArg) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("{User:")
	if i.ScreenName != nil {
		buf.WriteString(fmt.Sprintf("%s", *i.ScreenName))
	}
	if i.UserId != nil {
		buf.WriteString(fmt.Sprintf("(%s)", *i.UserId))
	}
	buf.WriteString(" IdentityId:")
	if i.IdentityId != nil {
		buf.WriteString(fmt.Sprintf("(%d)", *i.IdentityId))
	}
	buf.WriteString(fmt.Sprintf(" Client:(%s %s) Access:(%s %s)}",
		i.ClientToken, i.ClientSecret, i.AccessToken, i.AccessSecret))
	return buf.String()
}

func (arg *UsersShowArg) getValues() (v url.Values, err error) {
	if (arg.ScreenName == nil) && (arg.UserId == nil) {
		return nil, fmt.Errorf("ScreenName and UserId in arg should not both be empty.")
	}

	v = make(url.Values)
	if arg.ScreenName != nil {
		v.Add("screen_name", *arg.ScreenName)
	} else {
		v.Add("user_id", *arg.UserId)
	}
	return v, nil
}

type Users struct {
	UpdateInfoService
	log *log.Logger
}

func NewUsers(site_api string) *Users {
	log := log.New(os.Stderr, "exfe.twitter.users", log.LstdFlags)
	return &Users{
		UpdateInfoService: UpdateInfoService{
			SiteApi: site_api,
		},
		log: log,
	}
}

func (s *Users) GetInfo(arg *UsersShowArg, reply *UserInfo) error {
	s.log.Printf("show: %s", arg)

	client := oauth.CreateClient(arg.ClientToken, arg.ClientSecret, arg.AccessToken, arg.AccessSecret, "https://api.twitter.com/1/")

	params, err := arg.getValues()
	if err != nil {
		s.log.Printf("Can't get arg's value: %s", err)
		return err
	}

	retReader, err := client.Do("GET", "/users/show.json", params)
	if err != nil {
		s.log.Printf("Twitter access error: %s", err)
		return err
	}

	decoder := json.NewDecoder(retReader)
	err = decoder.Decode(reply)
	if err != nil {
		s.log.Printf("Can't parse twitter reply: %s", err)
		return err
	}

	if arg.IdentityId != nil {
		go func() {
			id := *arg.IdentityId
			err := s.UpdateUserInfo(id, reply, 4)
			if err != nil {
				s.log.Printf("Update identity(%d) info fail: %s", id, err)
			} else {
				s.log.Printf("Update identity(%d) info succeed", id)
			}
		}()
	}

	return nil
}
