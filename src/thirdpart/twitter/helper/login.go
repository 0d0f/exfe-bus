package main

import (
	"fmt"
	"model"
	"oauth"
	"os"
	"thirdpart"
	"thirdpart/twitter"
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

	client, err := request.AccessClient(temp_token, verifier, "http://api.twitter.com/1.1/")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("\nClient token: ", client.ClientToken.Token)
	fmt.Println("Client secret:", client.ClientToken.Secret)
	fmt.Println("Access token: ", client.AccessToken.Token)
	fmt.Println("Access secret:", client.AccessToken.Secret)
	fmt.Println("Api base uri: ", client.ApiBaseUri)

	helper := new(thirdpart.HelperFake)
	clientToken := &thirdpart.Token{client.ClientToken.Token, client.ClientToken.Secret}
	accessToken := &thirdpart.Token{client.AccessToken.Token, client.AccessToken.Secret}
	twitter := twitter.New(clientToken, accessToken, helper)

	to := &model.Identity{
		ID:       123,
		Name:     "tester",
		Nickname: "tester nick",
		Bio:      "bio",
		Timezone: "+0800",
		UserID:   789,
		Avatar:   "avatar",

		Provider:         "twitter",
		ExternalID:       "56591660",
		ExternalUsername: "googollee",
		OAuthToken:       `{"oauth_token":"56591660-EXgdQxYxUiWocQ5krbreCYIVixDt21NDBgjPR7kKr","oauth_token_secret":"V0YRtAozGBJdPnghMRVRzzV65dtpWeLwmK00DVQC5X0"}`,
	}

	// fmt.Println()
	// err = twitter.UpdateIdentity(to)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println()
	// err = twitter.UpdateFriends(to)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println()
	// to.ExternalID = "247228987"
	to.ExternalUsername = "exfe"
	err = twitter.Send(to, "private", "@googollee public")
	if err != nil {
		fmt.Println(err)
	}
}
