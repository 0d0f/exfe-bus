package main

import (
	"formatter"
	"model"
	"notifier"
)

type Conversation struct {
	conversation *notifier.Conversation
}

func NewConversation(localTemplate *formatter.LocalTemplate, config *model.Config) *Conversation {
	return &Conversation{
		conversation: notifier.NewConversation(localTemplate, config),
	}
}

// 发送Conversation的更新消息updates
//
// Cross内容太长，懒得写例子了……
//
func (c *Conversation) Update(updates []model.ConversationUpdate, i *int) error {
	return c.conversation.Update(updates)
}
