package twitter

import (
	"exfe/service"
	"fmt"
	"gobot"
	"log"
	"oauth"
	"os"
)

func Daemon(config *exfe_service.Config, quit chan int) {
	l := log.New(os.Stderr, "exfe.bot.twitter", log.LstdFlags)

	b := bot.NewBot(NewBot(config))
	b.RegisterPretreat("Idle")
	b.Register("TweetWithIOM")
	b.Register("Default")

	client := oauth.CreateClient(config.Twitter.Client_token,
		config.Twitter.Client_secret,
		config.Twitter.Access_token,
		config.Twitter.Access_secret,
		"https://userstream.twitter.com")
	reader, err := client.Do("GET", "/2/user.json", nil)
	if err != nil {
		l.Printf("can't connect to twitter: %s", err)
		quit <- 1
		return
	}
	r := NewStreamingReader(reader)

	input := make(chan *Tweet)
	e := make(chan error)
	go func() {
		for {
			fmt.Printf(".")
			tweet, err := r.ReadTweet()
			if err != nil {
				e <- err
			}
			input <- tweet
		}
	}()
	for {
		select {
		case i := <-input:
			l.Printf("receive: %s", i)
			err := b.Feed(i)
			if err != nil {
				l.Printf("process input error: %s", err)
			}
		case err := <-e:
			l.Printf("convert input error: %s", err)
		case <-quit:
			l.Printf("quit")
			quit <- 1
			return
		}
	}
}
