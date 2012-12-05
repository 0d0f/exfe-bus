package thirdpart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-aws/smtp"
	"github.com/googollee/go-logger"
	"github.com/googollee/go-multiplexer"
	"model"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type updateFriendsArg struct {
	UserID     int64             `json:"user_id"`
	Identities []*model.Identity `json:"identities"`
}

type HelperImp struct {
	config        *model.Config
	emailCheckers *multiplexer.Hetero
	emailSender   *multiplexer.Homo
	emailFrom     string
}

func NewHelper(config *model.Config) *HelperImp {
	auth := smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.Host)
	return &HelperImp{
		config: config,
		emailCheckers: multiplexer.NewHetero(func(key string) (multiplexer.Instance, error) {
			return NewSmtpCheckerInstance(key, config.Log)
		}, time.Duration(config.Email.IdleTimeoutInSec)*time.Second, time.Duration(config.Email.IntervalInSec)*time.Second),
		emailSender: multiplexer.NewHomo(func() (multiplexer.Instance, error) {
			return NewSmtpSenderInstance(config.Log, config.Email.Host, auth)
		}, 5, time.Duration(config.Email.IdleTimeoutInSec)*time.Second, time.Duration(config.Email.IntervalInSec)*time.Second),
		emailFrom: fmt.Sprintf("x@%s", config.Email.Domain),
	}
}

func (h *HelperImp) Log() *logger.Logger {
	return h.config.Log
}

func (h *HelperImp) UpdateFriends(to *model.Recipient, externalUsers []ExternalUser) error {
	arg := updateFriendsArg{
		UserID:     to.UserID,
		Identities: make([]*model.Identity, len(externalUsers)),
	}
	for i, u := range externalUsers {
		user := &model.Identity{
			Name:             u.Name(),
			Provider:         u.Provider(),
			ExternalID:       u.ExternalID(),
			ExternalUsername: u.ExternalUsername(),
			Bio:              u.Bio(),
			Avatar:           u.Avatar(),
		}
		arg.Identities[i] = user
	}
	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	err := e.Encode(arg)
	if err != nil {
		return fmt.Errorf("encoding user error: %s", err)
	}
	url := fmt.Sprintf("%s/v2/Gobus/AddFriends", h.config.SiteApi)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return fmt.Errorf("update %s friends fail: %s", to, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("update %s friends fail: %s", to, resp.Status)
	}
	return nil
}

func (h *HelperImp) UpdateIdentity(to *model.Recipient, externalUser ExternalUser) error {
	params := make(url.Values)
	params.Set("id", fmt.Sprintf("%d", to.IdentityID))
	params.Set("provider", externalUser.Provider())
	params.Set("external_id", externalUser.ExternalID())
	params.Set("name", externalUser.Name())
	params.Set("bio", externalUser.Bio())
	params.Set("avatar_filename", externalUser.Avatar())
	params.Set("external_username", externalUser.ExternalUsername())

	url := fmt.Sprintf("%s/v2/gobus/UpdateIdentity", h.config.SiteApi)
	resp, err := http.PostForm(url, params)
	if err != nil {
		return fmt.Errorf("update with %v failed: %s", params, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("update with %v failed: %s", params, resp.Status)
	}
	return nil
}

func (h *HelperImp) SendEmail(to string, content string) (id string, err error) {
	mail_split := strings.Split(to, "@")
	if len(mail_split) != 2 {
		return "", fmt.Errorf("mail(%s) not valid.", to)
	}
	host := mail_split[1]
	err = nil

	h.emailCheckers.Do(host, func(i multiplexer.Instance) {
		s := i.(*SmtpInstance)
		err = s.conn.Mail(h.emailFrom)
		if err != nil {
			return
		}
		err = s.conn.Rcpt(to)
		if err != nil {
			return
		}
	})
	if err != nil {
		return "", fmt.Errorf("mail check fail: %s", err)
	}

	h.emailSender.Do(func(i multiplexer.Instance) {
		c := i.(*SmtpInstance).conn
		err = c.Mail(h.emailFrom)
		if err != nil {
			return
		}
		err = c.Rcpt(to)
		if err != nil {
			return
		}
		var w *smtp.DataWriter
		w, err = c.Data()
		if err != nil {
			return
		}
		_, err = w.Write([]byte(content))
		if err != nil {
			return
		}
		err = w.Close()
		if err != nil {
			return
		}
		id = w.MessageID()
	})
	if err != nil {
		return "", fmt.Errorf("mail send fail: %s", err)
	}
	return id, nil
}
