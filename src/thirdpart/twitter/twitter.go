package twitter

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/mrjones/oauth"
	"logger"
	"model"
	"strings"
	"thirdpart"
)

const twitterApiBase = "https://api.twitter.com/1.1/"

var provider = oauth.ServiceProvider{
	RequestTokenUrl:   "http://api.twitter.com/oauth/request_token",
	AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
	AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
}

type Twitter struct {
	token  *oauth.AccessToken
	oauth  broker.OAuth
	helper thirdpart.Helper
	config *model.Config
}

func New(config *model.Config, helper thirdpart.Helper) *Twitter {
	return &Twitter{
		oauth: broker.NewOAuth(config.Thirdpart.Twitter.ClientToken, config.Thirdpart.Twitter.ClientSecret, provider),
		token: &oauth.AccessToken{
			Token:  config.Thirdpart.Twitter.AccessToken,
			Secret: config.Thirdpart.Twitter.AccessSecret,
		},
		helper: helper,
		config: config,
	}
}

func (t *Twitter) Provider() string {
	return "twitter"
}

type twitterReply struct {
	Id        string      `json:"id_str"`
	Recipient twitterInfo `json:"recipient"`
}

func (t *Twitter) Post(id, text string) (string, error) {
	text = strings.Trim(text, " \n\r")
	privateMessage := text
	publicMessage := text
	if lastEnter := strings.LastIndex(text, "\n"); lastEnter >= 0 {
		privateMessage = text[:lastEnter]
		publicMessage = text[lastEnter+1:]
	}
	ret, err := t.sendPrivate(id, privateMessage)
	if err != nil {
		ret, err = t.sendPublic(id, publicMessage)
	}

	if err != nil {
		return "", err
	}
	return ret, nil
}

func (t *Twitter) sendPrivate(id, text string) (string, error) {
	text = strings.Trim(text, " \n\r")
	params := map[string]string{
		"screen_name": id,
		"text":        text,
	}
	resp, err := broker.HttpResponse(t.oauth.Post(twitterApiBase+"direct_messages/new.json", params, t.token))
	if err != nil {
		return "", fmt.Errorf("send to %s@twitter fail: %s", id, err)
	}
	defer resp.Close()
	decoder := json.NewDecoder(resp)
	var reply twitterReply
	err = decoder.Decode(&reply)
	if err != nil {
		return "", fmt.Errorf("parse %s@twitter reply error: %s", id, err)
	}
	go func() {
		to := &model.Recipient{
			ExternalUsername: id,
			Provider:         "twitter",
		}
		err := t.helper.UpdateIdentity(to, reply.Recipient)
		if err != nil {
			logger.ERROR("can't update %s@twitter identity: %s", id, err)
		}
	}()
	return reply.Id, nil
}

func (t *Twitter) sendPublic(id, text string) (string, error) {
	text = strings.Trim(text, " \n\r")
	params := map[string]string{
		"status": fmt.Sprintf("@%s %s", id, text),
	}
	resp, err := broker.HttpResponse(t.oauth.Post(twitterApiBase+"statuses/update.json", params, t.token))
	if err != nil {
		return "", fmt.Errorf("send to %s@twitter fail: %s", id, err)
	}
	defer resp.Close()
	decoder := json.NewDecoder(resp)
	var reply twitterReply
	err = decoder.Decode(&reply)
	if err != nil {
		return "", fmt.Errorf("parse %s@twitter reply error: %s", id, err)
	}
	return reply.Id, nil
}

func (t *Twitter) UpdateIdentity(to *model.Recipient) error {
	k, v := t.identity(to)
	params := map[string]string{k: v}
	resp, err := broker.HttpResponse(t.oauth.Get(twitterApiBase+"users/show.json", params, t.token))
	if err != nil {
		return fmt.Errorf("get %s users/show(%v) failed: %s", to, params, err)
	}
	defer resp.Close()
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
	var access model.OAuthToken
	err := json.Unmarshal([]byte(to.AuthData), &access)
	if err != nil {
		return fmt.Errorf("can't convert %s's AuthData: %s", to, err)
	}

	token := &oauth.AccessToken{
		Token:  access.Token,
		Secret: access.Secret,
	}
	k, v := t.identity(to)
	params := map[string]string{k: v}
	resp, err := broker.HttpResponse(t.oauth.Get(twitterApiBase+"friends/ids.json", params, token))
	if err != nil {
		return fmt.Errorf("get %s friends/ids(%v) failed: %s", to, params, err)
	}
	defer resp.Close()
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
		} else {
			friendIDs = nil
		}

		params := map[string]string{"user_id": join(ids, ",")}
		logger.DEBUG("twitter lookup: %v", params)
		resp, err := broker.HttpResponse(t.oauth.Get(twitterApiBase+"users/lookup.json", params, token))
		if err != nil {
			return fmt.Errorf("get %s users/lookup.json(%v) fail: %s", to, params, err)
		}
		defer resp.Close()
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

type twitterIDs struct {
	Previous string   `json:"previous_cursor_str"`
	Next     string   `json:"next_cursor_str"`
	IDs      []uint64 `json:"ids"`
}

type twitterInfo struct {
	ID              uint64  `json:"id"`
	ScreenName      string  `json:"screen_name"`
	ProfileImageUrl string  `json:"profile_image_url"`
	Description     *string `json:"description"`
	Name_           string  `json:"name"`
}

func (i twitterInfo) ExternalID() string {
	return fmt.Sprintf("%d", i.ID)
}

func (i twitterInfo) Provider() string {
	return "twitter"
}

func (i twitterInfo) ExternalUsername() string {
	return i.ScreenName
}

func (i twitterInfo) Name() string {
	return i.Name_
}

func (i twitterInfo) Bio() string {
	if i.Description == nil {
		return ""
	}
	return *i.Description
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
