package email

import (
	"broker"
	"fmt"
	"formatter"
	"github.com/googollee/goimap"
	"gobot"
	"launchpad.net/tomb"
	"model"
	"net"
	"strings"
	"time"
)

type EmailBotServer struct {
	bot    *bot.Bot
	conn   net.Conn
	client *imap.IMAPClient
	config *model.Config
}

func NewEmailBotServer(c *model.Config, localTemplate *formatter.LocalTemplate, sender *broker.Sender) *EmailBotServer {
	b := bot.NewBot(NewEmailBot(c, localTemplate, sender))
	b.Register("EmailWithCrossID")
	b.Register("Default")
	return &EmailBotServer{
		bot:    b,
		conn:   nil,
		client: nil,
		config: c,
	}
}

func (s *EmailBotServer) Conn() error {
	conn, err := net.DialTimeout("tcp", s.config.Bot.Email.IMAPHost, time.Second)
	if err != nil {
		return err
	}
	host := s.config.Bot.Email.IMAPHost[:strings.Index(s.config.Bot.Email.IMAPHost, ":")]
	conn.SetDeadline(time.Now().Add(time.Second))
	client, err := imap.NewClient(conn, host)
	if err != nil {
		return err
	}
	err = client.Login(s.config.Bot.Email.IMAPUser, s.config.Bot.Email.IMAPPassword)
	if err != nil {
		return err
	}
	s.conn = conn
	s.client = client
	return nil
}

func (s *EmailBotServer) Close() error {
	s.conn.SetDeadline(time.Now().Add(time.Second))
	return s.client.Close()
}

func (s *EmailBotServer) Serve() error {
	s.conn.SetDeadline(time.Now().Add(time.Second))
	s.client.Select(imap.Inbox)
	ids, err := s.client.Search("unseen")
	if err != nil {
		return err
	}
	for _, id := range ids {
		s.conn.SetDeadline(time.Now().Add(time.Second))
		s.config.Log.Debug("get mail id: %s\n", id)
		msg, err := s.client.GetMessage(id)
		if err != nil {
			return fmt.Errorf("Get message(%v) error: %s", id, err)
		}
		err = s.bot.Feed(msg)
		switch err {
		case nil:
			s.client.Do(fmt.Sprintf("copy %s posted", id))
		default:
			s.client.Do(fmt.Sprintf("copy %s error", id))
			if err != badrequest {
				return fmt.Errorf("Process message(%v) error: %s", id, err)
			}
		}
		s.client.StoreFlag(id, imap.Deleted)
	}
	return nil
}

func Daemon(config *model.Config, localTemplate *formatter.LocalTemplate, sender *broker.Sender) *tomb.Tomb {
	t := tomb.Tomb{}

	go func() {
		defer t.Done()

		for {
			s := NewEmailBotServer(config, localTemplate, sender)

			for i := 0; ; i++ {
				err := s.Conn()
				if err == nil {
					break
				}
				time.After(time.Duration(s.config.Bot.Email.TimeoutInSecond) * time.Second)
				if i > 10 {
					i = 0
					config.Log.Crit("email connect error: %s", err)
				}
				config.Log.Err("email connect error: %s", err)
			}
			for {
				err := s.Serve()
				if err != nil {
					config.Log.Err("email error: %s", err)
					break
				}
				select {
				case <-t.Dying():
					return
				case <-time.After(time.Duration(s.config.Bot.Email.TimeoutInSecond) * time.Second):
				}
			}
		}
	}()

	return &t
}
