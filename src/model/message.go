package model

type Message struct {
	Services   string      `json:"services"`
	Ticket     string      `json:"ticket"`
	Recipients []Recipient `json:"recipients"`
	Data       interface{} `json:"data"`
}
