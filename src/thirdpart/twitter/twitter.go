package twitter

import (
	"encoding/json"
	"fmt"
	"formatter"
	"io"
	"model"
	"net/url"
	"regexp"
	"strings"
	"thirdpart"
	"unicode/utf8"
)

type Broker interface {
	Do(accessToken *thirdpart.Token, cmd, path string, params url.Values) (io.ReadCloser, error)
}

type Twitter struct {
	broker      Broker
	accessToken *thirdpart.Token
	helper      thirdpart.Helper
}

const twitterApiBase = "https://api.twitter.com/1.1/"
const provider = "twitter"

func New(accessToken, accessSecret string, broker Broker, helper thirdpart.Helper) *Twitter {
	return &Twitter{
		broker: broker,
		accessToken: &thirdpart.Token{
			Token:  accessToken,
			Secret: accessSecret,
		},
		helper: helper,
	}
}

func (t *Twitter) Provider() string {
	return provider
}

type twitterReply struct {
	IDstr string `json:"id_str"`
}

func (t *Twitter) Send(to *model.Recipient, privateMessage string, publicMessage string, data *model.InfoData) (string, error) {
	ids, err := t.sendPrivate(to, privateMessage)
	if err != nil && strings.Index(err.Error(), `"code":150`) > 0 {
		ids, err = t.sendPublic(to, publicMessage)
	}

	if ids != "" {
		ids = ids[1:]
	}
	return ids, nil
}

func (t *Twitter) sendPrivate(to *model.Recipient, message string) (string, error) {
	params := make(url.Values)
	params.Set(t.identity(to))
	var resp io.ReadCloser
	ids := ""
	for _, line := range strings.Split(message, "\n") {
		line = strings.Trim(line, " \r\n\t")
		if line == "" {
			continue
		}
		cutter, err := formatter.CutterParse(line, twitterLen)
		if err != nil {
			return "", fmt.Errorf("parse %s private message failed: %s", to, err)
		}
		for _, content := range cutter.Limit(140) {
			params.Set("text", content)
			resp, err = t.broker.Do(t.accessToken, "POST", "direct_messages/new.json", params)
			if err != nil {
				return "", err
			}
			decoder := json.NewDecoder(resp)
			var reply twitterReply
			err = decoder.Decode(&reply)
			if err != nil {
				return "", fmt.Errorf("parse %s reply error: %s", to, err)
			}
			ids = fmt.Sprintf("%s,%s", ids, reply.IDstr)
		}
	}
	return ids, nil
}

func (t *Twitter) sendPublic(to *model.Recipient, message string) (string, error) {
	ids := ""
	params := make(url.Values)
	var resp io.ReadCloser
	for _, line := range strings.Split(message, "\n") {
		line = strings.Trim(line, " \r\n\t")
		if line == "" {
			continue
		}
		cutter, err := formatter.CutterParse(line, twitterLen)
		if err != nil {
			return "", fmt.Errorf("parse %s public message failed: %s", to, err)
		}
		for _, content := range cutter.Limit(140 - len(to.ExternalUsername) - 2) {
			params.Set("status", fmt.Sprintf("@%s %s", to.ExternalUsername, content))
			resp, err = t.broker.Do(t.accessToken, "POST", "statuses/update.json", params)
			if err != nil {
				return "", fmt.Errorf("send to %s fail: %s", to, err)
			}
			decoder := json.NewDecoder(resp)
			var reply twitterReply
			err = decoder.Decode(&reply)
			if err != nil {
				return "", fmt.Errorf("parse %s reply error: %s", to, err)
			}
			ids = fmt.Sprintf("%s,%s", ids, reply.IDstr)
		}
	}
	return ids, nil
}

func (t *Twitter) UpdateIdentity(to *model.Recipient) error {
	params := make(url.Values)
	params.Set(t.identity(to))
	resp, err := t.broker.Do(t.accessToken, "GET", "users/show.json", params)
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

	params := make(url.Values)
	params.Set(t.identity(to))
	resp, err := t.broker.Do(access, "GET", "friends/ids.json", params)
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
		resp, err := t.broker.Do(access, "GET", "users/lookup.json", params)
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

var urlRegex = regexp.MustCompile(`(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`)

const twitterUrl = "http://t.co/12345678"

func twitterLen(content string) int {
	content = urlRegex.ReplaceAllString(content, twitterUrl)
	return utf8.RuneCountInString(content)
}
