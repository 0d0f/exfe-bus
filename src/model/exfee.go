package model

type Exfee struct {
	ID          uint64       `json:"id"`
	Invitations []Invitation `json:"invitations"`
}
