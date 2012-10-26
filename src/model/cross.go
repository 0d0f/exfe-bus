package model

import (
	"fmt"
)

type Cross struct {
	ID          uint64    `json:"id"`
	By          Identity  `json:"by_identity"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Time        CrossTime `json:"time"`
	Place       Place     `json:"place"`
	Exfee       Exfee     `json:"exfee"`
}

func (c Cross) Equal(other *Cross) bool {
	return c.ID == other.ID
}

func (c Cross) String() string {
	return fmt.Sprintf("Cross:%d", c.ID)
}

type CrossUpdate struct {
	To       Recipient `json:"to"`
	OldCross Cross     `json:"old_cross"`
	Cross    Cross     `json:"cross"`
	By       Identity  `json:"by"`
}
