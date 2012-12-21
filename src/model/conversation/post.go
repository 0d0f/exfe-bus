package conversation

import (
	"model"
)

type ConversationPost struct {
	model.Meta
	Content string `json:"content"`
	Via     string `json:"via"`
	ExfeeID uint64 `json:"exfee_id"`
}
