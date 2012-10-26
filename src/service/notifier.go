package main

import (
	"formatter"
	"gobus"
	"model"
	"notifier"
	"service/args"
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
func (c *Conversation) Update(meta *gobus.HTTPMeta, updates args.ConversationUpdateArg, i *int) error {
	*i = 0
	return c.conversation.Update(updates)
}
