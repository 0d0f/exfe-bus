package model

import (
	"fmt"
)

type Recipient struct {
	ExternalID       string `json:"external_id"`
	ExternalUsername string `json:"external_username"`
	AuthData         string `json:"auth_data"`
	Provider         string `json:"provider"`
	IdentityID       int64  `json:"identity_id"`
	UserID           uint64 `json:"user_id"`
	Token            string `json:"token"`
}

func (r Recipient) String() string {
	return fmt.Sprintf("Recipient:%s(%s)@%s(i%d/u%d/t%s)", r.ExternalUsername, r.ExternalID, r.Provider, r.IdentityID, r.UserID, r.Token)
}
