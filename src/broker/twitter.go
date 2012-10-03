package broker

import (
	"io"
	"net/url"
	"oauth"
	"thirdpart"
)

type Twitter struct {
	client      *oauth.OAuthClient
	clientToken *thirdpart.Token
	accessToken *thirdpart.Token
}

const twitterApiBase = "https://api.twitter.com/1.1/"

func NewTwitter(client, access *thirdpart.Token) *Twitter {
	return &Twitter{
		client:      oauth.CreateClient(client.Token, client.Secret, access.Token, access.Secret, twitterApiBase),
		clientToken: client,
		accessToken: access,
	}
}

func (t *Twitter) Do(cmd, url string, params url.Values) (io.ReadCloser, error) {
	return t.client.Do(cmd, url, params)
}
