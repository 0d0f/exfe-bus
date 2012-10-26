package thirdpart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"model"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
)

type Token struct {
	Token  string
	Secret string
}

type ExternalUser interface {
	ExternalID() string
	Provider() string
	Name() string
	ExternalUsername() string
	Bio() string
	Avatar() string
}

type MessageType string

const (
	ShortMessage    MessageType = "short_message.txt"
	LongTextMessage             = "long_message.txt"
	HTMLMessage                 = "html"
	EmailMessage                = "email"
)

func MessageTypeFromProvider(provider string) (MessageType, error) {
	switch provider {
	case "email":
		return EmailMessage, nil
	case "facebook":
		return EmailMessage, nil
	case "iOS":
		return ShortMessage, nil
	case "android":
		return ShortMessage, nil
	case "twitter":
		return ShortMessage, nil
	}
	return "", fmt.Errorf("unknow provider: %s", provider)
}

type DataType string

func (t DataType) String() string {
	return string(t)
}

const (
	CrossInvitation DataType = "i"
	CrossUpdate              = "u"
	CrossRemove              = "r"
	Conversation             = "c"
)

type InfoData struct {
	CrossID uint64   `json:"cross_id"`
	Type    DataType `json:"type"`
}

func (i InfoData) String() string {
	return fmt.Sprintf("{cross:%d type:%s}", i.CrossID, i.Type)
}

type Sender interface {
	Provider() string
	Send(to *model.Recipient, privateMessage string, publicMessage string, data *InfoData) (id string, err error)
}

type Updater interface {
	Provider() string
	UpdateFriends(to *model.Recipient) error
	UpdateIdentity(to *model.Recipient) error
}

type Helper interface {
	UpdateIdentity(to *model.Recipient, externalUser ExternalUser) error
	UpdateFriends(to *model.Recipient, externalUsers []ExternalUser) error
	SendEmail(to string, content string) (id string, err error)
}

type updateFriendsArg struct {
	UserID     uint64            `json:"user_id"`
	Identities []*model.Identity `json:"identities"`
}

type HelperImp struct {
	config      *model.Config
	emailServer string
	emailAuth   smtp.Auth
	emailFrom   string
}

func NewHelper(config *model.Config) *HelperImp {
	return &HelperImp{
		config:      config,
		emailServer: fmt.Sprintf("%s:%d", config.Email.Host, config.Email.Port),
		emailAuth:   smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.Host),
		emailFrom:   fmt.Sprintf("exfe@%s", config.Email.Domain),
	}
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
	mx, err := net.LookupMX(host)
	if err != nil {
		return "", fmt.Errorf("lookup mail exchange fail: %s", err)
	}
	if len(mx) == 0 {
		return "", fmt.Errorf("can't find mail exchange for %s", to)
	}
	s, err := smtp.Dial(fmt.Sprintf("%s:25", mx[0].Host))
	if err != nil {
		return "", fmt.Errorf("dial to mail exchange %s fail: %s", mx[0].Host, err)
	}
	err = s.Mail(h.emailFrom)
	if err != nil {
		return "", fmt.Errorf("set smtp mail fail: %s", err)
	}
	err = s.Rcpt(to)
	if err != nil {
		return "", fmt.Errorf("set smtp rcpt fail: %s", err)
	}
	s.Quit()
	err = smtp.SendMail(h.emailServer, h.emailAuth, h.emailFrom, []string{to}, []byte(content))
	if err != nil {
		return "", fmt.Errorf("mail send fail: %s", err)
	}
	return "", nil
}

type FakeHelper struct {
}

func (h *FakeHelper) UpdateFriends(to *model.Recipient, externalUsers []ExternalUser) error {
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
	url := fmt.Sprintf("/v2/Gobus/AddFriends")
	fmt.Println("url:", url)
	fmt.Println("post:", buf.String())
	return nil
}

func (h *FakeHelper) UpdateIdentity(to *model.Recipient, externalUser ExternalUser) error {
	params := make(url.Values)
	params.Set("id", fmt.Sprintf("%d", to.IdentityID))
	params.Set("provider", externalUser.Provider())
	params.Set("external_id", externalUser.ExternalID())
	params.Set("name", externalUser.Name())
	params.Set("bio", externalUser.Bio())
	params.Set("avatar_filename", externalUser.Avatar())
	params.Set("external_username", externalUser.ExternalUsername())

	url := fmt.Sprintf("/v2/gobus/UpdateIdentity")
	fmt.Println("url:", url)
	fmt.Println("post:", params.Encode())
	return nil
}

func (h *FakeHelper) SendEmail(to string, content string) (id string, err error) {
	fmt.Printf("send mail to %s, content: %s\n", to, content)
	return "", nil
}
