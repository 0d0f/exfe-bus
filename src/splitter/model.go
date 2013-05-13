package splitter

import (
	"model"
)

type Pack struct {
	MergeKey string
	Method   string
	Service  string
	Update   string
	Ontime   int64
	Data     map[string]interface{}
}

type BigPack struct {
	Recipients []model.Recipient      `json:"recipients"`
	MergeKey   string                 `json:"merge_key"`
	Method     string                 `json:"method"`
	Service    string                 `json:"service"`
	Update     string                 `json:"update"`
	Ontime     int64                  `json:"ontime"`
	Data       map[string]interface{} `json:"data"`
}
