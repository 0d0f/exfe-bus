package broker

import (
	"fmt"
	"github.com/mrjones/oauth"
	"io"
	"io/ioutil"
	"model"
	"net/http"
)

type Twitter interface {
	Do(accessToken model.OAuthToken, cmd, url string, params map[string]string) (io.ReadCloser, error)
}

type TwitterImpl struct {
	consumer *oauth.Consumer
}

const twitterApiBase = "https://api.twitter.com/1.1/"

func NewTwitter(clientToken, clientSecret string) *TwitterImpl {
	provider := oauth.ServiceProvider{
		RequestTokenUrl:   "http://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	}
	consumer := oauth.NewConsumer(clientToken, clientSecret, provider)
	return &TwitterImpl{
		consumer: consumer,
	}
}

func (t *TwitterImpl) Do(accessToken model.OAuthToken, cmd, url string, params map[string]string) (io.ReadCloser, error) {
	token := &oauth.AccessToken{
		Token:  accessToken.Token,
		Secret: accessToken.Secret,
	}
	var err error
	var resp *http.Response
	switch cmd {
	case "GET":
		resp, err = t.consumer.Get(url, params, token)
	case "POST":
		resp, err = t.consumer.Post(url, params, token)
	case "DELETE":
		resp, err = t.consumer.Delete(url, params, token)
	default:
		return nil, fmt.Errorf("can't handle %s", cmd)
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %s", resp.Status, content)
	}
	return resp.Body, nil
}
