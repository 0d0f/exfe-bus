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

func (i *TwitterUserInfo) MakeUrlValues(id int64) (v url.Values) {
	v.Add("id", fmt.Sprintf("%d", id))
	v.Add("provider", "twitter")
	v.Add("external_identity", fmt.Sprintf("%d", i.Id))
	v.Add("name", i.Name)
	v.Add("bio", i.Description)
	v.Add("avatar_url", i.Profile_image_url)
	v.Add("external_username", i.Screen_name)
	return
}
