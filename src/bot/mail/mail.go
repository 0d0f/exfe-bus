package mail

import (
	"broker"
	"bytes"
	"code.google.com/p/go-imap/go1/imap"
	"crypto/tls"
	"fmt"
	"formatter"
	"github.com/googollee/go-logger"
	"launchpad.net/tomb"
	"model"
	"net"
	"net/mail"
	"strings"
	"time"
)

type MessageIDSaver interface {
	Save(id []string, crossID string) error
	Check(id []string) (string, bool, error)
}

type Worker struct {
	Tomb tomb.Tomb

	config   *model.Config
	log      *logger.SubLogger
	templ    *formatter.LocalTemplate
	platform *broker.Platform
	saver    MessageIDSaver
}

func New(config *model.Config, templ *formatter.LocalTemplate, platform *broker.Platform, saver MessageIDSaver) (*Worker, error) {

	return &Worker{
		config:   config,
		log:      config.Log.SubPrefix("mail"),
		templ:    templ,
		platform: platform,
		saver:    saver,
	}, nil
}

func (w *Worker) Daemon() {
	defer w.Tomb.Done()

	timeout := time.Duration(w.config.Bot.Email.TimeoutInSecond) * time.Second

	for true {
		select {
		case <-w.Tomb.Dying():
			w.log.Notice("quitted")
			return
		case <-time.After(timeout):
			w.process()
		}
	}
}

func (w *Worker) process() {
	w.log.Debug("process...")
	w.log.Notice("process...")
	conn, imapConn, err := w.login()
	if err != nil {
		w.log.Err("can't connect to %s: %s", w.config.Bot.Email.IMAPHost, err)
		return
	}
	w.log.Notice("logined")
	defer imapConn.Logout(broker.NetworkTimeout)
	conn.SetDeadline(time.Now().Add(broker.ProcessTimeout))

	_, err = imapConn.Select("INBOX", false)
	if err != nil {
		w.log.Err("can't select INBOX: %s", err)
		return
	}
	w.log.Notice("selected inbox")

	cmd, err := imap.Wait(imapConn.Search("UNSEEN"))
	if err != nil {
		w.log.Err("can't seach UNSEEN: %s", err)
		return
	}
	w.log.Notice("done search unseen")
	var ids []uint32
	for _, resp := range cmd.Data {
		ids = append(ids, resp.SearchResults()...)
	}

	var errorIds []uint32
	var okIds []uint32
	for _, id := range ids {
		conn.SetDeadline(time.Now().Add(broker.ProcessTimeout))

		msg, err := w.getMail(imapConn, id)
		if err != nil {
			w.log.Err("can't get mail %d: %s", id, err)
			errorIds = append(errorIds, id)
			continue
		}
		w.log.Notice("got %d", id)
		parser, err := NewParser(msg, w.config)
		if err != nil {
			w.log.Err("parse mail %d failed: %s", id, err)
			errorIds = append(errorIds, id)
			continue
		}
		to, toID := parser.GetTypeID()
		if to == "" {
			crossID, exist, err := w.saver.Check(parser.GetIDs())
			if err != nil {
				w.log.Crit("saver check %s failed: %s", id, err)
				errorIds = append(errorIds, id)
				continue
			}
			if exist {
				to, toID = "cross_id", crossID
			}
		}
		cross := parser.GetCross()
		if to == "" {
			crossID, code, err := w.platform.BotCrossGather(cross)
			if err != nil {
				if code < 500 {
					w.sendHelp(code, err, parser)
				} else {
					w.log.Err("can't gather (%s) with: %+v", err, cross)
				}
				errorIds = append(errorIds, id)
				continue
			}
			to, toID = "cross_id", fmt.Sprintf("%d", crossID)
		} else {
			post := parser.GetPost()
			if post != "" {
				_, err := w.platform.BotPostConversation(parser.from.Address, post, parser.Date(), to, toID)
				if err != nil {
					w.log.Err("%s can't post %s with: %s", parser.from.Address, post, err)
					errorIds = append(errorIds, id)
					continue
				}
			}

			if !parser.HasICS() {
				cross.Title = ""
				cross.Description = ""
				cross.Place = nil
				cross.Time = nil
			}
			_, err = w.platform.BotCrossUpdate(to, toID, cross, cross.By)
			if err != nil {
				w.log.Err("%s can't update %s %s: %s", parser.from.Address, to, toID, err)
				errorIds = append(errorIds, id)
				continue
			}
		}
		if to == "cross_id" {
			err = w.saver.Save(parser.GetIDs(), toID)
			if err != nil {
				w.log.Crit("saver save %s failed: %s", id, err)
			}
		}
		w.log.Notice("parsed %d", id)
		okIds = append(okIds, id)
		w.log.Notice("handled %d", id)
	}
	w.log.Notice("id:%v, ok:%v, err:%v", ids, okIds, errorIds)
	if err := w.copy(imapConn, okIds, "posted"); err != nil {
		w.log.Err("can't copy %v to posted: %s", errorIds, err)
	}
	w.log.Notice("copied %v to posted", okIds)
	if err := w.copy(imapConn, errorIds, "error"); err != nil {
		w.log.Err("can't copy %v to error: %s", errorIds, err)
	}
	w.log.Notice("copied %v to error", errorIds)
	if err := w.delete(imapConn, ids); err != nil {
		w.log.Err("can't remove %v from inbox: %s", ids, err)
	}
	w.log.Notice("removed %v", ids)
}

