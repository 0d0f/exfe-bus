package model

import (
	"fmt"
	"time"
)

type Post struct {
	ID           uint64   `json:"id"`
	By           Identity `json:"by_identity"`
	Content      string   `json:"content"`
	PostableID   uint64   `json:"postable_id"`
	PostableType string   `json:"postable_type"`
	Via          string   `json:"via"`
	CreatedAt    string   `json:"created_at"`
	/*Relative map[string]string*/
}

func (p *Post) CreatedAtInZone(timezone string) (string, error) {
	createdAt := p.CreatedAt
	if len(createdAt) > 19 {
		createdAt = createdAt[:19]
	}
	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return "", err
	}
	loc, err := LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("Parse target zone error: %s", err)
	}
	t = t.In(loc)
	return t.Format("03:04PM Mon, Jan 2"), nil
}