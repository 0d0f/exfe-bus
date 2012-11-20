package facebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"model"
	"net/http"
	"thirdpart"
)

type Facebook struct {
	helper thirdpart.Helper
}

const provider = "facebook"

func New(helper thirdpart.Helper) *Facebook {
	return &Facebook{
		helper: helper,
	}
}

func (f *Facebook) Provider() string {
	return provider
}

func (f *Facebook) Send(to *model.Recipient, privateMessage string, publicMessage string, info *model.InfoData) (string, error) {
	return f.helper.SendEmail(fmt.Sprintf("%s@facebook.com", to.ExternalUsername), privateMessage)
}

func (f *Facebook) UpdateFriends(to *model.Recipient) error {
	idToken, err := f.getToken(to)
	if err != nil {
		return fmt.Errorf("can't convert %s's AuthData(%s): %s", to, to.AuthData, err)
	}
	url := fmt.Sprintf("https://graph.facebook.com/%s/friends?access_token=%s", to.ExternalID, idToken.Token)
	for {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("facebook get friends from %s error: %s", url, err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("facebook get body from %s fail: %s", url, err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("facebook get friends from %s fail: (%s) %s", url, resp.Status, string(body))
		}
		var friends facebookFriendsReply
		err = json.Unmarshal(body, &friends)
		if err != nil {
			return fmt.Errorf("facebook get friends json error: %s", err)
		}
		if len(friends.Data) == 0 {
			break
		}
		users := make([]thirdpart.ExternalUser, 0)
		for _, friend := range friends.Data {
			user, err := f.getInfo(idToken, friend.Id)
			if err != nil {
				f.helper.Log().Err("can't get %s facebook infomation: %s", friend.Id, err)
				continue
			}
			if user.ExternalUsername() == "" {
				f.helper.Log().Err("facebook user %d doesn't have username, ignored", friend.Id)
				continue
			}
			users = append(users, user)
		}
		err = f.helper.UpdateFriends(to, users)
		if err != nil {
			return fmt.Errorf("update %s friends error: %s", to, err)
		}
		url = friends.Paging.Next
	}
	return nil
}

func (f *Facebook) UpdateIdentity(to *model.Recipient) error {
	idToken, err := f.getToken(to)
	if err != nil {
		return fmt.Errorf("can't convert %s's AuthData(%s): %s", to, to.AuthData, err)
	}
	user, err := f.getInfo(idToken, "me")
	if err != nil {
		return err
	}
	err = f.helper.UpdateIdentity(to, user)
	if err != nil {
		return fmt.Errorf("update %s info error: %s", to, err)
	}
	return nil
}

func (f Facebook) getInfo(idToken *facebookIdentityToken, id string) (*facebookUser, error) {
	url := fmt.Sprintf("https://graph.facebook.com/%s?access_token=%s", id, idToken.Token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("facebook get %s info from %s error: %s", id, url, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("facebook get %s info body from %s fail: %s", id, url, err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("facebook get %s info from %s fail: (%s) %s", id, url, resp.Status, string(body))
	}
	var user facebookUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, fmt.Errorf("facebook get %s info json error: %s", id, err)
	}
	return &user, nil
}

func (f *Facebook) getToken(to *model.Recipient) (*facebookIdentityToken, error) {
	var idToken facebookIdentityToken
	err := json.Unmarshal([]byte(to.AuthData), &idToken)
	if err != nil {
		return nil, err
	}
	if idToken.Token == "" {
		return nil, fmt.Errorf("can't find token info")
	}
	return &idToken, nil
}

type facebookIdentityToken struct {
	Token string `json:"oauth_token"`
}

type facebookPaging struct {
	Next string `json:"next"`
}

type facebookFriendsReply struct {
	Data   []facebookUser `json:"data"`
	Paging facebookPaging `json:"paging"`
}

type facebookUser struct {
	Id       string `json:"id"`
	Name_    string `json:"name"`
	Username string `json:"username"`
	Link     string `json:"link"`
}

func (f facebookUser) ExternalID() string {
	return f.Id
}

func (f facebookUser) Provider() string {
	return provider
}

func (f facebookUser) ExternalUsername() string {
	return f.Username
}

func (f facebookUser) Name() string {
	return f.Name_
}

func (f facebookUser) Bio() string {
	return f.Link
}

func (f facebookUser) Avatar() string {
	return fmt.Sprintf("http://graph.facebook.com/%s/picture", f.Id)
}
