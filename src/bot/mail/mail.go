package mail

import (
	"broker"
	"bytes"
	"code.google.com/p/go-imap/go1/imap"
	"crypto/tls"
	"fmt"
	"formatter"
	"github.com/googollee/go-aws/s3"
	"launchpad.net/tomb"
	"logger"
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
	templ    *formatter.LocalTemplate
	platform *broker.Platform
	saver    MessageIDSaver
	bucket   *s3.Bucket
}

func New(config *model.Config, templ *formatter.LocalTemplate, platform *broker.Platform, saver MessageIDSaver) (*Worker, error) {
	aws := s3.New(config.AWS.S3.Domain, config.AWS.S3.Key, config.AWS.S3.Secret)
	aws.SetACL(s3.ACLPublicRead)
	aws.SetLocationConstraint(s3.LC_AP_SINGAPORE)
	bucket, err := aws.GetBucket(fmt.Sprintf("%s-email", config.AWS.S3.BucketPrefix))
	if err != nil {
		return nil, err
	}
	return &Worker{
		config:   config,
		templ:    templ,
		platform: platform,
		saver:    saver,
		bucket:   bucket,
	}, nil
}

func (w *Worker) Daemon() {
	defer w.Tomb.Done()

	timeout := time.Duration(w.config.Bot.Email.TimeoutInSecond) * time.Second

	for true {
		select {
		case <-w.Tomb.Dying():
			return
		case <-time.After(timeout):
			w.process()
		}
	}
}

func (w *Worker) process() {
	conn, imapConn, err := w.login()
	if err != nil {
		logger.ERROR("can't connect to %s: %s", w.config.Bot.Email.IMAPHost, err)
		return
	}
	defer imapConn.Logout(broker.NetworkTimeout)
	conn.SetDeadline(time.Now().Add(broker.ProcessTimeout))

	_, err = imapConn.Select("INBOX", false)
	if err != nil {
		logger.ERROR("can't select INBOX: %s", err)
		return
	}

	cmd, err := imap.Wait(imapConn.Search("UNSEEN"))
	if err != nil {
		logger.ERROR("can't seach UNSEEN: %s", err)
		return
	}
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
			logger.ERROR("can't get mail %d: %s", id, err)
			errorIds = append(errorIds, id)
			continue
		}
		parser, err := NewParser(msg, w.config, w.bucket)
		if err != nil {
			logger.ERROR("parse mail %d failed: %s", id, err)
			errorIds = append(errorIds, id)
			continue
		}
		if strings.HasSuffix(parser.from.Address, "googlemail.com") {
			errorIds = append(errorIds, id)
			continue
		}
		to, toID := parser.GetTypeID()
		fromCalendar := false
		if to == "" {
			crossID, exist, err := w.saver.Check(parser.referenceIDs)
			if err != nil {
				logger.ERROR("saver check %s failed: %s", id, err)
				errorIds = append(errorIds, id)
				continue
			}
			if exist {
				to, toID = "cross_id", crossID
			}
		}
		if to == "" && parser.event != nil {
			crossID, exist, err := w.saver.Check([]string{parser.event.ID})
			if err != nil {
				logger.ERROR("saver check %s failed: %s", id, err)
				errorIds = append(errorIds, id)
				continue
			}
			if exist {
				to, toID = "cross_id", crossID
				fromCalendar = true
			}
		}
		cross := parser.GetCross()
		if to == "" {
			cross, err := w.platform.BotCrossGather(cross)
			if err != nil {
				if warning, ok := err.(broker.Warning); ok {
					w.sendHelp(warning, parser)
				}
				errorIds = append(errorIds, id)
				continue
			}
			to, toID = "cross_id", fmt.Sprintf("%d", cross.ID)
		} else {
			if post := parser.GetPost(); post != "" && !fromCalendar {
				err := w.platform.BotPostConversation(parser.from.Address, post, parser.Date(), parser.addrList, to, toID)
				if err != nil {
					errorIds = append(errorIds, id)
					continue
				}
			}

			cross.Title = ""
			cross.Description = ""
			if !parser.HasICS() {
				cross.Place = nil
				cross.Time = nil
			}
			if cross.Place != nil || cross.Time != nil || len(cross.Exfee.Invitations) != 0 {
				err = w.platform.BotCrossUpdate(to, toID, cross, cross.By)
				if err != nil {
					if warning, ok := err.(broker.Warning); ok {
						logger.ERROR("%s can't update %s %s: %s", parser.from.Address, to, toID, warning)
					}
					errorIds = append(errorIds, id)
					continue
				}
			}
		}
		if to == "cross_id" {
			err = w.saver.Save(parser.GetIDs(), toID)
			if err != nil {
				logger.ERROR("saver save %s failed: %s", id, err)
			}
		}
		okIds = append(okIds, id)
	}
	if err := w.copy(imapConn, okIds, "posted"); err != nil {
		logger.ERROR("can't copy %v to posted: %s", errorIds, err)
	}
	if err := w.copy(imapConn, errorIds, "error"); err != nil {
		logger.ERROR("can't copy %v to error: %s", errorIds, err)
	}
	if err := w.delete(imapConn, ids); err != nil {
		logger.ERROR("can't remove %v from inbox: %s", ids, err)
	}
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
	logger.DEBUG("caps: %v", conn.Caps)
	if conn.Caps["STARTTLS"] {
		conn.StartTLS(nil)
	}

	if conn.State() == imap.Login {
		logger.DEBUG("user: %#v, password: %#v", w.config.Bot.Email.IMAPUser, w.config.Bot.Email.IMAPPassword)
		_, err = conn.Login(w.config.Bot.Email.IMAPUser, w.config.Bot.Email.IMAPPassword)
		logger.DEBUG("login err: %s", err)
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

func (w *Worker) sendHelp(err error, parser *Parser) error {
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
		logger.ERROR("template(conversation_reply.email) failed: %s", err)
	}

	to := model.Recipient{
		Provider:         "email",
		Name:             parser.from.Name,
		ExternalID:       parser.from.Address,
		ExternalUsername: parser.from.Address,
	}
	_, _, _, err = w.platform.Send(to, buf.String())
	return err
}
