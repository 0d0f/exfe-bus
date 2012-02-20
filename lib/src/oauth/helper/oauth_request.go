package main

import (
	"oauth"
	"fmt"
	"os"
)

func main() {
	request := oauth.CreateOAuthRequest("VC3OxLBNSGPLOZ2zkgisA", "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A",
		"http://api.twitter.com/oauth/request_token", "http://api.twitter.com/oauth/authenticate", "http://api.twitter.com/oauth/access_token")

	temp_token, auth_url, err := request.AuthorizeUrl("http://callback", nil, nil)
	if (err != nil) {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(auth_url)
	var verifier string
	_, err = fmt.Scan(&verifier)
	if (err != nil) {
		fmt.Println(err)
		os.Exit(-1)
	}

	client, err := request.AccessClient(temp_token, verifier, "http://api.twitter.com/1/")
	if (err != nil) {
		fmt.Println(err)
		os.Exit(-1)
	}

	f, err := os.Create("./oauth_client.json")
	defer func () { f.Close() }()
	if (err != nil) {
		fmt.Println(err)
		os.Exit(-1)
	}

	client.Dump(f)
}
