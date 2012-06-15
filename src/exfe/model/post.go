package exfe_model

import (
	"time"
	"fmt"
)

type Post struct {
	Id uint64
	By_identity Identity
	Content string
	Postable_id uint64
	Postable_type string
	Via string
	Created_at string
	/*Relative map[string]string*/
}

func (p *Post) CreatedAt(timezone string) (string, error) {
	t, err := time.Parse("2006-01-02 15:04:05 -0700", p.Create_at)
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
