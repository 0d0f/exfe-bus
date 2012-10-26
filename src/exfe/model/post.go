package exfe_model

import (
	"fmt"
	"time"
)

type Post struct {
	Id            uint64
	By_identity   Identity
	Content       string
	Postable_id   uint64
	Postable_type string
	Via           string
	Created_at    string
	/*Relative map[string]string*/
}

func (p *Post) CreatedAt(timezone string) (string, error) {
	t, err := time.Parse("2006-1-2 15:4:5", p.Created_at[:19])
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