func (w *Worker) copy(conn *imap.Client, ids []uint32, folder string) error {
	if len(ids) == 0 {
		return nil
	}
	set := new(imap.SeqSet)
	set.AddNum(ids...)
	_, err := imap.Wait(conn.Copy(set, folder))
	return err
}

func (w *Worker) delete(conn *imap.Client, ids []uint32) error {
	if len(ids) == 0 {
		return nil
	}
	set := new(imap.SeqSet)
	set.AddNum(ids...)
	_, err := imap.Wait(conn.Store(set, "FLAGS", "\\Deleted"))
	return err
}

func (w *Worker) login() (net.Conn, *imap.Client, error) {
	c, err := net.DialTimeout("tcp", w.config.Bot.Email.IMAPHost, broker.NetworkTimeout)
	if err != nil {
		return nil, nil, err
	}
	tlsConn := tls.Client(c, nil)

	conn, err := imap.NewClient(tlsConn, strings.Split(w.config.Bot.Email.IMAPHost, ":")[0], broker.NetworkTimeout)
	if err != nil {
		return nil, nil, err
	}
	c.SetDeadline(time.Now().Add(broker.NetworkTimeout))

	conn.Data = nil
	if conn.Caps["STARTTLS"] {
		conn.StartTLS(nil)
	}

	if conn.State() == imap.Login {
		_, err = conn.Login(w.config.Bot.Email.IMAPUser, w.config.Bot.Email.IMAPPassword)
		if err != nil {
			return nil, nil, err
		}
	}

	return c, conn, nil
}

func (w *Worker) getMail(conn *imap.Client, id uint32) (*mail.Message, error) {
	buf := bytes.NewBuffer(nil)
	set := new(imap.SeqSet)
	set.AddNum(id)

	cmd, err := imap.Wait(conn.Fetch(set, "RFC822"))
	if err != nil {
		return nil, err
	}
	buf.Write(imap.AsBytes(cmd.Data[0].MessageInfo().Attrs["RFC822"]))

	msg, err := mail.ReadMessage(buf)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (w *Worker) sendHelp(code int, err error, parser *Parser) error {
	w.log.Info("|help|%s|", parser.from.Address)
	buf := bytes.NewBuffer(nil)
	type Email struct {
		From      *mail.Address
		Subject   string
		CrossID   string
		Date      time.Time
		MessageID string
		Text      string

		Config *model.Config
	}
	email := Email{
		From:      parser.from,
		Subject:   parser.subject,
		MessageID: parser.messageID,
		Text:      parser.content,
		Config:    w.config,
	}
	err = w.templ.Execute(buf, "en_US", "email/conversation_reply", email)
	if err != nil {
		w.log.Crit("template(conversation_reply.email) failed: %s", err)
	}

	to := model.Recipient{
		Provider:         "email",
		Name:             parser.from.Name,
		ExternalID:       parser.from.Address,
		ExternalUsername: parser.from.Address,
	}
	_, err = w.platform.Send(to, buf.String())
	return err
}
