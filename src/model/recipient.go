package model

import (
	"fmt"
	"strings"
)

type IdentityId string

func (i IdentityId) Split() (externalId, provider string, err error) {
	s := string(i)
	spliter := strings.LastIndex(s, "@")
	if spliter < 0 {
		err = fmt.Errorf("invalid identity id: %s", i)
		return
	}
	externalId = s[:spliter]
	provider = s[spliter+1:]
	return
}

type Recipient struct {
	IdentityID       int64  `json:"identity_id"`
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

func (r Recipient) Tunnel() Recipient {
	switch r.Provider {
	case "facebook":
		r.ExternalUsername = r.ExternalUsername + "@facebook.com"
		r.Provider = "email"
	case "google":
		r.Provider = "email"
	}
	return r
}

func (r Recipient) Equal(other *Recipient) bool {
	if r.UserID == other.UserID {
		return true
	}
	if r.IdentityID == other.IdentityID {
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
	return fmt.Sprintf("%s@%s(u%d)", r.ExternalUsername, r.Provider, r.UserID)
}
