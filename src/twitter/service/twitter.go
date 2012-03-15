package twitter_service

import (
	"fmt"
	"net/url"
)

type TwitterUserInfo struct {
	Id int64
	Screen_name string
	Profile_image_url string
	Description string
	Name string
}

func (i *TwitterUserInfo) MakeUrlValues(id uint64) (v url.Values) {
	v.Add("id", fmt.Sprintf("%d", id))
	v.Add("provider", "twitter")
	v.Add("external_identity", fmt.Sprintf("%d", i.Id))
	if i.Name != "" {
		v.Add("name", i.Name)
	}
	if i.Description != "" {
		v.Add("bio", i.Description)
	}
	if i.Profile_image_url != "" {
		v.Add("avatar_url", i.Profile_image_url)
	}
	if i.Screen_name != "" {
		v.Add("external_username", i.Screen_name)
	}
	return
}
