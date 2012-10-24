package main

import (
	"business"
	"formatter"
	"model"
)

type Conversation struct {
	conversation *business.Conversation
}

func NewConversation(localTemplate *formatter.LocalTemplate, config *model.Config) *Conversation {
	return &Conversation{
		conversation: business.NewConversation(localTemplate, config),
	}
}

// 发送Conversation的更新消息updates
//
// Cross内容太长，懒得写例子了……
//
func (c *Conversation) Update(updates []model.ConversationUpdate) error {
	return c.conversation.Update(updates)
}