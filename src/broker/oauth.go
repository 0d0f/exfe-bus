package broker

import (
	"github.com/mrjones/oauth"
)

type OAuth struct {
	*oauth.Consumer
}

func NewOAuth(client, secret string, provider oauth.ServiceProvider) OAuth {
	consumer := oauth.NewConsumer(client, secret, provider)
	consumer.HttpClient = HttpClient
	return OAuth{consumer}
}
