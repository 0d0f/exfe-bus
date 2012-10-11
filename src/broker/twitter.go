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

func NewTwitter(clientToken, clientSecret, accessToken, accessSecret string) *Twitter {
	return &Twitter{
		client: oauth.CreateClient(clientToken, clientSecret, accessToken, accessSecret, twitterApiBase),
		clientToken: &thirdpart.Token{
			Token:  clientToken,
			Secret: clientSecret,
		},
		accessToken: &thirdpart.Token{
			Token:  accessToken,
			Secret: accessSecret,
		},
	}
}

func (t *Twitter) Do(cmd, url string, params url.Values) (io.ReadCloser, error) {
	return t.client.Do(cmd, url, params)
}
