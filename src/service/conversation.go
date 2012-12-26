package main

import (
	"broker"
	"conversation"
	"gobus"
	"model"
)

type Conversation_ struct {
	conversation *conversation.Conversation
}

func NewConversation_(config *model.Config, db *broker.DBMultiplexer, redis *broker.RedisMultiplexer, dispatcher *gobus.Dispatcher) (*Conversation_, error) {
	repo, err := NewPostRepository(config, db, redis, dispatcher)
	if err != nil {
		return nil, err
	}
	return &Conversation_{
		conversation: conversation.New(repo),
	}, nil
}
