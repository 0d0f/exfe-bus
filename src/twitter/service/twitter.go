package twitter_service

import (
	"fmt"
	"net/url"
	"net/http"
)

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
	v.Add("external_identity", fmt.Sprintf("%d", i.Id))
	if i.Name != nil {
		v.Add("name", *i.Name)
	}
	if i.Description != nil {
		v.Add("bio", *i.Description)
	}
	if i.Profile_image_url != nil {
		v.Add("avatar_url", *i.Profile_image_url)
	}
	if i.Screen_name != nil {
		v.Add("external_username", *i.Screen_name)
	}
	return
}

type UpdateInfoService struct {
	SiteUrl string
}

func (s *UpdateInfoService) UpdateUserInfo(id uint64, i *UserInfo) error {
	url := fmt.Sprintf("%s/identity/update", s.SiteUrl)
	resp, err := http.PostForm(url, i.makeUrlValues(id))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Update to %s failed: %s", url, resp.Status)
	}
	return nil
}