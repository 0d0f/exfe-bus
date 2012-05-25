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
	Avatar_updated_at string
}

type Invitation struct {
	Id uint64
	Type string
	Token string
	Host bool
	Identity Identity
	Rsvp_status string
	By_identity Identity
	Via string
}
