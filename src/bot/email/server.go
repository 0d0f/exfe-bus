package email

import (
	"fmt"
	"github.com/googollee/goimap"
	"gobot"
	"launchpad.net/tomb"
	"model"
	"time"
)

type EmailBotServer struct {
	bot    *bot.Bot
	conn   *imap.IMAPClient
	config *model.Config
}

func NewEmailBotServer(c *model.Config) *EmailBotServer {
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
	conn, err := imap.NewClient(s.config.Bot.Email.IMAPHost)
	if err != nil {
		return err
	}
	err = conn.Login(s.config.Bot.Email.IMAPUser, s.config.Bot.Email.IMAPPassword)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *EmailBotServer) Close() error {
	return s.conn.Close()
}

func (s *EmailBotServer) Serve() error {
	s.conn.Select(imap.Inbox)
	ids, err := s.conn.Search("unseen")
	if err != nil {
		return err
	}
	for _, id := range ids {
		s.config.Log.Debug("get mail id: %s\n", id)
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

func Daemon(config *model.Config) *tomb.Tomb {
	t := tomb.Tomb{}

	go func() {
		defer t.Done()

		for {
			s := NewEmailBotServer(config)

			for i := 0; ; i++ {
				err := s.Conn()
				if err == nil {
					break
				}
				if i > 10 {
					i = 0
					config.Log.Crit("email connect error: %s", err)
				}
				config.Log.Err("email connect error: %s", err)
			}
			for {
				err := s.Serve()
				if err != nil {
					config.Log.Crit("email error: %s", err)
					s.Close()
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
