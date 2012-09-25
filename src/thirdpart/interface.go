package thirdpart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"model"
	"net/http"
	"net/url"
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

type MessageType int

const (
	ShortMessage MessageType = iota
	LongTextMessage
	HTMLMessage
	EmailMessage
)

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
	CrossID uint64
	Type    DataType
}

type Sender interface {
	Provider() string
	MessageType() MessageType
	Send(to *model.Recipient, privateMessage string, publicMessage string, data *InfoData) (id string, err error)
}

type Updater interface {
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
	config *model.Config
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
	return "", fmt.Errorf("not implemented")
}

type HelperFake struct {
}

func (h *HelperFake) UpdateFriends(to *model.Recipient, externalUsers []ExternalUser) error {
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

func (h *HelperFake) UpdateIdentity(to *model.Recipient, externalUser ExternalUser) error {
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

func (h *HelperFake) SendEmail(to string, content string) (id string, err error) {
	fmt.Printf("send mail to %s, content: %s\n", to, content)
	return "", nil
}
