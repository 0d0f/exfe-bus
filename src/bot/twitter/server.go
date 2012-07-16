package twitter

import (
	"exfe/service"
	"fmt"
	"gobot"
	"log"
	"oauth"
	"os"
)

type Server struct {
	config *exfe_service.Config
	bot    *bot.Bot
	reader *StreamingReader
}

func NewServer(c *exfe_service.Config) *Server {
	b := bot.NewBot(NewBot(c))
	b.RegisterPretreat("Idle")
	b.Register("TweetWithIOM")
	b.Register("Default")
	return &Server{
		config: c,
		bot:    b,
	}
}

func (s *Server) Conn() error {
	client := oauth.CreateClient(s.config.Twitter.Client_token,
		s.config.Twitter.Client_secret,
		s.config.Twitter.Access_token,
		s.config.Twitter.Access_secret,
		"https://userstream.twitter.com")
	reader, err := client.Do("GET", "/2/user.json", nil)
	if err != nil {
		return err
	}
	s.reader = NewStreamingReader(reader)
	return nil
}

func Daemon(config *exfe_service.Config, quit chan int) {
	l := log.New(os.Stderr, "exfe.bot.twitter", log.LstdFlags)
	s := NewServer(config)
	err := s.Conn()
	if err != nil {
		l.Printf("can't connect to twitter: %s", err)
		quit <- 1
	}
	input := make(chan *Tweet)
	e := make(chan error)
	reader := s.reader
	go func() {
		for {
			fmt.Printf(".")
			tweet, err := reader.ReadTweet()
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
			err := s.bot.Feed(i)
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
