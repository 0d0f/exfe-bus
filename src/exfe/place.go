package exfe

import (
	"fmt"
)

type Place struct {
	Id uint64
	Type string
	Title string
	Description string
	Lng string
	Lat string
	Provider string
	External_id string
}

func (p *Place) String() string {
	if p.Title == "" {
		return ""
	}

	if p.Description == "" {
		return fmt.Sprintf("at %s", p.Title)
	}

	return fmt.Sprintf("at %s, %s", p.Title, p.Description)
}
