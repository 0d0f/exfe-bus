package thirdpart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-aws/smtp"
	"github.com/googollee/go-logger"
	"model"
	"net"
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
	config    *model.Config
	emailFrom string
	emailHost string
	auth      smtp.Auth
}

func NewHelper(config *model.Config) *HelperImp {
	auth := smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.Host)
	return &HelperImp{
		config:    config,
		emailFrom: fmt.Sprintf("x@%s", config.Email.Domain),
		emailHost: config.Email.Host,
		auth:      auth,
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
	defer resp.Body.Close()
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
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("update with %v failed: %s", params, resp.Status)
	}
	return nil
}

func (h *HelperImp) SendEmail(to string, content string) (string, error) {
	mail_split := strings.Split(to, "@")
	if len(mail_split) != 2 {
		return "", fmt.Errorf("mail(%s) not valid.", to)
	}
	host := mail_split[1]
	addr := ""
	var s *smtp.Client
	var conn net.Conn

	mx, err := net.LookupMX(host)
	if err != nil {
		h.Log().Notice("lookup mail exchange fail: %s", err)
		goto SEND
	}
	if len(mx) == 0 {
		h.Log().Notice("unreach mail exchange: %s", host)
		goto SEND
	}
	addr = fmt.Sprintf("%s:25", mx[0].Host)
	conn, err = net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		h.Log().Notice("conn %s fail: %s", addr, err)
		goto SEND
	}
	conn.SetDeadline(time.Now().Add(time.Second * 10))
	s, err = smtp.NewClient(conn, host)
	if err != nil {
		h.Log().Notice("new smtp client %s fail: %s", mx[0].Host, err)
		goto SEND
	}
	err = s.Mail(h.emailFrom)
	if err != nil {
		h.Log().Notice("mail smtp %s command mail fail: %s", host, err)
		goto SEND
	}
	err = s.Rcpt(to)
	if err != nil {
		return "", fmt.Errorf("can't find mail: %s", to)
	}
	s.Quit()

SEND:
	id, err := smtp.SendMailTimeout(h.emailHost+":25", h.auth, h.emailFrom, []string{to}, []byte(content), time.Second*10)
	return id, err
}
