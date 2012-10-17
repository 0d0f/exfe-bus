package model

type Cross struct {
	ID          uint64    `json:"id"`
	By          Identity  `json:"by_identity"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Time        CrossTime `json:"time"`
	Place       Place     `json:"place"`
	Exfee       Exfee     `json:"exfee"`
}
