package model

import (
	"fmt"
)

type Recipient struct {
	IdentityID       uint64 `json:"identity_id"`
	UserID           int64  `json:"user_id"`
	Name             string `json:"name"`
	AuthData         string `json:"auth_data"`
	Timezone         string `json:"timezone"`
	Token            string `json:"token"`
	Language         string `json:"language"`
	Provider         string `json:"provider"`
	ExternalID       string `json:"external_id"`
	ExternalUsername string `json:"external_username"`
}

func (r Recipient) Equal(other *Recipient) bool {
	return r.IdentityID == other.IdentityID && r.UserID == other.UserID
}

func (r Recipient) SameUser(other *Identity) bool {
	if r.UserID == other.UserID {
		return true
	}
	if r.IdentityID == other.ID {
		return true
	}
	if r.Provider == other.Provider {
		if r.ExternalID == other.ExternalID {
			return true
		}
		if r.ExternalUsername == other.ExternalUsername {
			return true
		}
	}
	return false
}

func (r Recipient) ID() string {
	if r.ExternalID != "" {
		return fmt.Sprintf("%s@%s", r.ExternalID, r.Provider)
	}
	return fmt.Sprintf("%s@%s", r.ExternalUsername, r.Provider)
}

func (r Recipient) String() string {
	return fmt.Sprintf("Recipient:%s(%s)@%s(i%d/u%d/t%s)", r.ExternalUsername, r.ExternalID, r.Provider, r.IdentityID, r.UserID, r.Token)
}
