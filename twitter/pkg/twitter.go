package twitter

import (
	"fmt"
)

type Response struct {
	Error string
	Result string
}

type Tweet struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string
	Tweet        string
}

func (t *Tweet) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) Tweet:%s}", t.ClientToken, t.ClientSecret, t.AccessToken, t.AccessSecret, t.Tweet)
}

type DirectMessage struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string
	Message      string
	ToUserName   string
	ToUserId     string
}

func (m *DirectMessage) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) ToUser:%s(%s) Message:%s}",
		m.ClientToken, m.ClientSecret, m.AccessToken, m.AccessSecret, m.ToUserName, m.ToUserId, m.Message)
}

type Friendship struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string

	UserA string
	UserB string
}

func (f *Friendship) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) UserA %s UserB %s}",
		f.ClientToken, f.ClientSecret, f.AccessToken, f.AccessSecret,
		f.UserA, f.UserB)
}

type UserInfo struct {
	ClientToken  string
	ClientSecret string
	AccessToken  string
	AccessSecret string

	UserId     string
	ScreenName string
}

func (i *UserInfo) GoString() string {
	return fmt.Sprintf("{Client:(%s %s) Access:(%s %s) User %s(%s)}",
		i.ClientToken, i.ClientSecret, i.AccessToken, i.AccessSecret, i.ScreenName, i.UserId)
}
