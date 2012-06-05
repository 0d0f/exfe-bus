package main

import (
	"oauth"
	"io"
	"fmt"
	"encoding/json"
	"bytes"
	"strings"
)

type DirectMessage struct {
	Sender struct {
		Id_str string
		Screen_name string
	}
	Text string
}

type Tweet struct {
	Entities struct {
		User_mentions []struct {
			Screen_name string
			Id_str string
		}
	}
	Text string
	In_reply_to_status_id_str *string
	User *struct {
		Screen_name string
		Id_str string
	}
	Direct_message *DirectMessage
}

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
	c, _ := connStreaming("VC3OxLBNSGPLOZ2zkgisA", "Lg6b5eHdPLFPsy4pI2aXPn6qEX6oxTwPyS0rr2g4A", "491159882-kPatxldk69pToBLBaiCGabZcCkxxTwjXBc6UQwBU", "uIYJyzC8kavcZkj25wVDjlI7GdG0o1X0dqJJSJum1Y")

	for t := range c {
		fmt.Println(t)
		if t.Direct_message != nil {
			fmt.Println(*t.Direct_message)
		}
	}
}
