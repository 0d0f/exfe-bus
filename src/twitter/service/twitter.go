package twitter_service

import (
	"fmt"
	"net/http"
	"net/url"
)

const UrlLength = 20
const MaxTweetLength = 140
const MaxTwitterIDLength = 16

type UserInfo struct {
	Id                uint64
	Screen_name       *string
	Profile_image_url *string
	Description       *string
	Name              *string
}

func (i *UserInfo) makeUrlValues(id int64) (v url.Values) {
	v = make(url.Values)
	v.Add("id", fmt.Sprintf("%d", id))
	v.Add("provider", "twitter")
	v.Add("external_id", fmt.Sprintf("%d", i.Id))
	if i.Name != nil {
		v.Add("name", *i.Name)
	}
	if i.Description != nil {
		v.Add("bio", *i.Description)
	}
	if i.Profile_image_url != nil {
		v.Add("avatar_filename", *i.Profile_image_url)
	}
	if i.Screen_name != nil {
		v.Add("external_username", *i.Screen_name)
	}
	return
}

type UpdateInfoService struct {
	SiteApi string
}

func (s *UpdateInfoService) UpdateUserInfo(id int64, i *UserInfo, _ int) error {
	url := fmt.Sprintf("%s/v2/gobus/UpdateIdentity", s.SiteApi)
	resp, err := http.PostForm(url, i.makeUrlValues(id))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Update to %s failed: %s", url, resp.Status)
	}
	return nil
}

func makeText(message string, attachments []string) (string, error) {
	length := 0
	for _, a := range attachments {
		if a[:7] == "http://" || a[:8] == "https://" {
			length += UrlLength
		} else {
			length += len(a)
		}
		length += 1
	}
	if length > MaxTweetLength-MaxTwitterIDLength {
		return "", fmt.Errorf("too much attachments")
	}
	maxMessageLength := MaxTweetLength - length

	if len(message) > maxMessageLength {
		message = message[0:maxMessageLength]
	}
	for _, a := range attachments {
		message = fmt.Sprintf("%s %s", message, a)
	}

	return message, nil
}
