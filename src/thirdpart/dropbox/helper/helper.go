package main

import (
	"encoding/json"
	"fmt"
	"github.com/mrjones/oauth"
	"model"
	"net/url"
	"thirdpart/dropbox"
)

func main() {
	clientToken := "5exqb4gosclu5x0"
	clientSecret := "kr3oga1il0eu4mn"

	provider := oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.dropbox.com/1/oauth/request_token",
		AuthorizeTokenUrl: "https://www.dropbox.com/1/oauth/authorize",
		AccessTokenUrl:    "https://api.dropbox.com/1/oauth/access_token",
	}
	consumer := oauth.NewConsumer(clientToken, clientSecret, provider)
	reqToken, u, err := consumer.GetRequestTokenAndUrl("http://abc")
	if err != nil {
		panic(err)
	}
	fmt.Println(u + "&oauth_callback=" + url.QueryEscape("http://abc"))
	code := ""
	fmt.Scanln(&code)
	token, err := consumer.AuthorizeToken(reqToken, code)
	if err != nil {
		panic(err)
	}
	fmt.Println("client token", clientToken)
	fmt.Println("client secret", clientSecret)
	fmt.Println("access token:", token.Token)
	fmt.Println("access secret:", token.Secret)

	t := model.OAuthToken{
		Token:  token.Token,
		Secret: token.Secret,
	}
	str, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	to := model.Recipient{
		IdentityID: 123,
		Provider:   "dropbox",
		AuthData:   string(str),
	}
	var config model.Config
	config.Thirdpart.Dropbox.Key = clientToken
	config.Thirdpart.Dropbox.Secret = clientSecret
	config.AWS.S3.Domain = "127.0.0.1:1234"
	config.AWS.S3.BucketPrefix = "test"
	config.Log, err = logger.New(logger.Stderr, "test")
	if err != nil {
		panic(err)
	}
	dropbox, err := dropbox.New(&config)
	if err != nil {
		panic(err)
	}
	fmt.Println("grabbing...")
	photos, err := dropbox.Grab(to, "/Photos/under water")
	if err != nil {
		panic(err)
	}
	fmt.Println(photos)
}
