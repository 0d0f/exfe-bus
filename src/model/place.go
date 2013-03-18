package model

import (
	"fmt"
)

type Place struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Lng         string `json:"lng"`
	Lat         string `json:"lat"`
	Provider    string `json:"provider"`
	ExternalID  string `json:"external_id"`
}

func (p *Place) String() string {
	return fmt.Sprintf("Place(%d)", p.ID)
}

func (p *Place) Same(other *Place) bool {
	if p == nil || other == nil {
		return false
	}
	if p == other {
		return true
	}
	return p.Title == other.Title && p.Description == other.Description && p.Lng == other.Lng && p.Lat == other.Lat
}
