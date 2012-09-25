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
	UserID           uint64 `json:"connected_user_id"`
}

func (r Recipient) String() string {
	return fmt.Sprintf("Recipient:%s@%s", r.ExternalID, r.Provider)
}
