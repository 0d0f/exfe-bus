package exfe_model

type Identity struct {
	Id uint64
	Type string
	Token string
	Name string
	Nickname string
	Bio string
	Provider string
	Timezone string
	Connected_user_id uint64

	External_id string
	External_username string
	Avatar_filename string
}

type Invitation struct {
	Id uint64
	Type string
	Token string
	Host bool
	Mates uint64
	Identity Identity
	Rsvp_status string
	By_identity Identity
	Via string
}

func (i Invitation) IsAccepted() bool {
	return i.Rsvp_status == "ACCEPTED"
}
