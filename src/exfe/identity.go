package exfe

type Identity struct {
	Id uint64
	Type string
	Name string
	Nickname string
	Bio string
	Provider string
	Timezone string

	External_id string
	External_username string
	avatar_filename string
	avatar_updated_at string
}

 type Invitation struct {
	 Id uint64
	 Type string
	 Identity Identity
	 Rsvp_status uint64
	 By_identity Identity
	 Via string
 }
