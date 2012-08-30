package exfe_model

import (
	"fmt"
)

type Identity struct {
	Id                int64  `json:"id"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	Nickname          string `json:"nickname"`
	Bio               string `json:"bio"`
	Provider          string `json:"provider"`
	Timezone          string `json:"timezone"`
	Connected_user_id int64  `json:"connected_user_id"`

	External_id       string `json:"external_id"`
	External_username string `json:"external_username"`
	Avatar_filename   string `json:"avatar_filename"`
}

type Invitation struct {
	Id          uint64
	Type        string
	Token       string
	Host        bool
	Mates       uint64
	Identity    Identity
	Rsvp_status string
	By_identity Identity
	Via         string
}

func (i Invitation) IsAccepted() bool {
	return i.Rsvp_status == "ACCEPTED"
}

func (i Identity) UserId() string {
	return fmt.Sprintf("%d", i.Connected_user_id)
}

func (i Identity) ExternalId() string {
	return fmt.Sprintf("%s@%s", i.External_username, i.Provider)
}

func (i Identity) DiffId() string {
	if i.Connected_user_id != 0 {
		return i.UserId()
	}
	return i.ExternalId()
}
