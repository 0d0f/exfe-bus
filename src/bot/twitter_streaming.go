package main

import (
	"oauth"
	"io"
	"fmt"
	"encoding/json"
	"bytes"
	"strings"
	"twitter/service"
	"log"
	"io/ioutil"
	"exfe/service"
	"net/url"
	"net/http"
)

func find(data []byte, c rune) int {
	for i, d := range data {
		if rune(d) == c {
			return i
		}
	}
	return -1
}

func connTwitter(clientToken, clientSecret, accessToken, accessSecret string) io.ReadCloser {
	client := oauth.CreateClient(clientToken, clientSecret, accessToken, accessSecret, "https://userstream.twitter.com")
	reader, err := client.Do("GET", "/2/user.json", nil)
	if err != nil {
		panic(err)
	}
	reader, err = client.Do("GET", "/2/user.json", nil)
	if err != nil {
		panic(err)
	}
	return reader
}

func read(clientToken, clientSecret, accessToken, accessSecret string, reader io.ReadCloser, ret chan Tweet) {
	var cache []byte
	var buf [20]byte
	for {
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println(err)
			reader = connTwitter(clientToken, clientSecret, accessToken, accessSecret)
			continue
		}

		cache = parseBuf(buf[0:n], cache, ret)
	}
}

func parseBuf(buf []byte, cache []byte, ret chan Tweet) []byte {
	for {
		i := find(buf, '\r')
		if i < 0 {
			return append(cache, buf...)
		} else {
			cache = append(cache, buf[0:i]...)
			item := strings.Trim(string(cache), "\r\n")
			cache = nil
			buf = buf[(i + 1):]

			var t Tweet
			buf := bytes.NewBufferString(item)
			decoder := json.NewDecoder(buf)
			err := decoder.Decode(&t)
			if err == nil && (t.User != nil || t.Direct_message != nil) {
				ret <- t
			}
		}
	}
	return nil
}

func connStreaming(clientToken, clientSecret, accessToken, accessSecret string) (chan Tweet, error) {
	reader := connTwitter(clientToken, clientSecret, accessToken, accessSecret)
	ret := make(chan Tweet)
	go read(clientToken, clientSecret, accessToken, accessSecret, reader, ret)

	return ret, nil
}

func sendHelp(screen_name string) {
	f := &twitter_service.FriendshipsExistsArg{
		ClientToken:  config.Twitter.Client_token,
		ClientSecret: config.Twitter.Client_secret,
		AccessToken:  config.Twitter.Access_token,
		AccessSecret: config.Twitter.Access_secret,
		UserA:        screen_name,
		UserB:        config.Twitter.Screen_name,
	}
	var isFriend bool
	err := client.Do("GetFriendship", f, &isFriend, 10)
	if err != nil {
		log.Printf("Can't require user %s friendship: %s", screen_name, err)
		isFriend = false
	}

	if isFriend {
		dm := &twitter_service.DirectMessagesNewArg{
			ClientToken:  config.Twitter.Client_token,
			ClientSecret: config.Twitter.Client_secret,
			AccessToken:  config.Twitter.Access_token,
			AccessSecret: config.Twitter.Access_secret,
			Message:      helper,
			ToUserName:   &screen_name,
		}
		client.Send("SendDM", dm, 5)
	} else {
		tweet := &twitter_service.StatusesUpdateArg{
			ClientToken:  config.Twitter.Client_token,
			ClientSecret: config.Twitter.Client_secret,
			AccessToken:  config.Twitter.Access_token,
			AccessSecret: config.Twitter.Access_secret,
			Tweet:        fmt.Sprintf("@%s %s", screen_name, helper),
		}
		client.Send("SendTweet", tweet, 5)
	}
}

func processTwitter(config *exfe_service.Config) {
	c, _ := connStreaming(config.Twitter.Client_token, config.Twitter.Client_secret, config.Twitter.Access_token, config.Twitter.Access_secret)

	for t := range c {
		hash, post := t.parse()
		time := t.created_at()
		external_id := t.external_id()
		screen_name := t.screen_name()

		fmt.Println(hash, time, external_id, screen_name, post)

		if screen_name == "" {
			continue
		}

		if hash == "" && post != "" {
			sendHelp(screen_name)
			continue
		}

		params := make(url.Values)
		params.Add("per_user_hash", hash)
		params.Add("content", post)
		params.Add("external_id", external_id)
		params.Add("provider", "twitter")
		params.Add("time", time)
		resp, err := http.PostForm(fmt.Sprintln("%s/v2/gobus/PostConversation", config.Site_api), params)
		if err != nil {
			log.Printf("Send post to server error: %s", err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Get response body error: %s", err)
			continue
		}
		if resp.StatusCode == 500 {
			log.Printf("Server inner error: %s", string(body))
			continue
		}
		if resp.StatusCode == 400 {
			log.Printf("User status error: %s", string(body))
			continue
		}
	}
}
