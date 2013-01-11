package model

type Message struct {
	Service    string      `json:"service"`
	Ticket     string      `json:"ticket"`
	Recipients []Recipient `json:"recipients"`
	Data       interface{} `json:"data"`
}
