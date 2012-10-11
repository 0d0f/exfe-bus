package broker

import (
	"io"
	"net/url"
	"oauth"
	"thirdpart"
)

type Twitter struct {
	token *thirdpart.Token
}

const twitterApiBase = "https://api.twitter.com/1.1/"

func NewTwitter(clientToken, clientSecret string) *Twitter {
	return &Twitter{
		token: &thirdpart.Token{
			Token:  clientToken,
			Secret: clientSecret,
		},
	}
}

func (t *Twitter) Do(accessToken *thirdpart.Token, cmd, url string, params url.Values) (io.ReadCloser, error) {
	client := oauth.CreateClient(t.token.Token, t.token.Secret, accessToken.Token, accessToken.Secret, twitterApiBase)
	return client.Do(cmd, url, params)
}
