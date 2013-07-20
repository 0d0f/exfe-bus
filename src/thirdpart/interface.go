package thirdpart

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"model"
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

type Sender interface {
	Provider() string
	Send(to *model.Recipient, text string) (id string, err error)
}

type Updater interface {
	Provider() string
	UpdateFriends(to *model.Recipient) error
	UpdateIdentity(to *model.Recipient) error
}

type Photographer interface {
	Provider() string
	Grab(to model.Recipient, albumID string) ([]model.Photo, error)
	Get(to model.Recipient, pictures []string) ([]string, error)
}

type Helper interface {
	UpdateIdentity(to *model.Recipient, externalUser ExternalUser) error
	UpdateFriends(to *model.Recipient, externalUsers []ExternalUser) error
	SendEmail(to string, content string) (id string, err error)
}

type FakeHelper struct {
	log *logger.Logger
}

func (h *FakeHelper) Log() *logger.Logger {
	return h.log
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
	url := fmt.Sprintf("/v3/bus/addfriends")
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

	url := fmt.Sprintf("/v3/bus/updateidentity")
	fmt.Println("url:", url)
	fmt.Println("post:", params.Encode())
	return nil
}

func (h *FakeHelper) SendEmail(to string, content string) (id string, err error) {
	fmt.Printf("send mail to %s, content: %s\n", to, content)
	return "", nil
}
