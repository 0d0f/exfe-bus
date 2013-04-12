package model

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Identity struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	UserID   int64  `json:"connected_user_id,omitempty"`
	Avatar   string `json:"avatar_filename,omitempty"`

	Provider         string `json:"provider,omitempty"`
	ExternalID       string `json:"external_id,omitempty"`
	ExternalUsername string `json:"external_username,omitempty"`
	OAuthToken       string `json:"oauth_token,omitempty"`
}

func (i Identity) Equal(other Identity) bool {
	if i.ID == other.ID {
		return true
	}
	if i.Provider == other.Provider {
		if i.ExternalID == other.ExternalID {
			return true
		}
		if i.ExternalUsername == other.ExternalUsername {
			return true
		}
	}
	return false
}

func (i Identity) SameUser(other Identity) bool {
	if i.Equal(other) {
		return true
	}
	return i.UserID == other.UserID
}

func (i Identity) String() string {
	return fmt.Sprintf("Identity:(i%d/u%d)", i.ID, i.UserID)
}

type RsvpType string

const (
	RsvpNoresponse   RsvpType = "NORESPONSE"
	RsvpAccepted              = "ACCEPTED"
	RsvpInterested            = "INTERESTED"
	RsvpDeclined              = "DECLINED"
	RsvpRemoved               = "REMOVED"
	RsvpNotification          = "NOTIFICATION"
)

type Invitation struct {
	ID         uint64   `json:"id,omitempty"`
	Host       bool     `json:"host,omitempty"`
	Mates      uint64   `json:"mates,omitempty"`
	Identity   Identity `json:"identity,omitempty"`
	RsvpStatus RsvpType `json:"rsvp_status,omitempty"`
	By         Identity `json:"by_identity,omitempty"`
	UpdateBy   Identity `json:"update_by,omitempty"`
	Via        string   `json:"via,omitempty"`
}

func (i *Invitation) String() string {
	return i.Identity.Name
}

func (i Invitation) IsAccepted() bool {
	return i.RsvpStatus == RsvpAccepted
}

func (i Invitation) IsDeclined() bool {
	return i.RsvpStatus == RsvpDeclined
}

func (i Invitation) IsPending() bool {
	return !i.IsAccepted() && !i.IsDeclined()
}

func (i Invitation) IsUpdateBy(userId int64) bool {
	return i.UpdateBy.UserID == userId
}

func (i Invitation) Avatar() string {
	resp, err := http.Get(i.Identity.Avatar)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(b)
}

type OAuthToken struct {
	Token  string `json:"oauth_token"`
	Secret string `json:"oauth_token_secret"`
}
