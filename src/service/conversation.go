package main

import (
	"broker"
	"conversation"
	"fmt"
	"gobus"
	"model"
	"strconv"
	"time"
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

func (c *Conversation_) POST(meta *gobus.HTTPMeta, arg model.Post, reply *model.Post) error {
	values := meta.Request.URL.Query()
	via := values.Get("via")
	createdAt_ := values.Get("created_at")
	createdAt := time.Now().Unix()
	if createdAt_ != "" {
		var err error
		createdAt, err = strconv.ParseInt(createdAt_, 10, 64)
		if err != nil {
			return fmt.Errorf("can't pasrse created_at: %s", createdAt_)
		}
	}
	crossID, err := strconv.ParseUint(meta.Vars["cross_id"], 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse cross_id: %s", meta.Vars["cross_id"])
	}
	*reply, err = c.conversation.NewPost(crossID, arg, via, createdAt)
	return err
}

func (c *Conversation_) GET(meta *gobus.HTTPMeta, arg string, reply *[]model.Post) error {
	values := meta.Request.URL.Query()
	crossID, err := strconv.ParseUint(meta.Vars["cross_id"], 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse cross_id: %s", meta.Vars["cross_id"])
	}
	clearUserID, err := strconv.ParseInt(values.Get("clear_user"), 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse clear_user: %s", values.Get("clear_user"))
	}
	sinceTime := values.Get("since")
	untilTime := values.Get("until")
	minID, err := strconv.ParseUint(values.Get("min"), 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse min: %s", values.Get("min"))
	}
	maxID, err := strconv.ParseUint(values.Get("max"), 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse max: %s", values.Get("max"))
	}

	*reply, err = c.conversation.FindPosts(crossID, clearUserID, sinceTime, untilTime, minID, maxID)
	return err
}

func (c *Conversation_) DELETE(meta *gobus.HTTPMeta, arg string, reply *model.Post) error {
	crossID, err := strconv.ParseUint(meta.Vars["cross_id"], 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse cross_id: %s", meta.Vars["cross_id"])
	}
	postID, err := strconv.ParseUint(meta.Vars["post_id"], 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse post_id: %s", meta.Vars["post_id"])
	}

	posts, err := c.conversation.FindPosts(uint64(crossID), 0, "", "", postID, postID)
	if err != nil {
		return err
	}
	if len(posts) == 0 {
		return fmt.Errorf("can't find post with id %d", postID)
	}
	err = c.conversation.DeletePost(crossID, postID)
	if err != nil {
		return err
	}
	*reply = posts[0]
	return nil
}

func (c *Conversation_) Unread(meta *gobus.HTTPMeta, arg string, reply *int) error {
	crossID, err := strconv.ParseInt(meta.Vars["cross_id"], 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse cross_id: %s", meta.Vars["cross_id"])
	}
	userID, err := strconv.ParseInt(meta.Vars["user_id"], 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse user_id: %s", meta.Vars["user_id"])
	}

	*reply, err = c.conversation.GetUnreadCount(crossID, userID)
	return err
}
