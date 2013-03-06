package main

import (
	"fmt"
)

type LoadType int

const (
	Ping LoadType = iota
	Pong
	Send
	Respond
)

type Load struct {
	Type    LoadType `json:"type"`
	To      string   `json:"to"`
	Content string   `json:"content"`
}

func (l Load) String() string {
	return fmt.Sprintf("{type:%d to:%s}", l.Type, l.To)
}
