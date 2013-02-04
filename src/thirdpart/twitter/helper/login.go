package main

import (
	"fmt"
	"github.com/mrjones/oauth"
)

func main() {
	clientToken, clientSecret := "VC3OxLBNSGPLOZ2zkgisA", "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A"
	provider := oauth.ServiceProvider{
		RequestTokenUrl:   "http://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	}
	consumer := oauth.NewConsumer(clientToken, clientSecret, provider)
	reqToken, u, err := consumer.GetRequestTokenAndUrl("http://abc")
	if err != nil {
		panic(err)
	}
	fmt.Println(u)
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
