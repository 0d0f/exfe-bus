package model

import (
	"fmt"
)

type QueuePush struct {
	Service  string      `json:"service"`
	Method   string      `json:"method"`
	MergeKey string      `json:"merge_key"`
	Tos      []Recipient `json:"tos"` // it will expand and overwrite "to" field in data
	Data     interface{} `json:"data"`
}

func (a QueuePush) String() string {
	return fmt.Sprintf("{service:%s method:%s key:%s tos:%s}", a.Service, a.Method, a.MergeKey, a.Tos)
}
