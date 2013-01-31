package main

import (
	"fmt"
	"github.com/mrjones/oauth"
	"net/url"
)

func main() {
	provider := oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.dropbox.com/1/oauth/request_token",
		AuthorizeTokenUrl: "https://www.dropbox.com/1/oauth/authorize",
		AccessTokenUrl:    "https://api.dropbox.com/1/oauth/access_token",
	}
	clientToken := "5exqb4gosclu5x0"
	clientSecret := "kr3oga1il0eu4mn"
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
}
