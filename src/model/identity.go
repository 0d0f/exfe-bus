package model

import (
	"fmt"
	"strings"
)

type Identity struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	Locale   string `json:"locale,omitempty"`
	UserID   int64  `json:"connected_user_id,omitempty"`
	Avatar   string `json:"avatar_filename,omitempty"`

	Provider         string `json:"provider,omitempty"`
	ExternalID       string `json:"external_id,omitempty"`
	ExternalUsername string `json:"external_username,omitempty"`
	OAuthToken       string `json:"oauth_token,omitempty"`
}

func FromIdentityId(id string) Identity {
	spliter := strings.LastIndex(id, "@")
	if spliter < 0 {
		return Identity{
			ExternalUsername: id,
			Provider:         "",
		}
	}
	return Identity{
		ExternalUsername: id[:spliter],
		Provider:         id[spliter+1:],
	}
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

func (i Identity) ScreenId() string {
	switch i.Provider {
	case "email":
		return i.ExternalUsername
	case "phone":
		return i.ExternalUsername
	case "twitter":
		return "@" + i.ExternalUsername + "@" + i.Provider
	case "wechat":
		return "@" + i.Provider
	}
	return i.ExternalUsername + "@" + i.Provider
}

func (i Identity) Id() string {
	return fmt.Sprintf("%s@%s", i.ExternalUsername, i.Provider)
}

func (i Identity) ToRecipient() Recipient {
	return Recipient{
		IdentityID:       i.ID,
		UserID:           i.UserID,
		Timezone:         i.Timezone,
		Language:         i.Locale,
		Provider:         i.Provider,
		ExternalID:       i.ExternalID,
		ExternalUsername: i.ExternalUsername,
	}
}

type RsvpType string

const (
	Noresponse   RsvpType = "NORESPONSE"
	Accepted              = "ACCEPTED"
	Interested            = "INTERESTED"
	Declined              = "DECLINED"
	Removed               = "REMOVED"
	Notification          = "NOTIFICATION"
)

type Invitation struct {
	ID            uint64   `json:"id,omitempty"`
	Host          bool     `json:"host,omitempty"`
	Mates         uint64   `json:"mates,omitempty"`
	Identity      Identity `json:"identity,omitempty"`
	Response      RsvpType `json:"response,omitempty"`
	By            Identity `json:"by_identity,omitempty"`
	UpdatedBy     Identity `json:"updated_by,omitempty"`
	Via           string   `json:"via,omitempty"`
	Token         string   `json:"token,omitempty"`
	Notifications []string `json:"notification_identities,omitempty"`
}

func (i *Invitation) String() string {
	return i.Identity.Name
}

func (i Invitation) IsAccepted() bool {
	return i.Response == Accepted
}

func (i Invitation) IsDeclined() bool {
	return i.Response == Declined
}

func (i Invitation) IsPending() bool {
	return !i.IsAccepted() && !i.IsDeclined()
}

func (i Invitation) IsUpdatedBy(userId int64) bool {
	return i.UpdatedBy.UserID == userId
}

type OAuthToken struct {
	Token  string `json:"oauth_token"`
	Secret string `json:"oauth_token_secret"`
}
