package facebook

import (
	"broker"
	"encoding/json"
	"fmt"
	"logger"
	"model"
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

func (f *Facebook) Post(id, text string) (string, error) {
	return f.helper.SendEmail(id, text)
}

func (f *Facebook) UpdateFriends(to *model.Recipient) error {
	idToken, err := f.getToken(to)
	if err != nil {
		return fmt.Errorf("can't convert %s's AuthData(%s): %s", to, to.AuthData, err)
	}
	url := fmt.Sprintf("https://graph.facebook.com/%s/friends?access_token=%s", to.ExternalID, idToken.Token)
	for {
		resp, err := broker.HttpResponse(broker.Http("GET", url, "", nil))
		if err != nil {
			return fmt.Errorf("facebook get friends from %s error: %s", url, err)
		}
		defer resp.Close()
		var friends facebookFriendsReply
		decoder := json.NewDecoder(resp)
		err = decoder.Decode(&friends)
		if err != nil {
			return fmt.Errorf("facebook get friends json error: %s", err)
		}
		if len(friends.Data) == 0 {
			break
		}
		users := make([]thirdpart.ExternalUser, 0)
		c := make(chan *facebookUser)
		fmt.Sprintf("facebook bust:", len(friends.Data))
		for _, friend := range friends.Data {
			go func(id string) {
				user, err := f.getInfo(idToken, id)
				defer func() {
					c <- user
				}()
				if err != nil {
					logger.ERROR("can't get %s facebook infomation: %s", id, err)
					return
				}
				if user.ExternalUsername() == "" {
					logger.ERROR("facebook user %s doesn't have username, ignored", id)
					user = nil
					return
				}
			}(friend.Id)
		}
		for _ = range friends.Data {
			user := <-c
			if user != nil {
				users = append(users, user)
			}
		}
		err = f.helper.UpdateFriends(to, users)
		if err != nil {
			return fmt.Errorf("update %s friends error: %s", to, err)
		}
		url = friends.Paging.Next
		if url == "" {
			break
		}
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
	resp, err := broker.HttpResponse(broker.Http("GET", url, "", nil))
	if err != nil {
		return nil, fmt.Errorf("facebook get %s info from %s error: %s", id, url, err)
	}
	defer resp.Close()
	var user facebookUser
	decoder := json.NewDecoder(resp)
	err = decoder.Decode(&user)
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
