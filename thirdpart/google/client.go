package google

import (
	"code.google.com/p/goauth2/oauth"
	"exfe/model"
	"exfe/service"
	"github.com/googollee/go-logger"
	"net/http"
)

type Client struct {
	http *http.Client
	log  *logger.Logger
}

func NewClient(config *exfe_service.Config, id *exfe_model.Receiver) *Client {
	c := &oauth.Config{
		ClientId:     config.Google.ID,
		ClientSecret: config.Google.Secret,
		Scope:        "https://www.googleapis.com/oauth2/v1/userinfo https://www.google.com/m8/feeds",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
	}
	t := oauth.Transport{Config: c}

	token := oauth.Token{
		AccessToken: id.Data["token"],
	}

	t.Token = token
	return &Client{
		http: t.Client(),
	}
}
