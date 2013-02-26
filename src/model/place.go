package model

import (
	"fmt"
)

type Place struct {
	ID          uint64 `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Lng         string `json:"lng,omitempty"`
	Lat         string `json:"lat,omitempty"`
	Provider    string `json:"provider,omitempty"`
	ExternalID  string `json:"external_id,omitempty"`
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
