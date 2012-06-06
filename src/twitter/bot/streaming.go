package main

import (
	"exfe/service"
	"oauth"
	"io"
	"fmt"
	"encoding/json"
	"bytes"
	"strings"
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
			buf = buf[(i+1):]

			var t Tweet
			buf := bytes.NewBufferString(item)
			decoder := json.NewDecoder(buf)
			err := decoder.Decode(&t)
			if err == nil && (t.User != nil || t.Direct_message != nil) {
				ret <- t
			} else {
				// send helper
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

func main() {
	config := exfe_service.InitConfig()
	Init(config.Twitter.Screen_name)

	c, _ := connStreaming(config.Twitter.Client_token, config.Twitter.Client_secret, config.Twitter.Access_token, config.Twitter.Access_secret)

	for t := range c {
		hash, post := t.parse()
		time := t.created_at()
		external_id := t.external_id()

		fmt.Println(hash, time, external_id, post)
	}
}
