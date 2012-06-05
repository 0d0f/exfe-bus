package main

import (
	"fmt"
	"oauth"
	"os"
)

func main() {
	request := oauth.CreateOAuthRequest("VC3OxLBNSGPLOZ2zkgisA", "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A",
		"http://api.twitter.com/oauth/request_token", "http://api.twitter.com/oauth/authorize", "http://api.twitter.com/oauth/access_token")

	temp_token, auth_url, err := request.AuthorizeUrl("http://callback", nil, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(auth_url)
	var verifier string
	_, err = fmt.Scan(&verifier)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	client, err := request.AccessClient(temp_token, verifier, "http://api.twitter.com/1/")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("\nClient token: ", client.ClientToken.Token)
	fmt.Println("Client secret:", client.ClientToken.Secret)
	fmt.Println("Access token: ", client.AccessToken.Token)
	fmt.Println("Access secret:", client.AccessToken.Secret)
	fmt.Println("Api base uri: ", client.ApiBaseUri)
}
