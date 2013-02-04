package model

type Image struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}

type Photo struct {
	ID              int      `json:"id"`
	Caption         string   `json:"caption"`
	By              Identity `json:"by_identity"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	Provider        string   `json:"provider"`
	ExternalAlbumID string   `json:"external_album_id"`
	ExternalID      string   `json:"external_id"`
	Location        Place    `json:"location"`
	Images          struct {
		Fullsize Image `json:"fullsize"`
		Preview  Image `json:"preview"`
	} `json:"images"`
}
