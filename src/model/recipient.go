package model

import (
	"fmt"
)

type Recipient struct {
	IdentityID       int64  `json:"identity_id"`
	UserID           uint64 `json:"user_id"`
	Name             string `json:"name"`
	AuthData         string `json:"auth_data"`
	Timezone         string `json:"timezone"`
	Token            string `json:"token"`
	Language         string `json:"language"`
	Provider         string `json:"provider"`
	ExternalID       string `json:"external_id"`
	ExternalUsername string `json:"external_username"`
}

func (r Recipient) String() string {
	return fmt.Sprintf("Recipient:%s(%s)@%s(i%d/u%d/t%s)", r.ExternalUsername, r.ExternalID, r.Provider, r.IdentityID, r.UserID, r.Token)
}
