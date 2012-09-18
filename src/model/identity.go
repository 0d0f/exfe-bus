package model

type Identity struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	UserID   uint64 `json:"connected_user_id,omitempty"`
	Avatar   string `json:"avatar_filename,omitempty"`

	Provider         string `json:"provider,omitempty"`
	ExternalID       string `json:"external_id,omitempty"`
	ExternalUsername string `json:"external_username,omitempty"`
	OAuthToken       string `json:"oauth_token,omitempty"`
}

func (i Identity) IsSame(other *Identity) bool {
	return i.ID == other.ID || (i.Provider == other.Provider && i.ExternalID == other.ExternalID)
}
