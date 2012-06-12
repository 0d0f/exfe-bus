package twitter_service

import (
	"fmt"
	"net/url"
	"net/http"
)

const UrlLength = 20
const MaxTweetLength = 140

type UserInfo struct {
	Id uint64
	Screen_name *string
	Profile_image_url *string
	Description *string
	Name *string
}

func (i *UserInfo) makeUrlValues(id uint64) (v url.Values) {
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

func (s *UpdateInfoService) UpdateUserInfo(id uint64, i *UserInfo, _ int) error {
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

func makeText(message string, urls []string) string {
	maxMessageLength := MaxTweetLength - (len(urls) + 1/* space */) * UrlLength

	if len(message) > maxMessageLength {
		message = fmt.Sprintf("%s…", message[0:(maxMessageLength-1)])
	}
	for _, url := range urls {
		message = fmt.Sprintf("%s %s", message, url)
	}

	return message
}
