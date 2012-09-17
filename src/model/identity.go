package model

type Identity struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Bio      string `json:"bio"`
	Timezone string `json:"timezone"`
	UserID   uint64 `json:"connected_user_id"`
	Avatar   string `json:"avatar_filename"`

	Provider         string `json:"provider"`
	ExternalID       string `json:"external_id"`
	ExternalUsername string `json:"external_username"`
	OAuthToken       string `json:"oauth_token"`
}

func (i Identity) IsSame(other *Identity) bool {
	return i.ID == other.ID || (i.Provider == other.Provider && i.ExternalID == other.ExternalID)
}
