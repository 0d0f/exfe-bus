package exfe_model

import (
	"fmt"
)

type Identity struct {
	Id                uint64
	Type              string
	Name              string
	Nickname          string
	Bio               string
	Provider          string
	Timezone          string
	Connected_user_id int64

	External_id       string
	External_username string
	Avatar_filename   string
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
