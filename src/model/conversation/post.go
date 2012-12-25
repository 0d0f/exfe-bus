package conversation

import (
	"model"
)

type Post struct {
	model.Meta
	ID      uint64
	Content string `json:"content"`
	Via     string `json:"via"`
	ExfeeID uint64 `json:"exfee_id"`
	RefURI  string `json:"ref_uri"`
}

func (p *Post) ToPost() model.Post {
	return model.Post{
		ID:           p.ID,
		By:           p.By,
		Content:      p.Content,
		Via:          p.Via,
		CreatedAt:    p.CreatedAt.UTC().Format("2006-01-02 15:04:05 -0700"),
		Relationship: p.Relationship,
		ExfeeID:      p.ExfeeID,
		RefURI:       p.RefURI,
	}
}
