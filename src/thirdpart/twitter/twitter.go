package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"model"
	"net/url"
	"oauth"
	"strings"
	"thirdpart"
)

type Twitter struct {
	client      *oauth.OAuthClient
	clientToken *thirdpart.Token
	accessToken *thirdpart.Token
	helper      thirdpart.Helper
}

const twitterApiBase = "https://api.twitter.com/1.1/"
const provider = "twitter"

func New(client, access *thirdpart.Token, helper thirdpart.Helper) *Twitter {
	return &Twitter{
		client:      oauth.CreateClient(client.Token, client.Secret, access.Token, access.Secret, twitterApiBase),
		clientToken: client,
		accessToken: access,
		helper:      helper,
	}
}

func (t *Twitter) Provider() string {
	return provider
}

func (t *Twitter) MessageType() thirdpart.MessageType {
	return thirdpart.ShortMessage
}

type twitterReply struct {
	IDstr string `json:"id_str"`
}

func (t *Twitter) Send(to *model.Recipient, privateMessage string, publicMessage string, data *thirdpart.InfoData) (string, error) {
	params := make(url.Values)
	params.Set(t.identity(to))
	var err error
	var resp io.ReadCloser
	params.Set("text", privateMessage)
	resp, err = t.client.Do("POST", "direct_messages/new.json", params)
	if err != nil && strings.Index(err.Error(), `"code":150`) > 0 {
		params := make(url.Values)
		params.Set("status", publicMessage)
		resp, err = t.client.Do("POST", "statuses/update.json", params)
	}
	if err != nil {
		return "", fmt.Errorf("send to %s fail: %s", to, err)
	}
	decoder := json.NewDecoder(resp)
	var reply twitterReply
	err = decoder.Decode(&reply)
	if err != nil {
		return "", fmt.Errorf("parse %s reply error: %s", to, err)
	}
	return reply.IDstr, nil
}

func (t *Twitter) UpdateIdentity(to *model.Recipient) error {
	params := make(url.Values)
	params.Set(t.identity(to))
	resp, err := t.client.Do("GET", "users/show.json", params)
	if err != nil {
		return fmt.Errorf("get %s users/show(%v) failed: %s", to, params, err)
	}
	var info twitterInfo
	decoder := json.NewDecoder(resp)
	err = decoder.Decode(&info)
	if err != nil {
		return fmt.Errorf("can't parse %s users/show(%v) reply: %s", to, params, err)
	}
	err = t.helper.UpdateIdentity(to, info)
	if err != nil {
		return fmt.Errorf("update %s error: %s", to, err)
	}
	return nil
}

func (t *Twitter) UpdateFriends(to *model.Recipient) error {
	var idToken twitterIdentityToken
	err := json.Unmarshal([]byte(to.AuthData), &idToken)
	if err != nil {
		return fmt.Errorf("can't convert %s's AuthData: %s", to, err)
	}
	access := idToken.ToToken()
	client := oauth.CreateClient(t.clientToken.Token, t.clientToken.Secret, access.Token, access.Secret, twitterApiBase)

	params := make(url.Values)
	params.Set(t.identity(to))
	resp, err := client.Do("GET", "friends/ids.json", params)
	if err != nil {
		return fmt.Errorf("get %s friends/ids(%v) failed: %s", to, params, err)
	}
	var twitterIDs_ twitterIDs
	decoder := json.NewDecoder(resp)
	err = decoder.Decode(&twitterIDs_)
	if err != nil {
		return fmt.Errorf("parse %s friends/ids(%s) reply failed: %s", to, params, err)
	}

	friendIDs := twitterIDs_.IDs
	for len(friendIDs) > 0 {
		ids := friendIDs
		if len(friendIDs) > 100 {
			ids = friendIDs[:100]
			friendIDs = friendIDs[100:]
		}

		params := make(url.Values)
		params.Set("user_id", join(ids, ","))
		resp, err := client.Do("GET", "users/lookup.json", params)
		if err != nil {
			return fmt.Errorf("get %s users/lookup.json(%v) fail: %s", to, params, err)
		}
		var users []twitterInfo
		decoder := json.NewDecoder(resp)
		err = decoder.Decode(&users)
		if err != nil {
			return fmt.Errorf("parse %s users/lookup(%v) reply failed: %s", to, params, err)
		}

		users_ := make([]thirdpart.ExternalUser, len(users))
		for i, u := range users {
			users_[i] = u
		}
		err = t.helper.UpdateFriends(to, users_)
		if err != nil {
			return fmt.Errorf("update %s's friends fail: %s", to, err)
		}
	}
	return nil
}

func (t *Twitter) identity(id *model.Recipient) (key, value string) {
	if id.ExternalID != "" {
		return "user_id", id.ExternalID
	}
	return "screen_name", id.ExternalUsername
}

type twitterIdentityToken struct {
	Token  string `json:"oauth_token"`
	Secret string `json:"oauth_token_secret"`
}

func (t twitterIdentityToken) ToToken() *thirdpart.Token {
	return &thirdpart.Token{
		Token:  t.Token,
		Secret: t.Secret,
	}
}

type twitterIDs struct {
	Previous string   `json:"previous_cursor_str"`
	Next     string   `json:"next_cursor_str"`
	IDs      []uint64 `json:"ids"`
}

type twitterInfo struct {
	ID              uint64 `json:"id"`
	ScreenName      string `json:"screen_name"`
	ProfileImageUrl string `json:"profile_image_url"`
	Description     string `json:"description"`
	Name_           string `json:"name"`
}

func (i twitterInfo) ExternalID() string {
	return fmt.Sprintf("%d", i.ID)
}

func (i twitterInfo) Provider() string {
	return provider
}

func (i twitterInfo) ExternalUsername() string {
	return i.ScreenName
}

func (i twitterInfo) Name() string {
	return i.Name_
}

func (i twitterInfo) Bio() string {
	return i.Description
}

func (i twitterInfo) Avatar() string {
	return i.ProfileImageUrl
}

func join(a []uint64, spliter string) string {
	if len(a) == 0 {
		return ""
	}
	ret := fmt.Sprintf("%d", a[0])
	for _, i := range a[1:] {
		ret = fmt.Sprintf("%s,%d", ret, i)
	}
	return ret
}
