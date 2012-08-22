package email

import (
	"exfe/service"
	"fmt"
	"github.com/googollee/goimap"
	"gobot"
	"log"
	"os"
	"time"
)

type EmailBotServer struct {
	bot    *bot.Bot
	conn   *imap.IMAPClient
	config *exfe_service.Config
}

func NewEmailBotServer(c *exfe_service.Config) *EmailBotServer {
	b := bot.NewBot(NewEmailBot(c))
	b.Register("EmailWithCrossID")
	b.Register("Default")
	return &EmailBotServer{
		bot:    b,
		conn:   nil,
		config: c,
	}
}

func (s *EmailBotServer) Conn() error {
	conn, err := imap.NewClient(s.config.Bot.Imap_host)
	if err != nil {
		return err
	}
	err = conn.Login(s.config.Bot.Imap_user, s.config.Bot.Imap_password)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *EmailBotServer) Serve() error {
	s.conn.Select(imap.Inbox)
	ids, err := s.conn.Search("unseen")
	if err != nil {
		return err
	}
	for _, id := range ids {
		fmt.Printf("get mail id: %s\n", id)
		msg, err := s.conn.GetMessage(id)
		if err != nil {
			return fmt.Errorf("Get message(%v) error: %s", id, err)
		}
		err = s.bot.Feed(msg)
		switch err {
		case nil:
			s.conn.Do(fmt.Sprintf("copy %s posted", id))
		default:
			s.conn.Do(fmt.Sprintf("copy %s error", id))
			if err != badrequest {
				return fmt.Errorf("Process message(%v) error: %s", id, err)
			}
		}
		s.conn.StoreFlag(id, imap.Deleted)
	}
	return nil
}

func Daemon(config *exfe_service.Config, quit chan int) {
	l := log.New(os.Stderr, "exfe.bot.email", log.LstdFlags)
	s := NewEmailBotServer(config)
	for {
		err := s.Conn()
		if err == nil {
			break
		}
		l.Printf("email connect error: %s", err)
	}
	for {
		err := s.Serve()
		if err != nil {
			l.Printf("email error: %s", err)
			quit <- 1
			return
		}
		select {
		case <-quit:
			quit <- 1
			return
		case <-time.After(s.config.Bot.Imap_time_out * time.Second):
		}
	}
}
