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

// 发一条新的Post到cross_id
//
// 例子：
//
//     > curl http://panda.0d0f.com:23333/cross/100354/Conversation?via=web&created_at=1293479544 -d '{"by_identity":{"id":572},"content":"@googollee@twitter blablabla"}'
//
// 返回：
//
//     {"id":11,"by_identity":{"id":572,"name":"Googol","connected_user_id":-572,"avatar_filename":"http://api.panda.0d0f.com/v2/avatar/default?name=Googol","provider":"email","external_id":"googollee@163.com","external_username":"googollee@163.com"},"content":"@googollee@twitter blablabla","via":"web","created_at":"2010-12-27 19:52:24 +0000","relationship":[{"uri":"identity://573","relation":"mention"}],"exfee_id":110220,"ref_uri":"cross://100354"}
//
// content解析方式： 
//
// 默认：
//
//     "@exfe@twitter look at this image http://instagr.am/xxxx\n cool!"
//      =>
//     "@exfe@twitter look at this image {{url:http://instagr.am/xxxx}}\n cool!"
//     relationship: [{"mention": "identity://123"}, {"url":"http://instagr.am/xxxx"}]
//
//特殊格式解析：
//
//     "@exfe@twitter look at this image {{image:http://instagr.am/xxxx.jpg}}\n cool!"
//      =>
//     "@exfe@twitter look at this image {{image:http://instagr.am/xxxx.jpg}}\n cool!"
//     relationship: [{"mention": "identity://123"}, {"image":"http://instagr.am/xxxx.jpg"}]
//
//     "@exfe@twitter look at this image {{webpage:http://instagr.am/xxxx}}\n cool!"
//      =>
//     "@exfe@twitter look at this image {{webpage:http://instagr.am/xxxx}}\n cool!"
//     relationship: [{"mention": "identity://123"}, {"webpage":"http://instagr.am/xxxx"}]
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

// 查询Posts
//
// 例子：
//
//     > curl "http://panda.0d0f.com:23333/cross/100354/Conversation?method=GET&clear_user=378&since=2010-12-27+19:52:24&until=2010-12-27+19:52:24&min=11&max=11" -d '""'
//
// 返回：
//
//     [{"id":11,"by_identity":{"id":572,"name":"Googol","connected_user_id":-572,"avatar_filename":"http://api.panda.0d0f.com/v2/avatar/default?name=Googol","provider":"email","external_id":"googollee@163.com","external_username":"googollee@163.com"},"content":"@googollee@twitter blablabla","via":"web","created_at":"2010-12-27 19:52:24 +0000","relationship":[{"uri":"identity://573","relation":"mention"}],"exfee_id":110220,"ref_uri":"cross://100354"}]
func (c *Conversation_) GET(meta *gobus.HTTPMeta, arg string, reply *[]model.Post) error {
	values := meta.Request.URL.Query()
	crossID, err := strconv.ParseUint(meta.Vars["cross_id"], 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse cross_id: %s", meta.Vars["cross_id"])
	}
	clearUserID, err := strconv.ParseInt(values.Get("clear_user"), 10, 64)
	if err != nil {
		if values.Get("clear_user") == "" {
			clearUserID = 0
		} else {
			return fmt.Errorf("can't parse clear_user: %s", values.Get("clear_user"))
		}
	}
	sinceTime := values.Get("since")
	untilTime := values.Get("until")
	minID, err := strconv.ParseUint(values.Get("min"), 10, 64)
	if err != nil {
		if values.Get("min") == "" {
			minID = 0
		} else {
			return fmt.Errorf("can't parse min: %s", values.Get("min"))
		}
	}
	maxID, err := strconv.ParseUint(values.Get("max"), 10, 64)
	if err != nil {
		if values.Get("max") == "" {
			maxID = 0
		} else {
			return fmt.Errorf("can't parse max: %s", values.Get("max"))
		}
	}

	*reply, err = c.conversation.FindPosts(crossID, clearUserID, sinceTime, untilTime, minID, maxID)
	return err
}

// 删除一条Post
//
// 例子：
//
//     > curl "http://panda.0d0f.com:23333/cross/100354/Conversation/11?method=DELETE" -d '""'
//
// 返回：
//
//     {"id":11,"by_identity":{"id":572,"name":"Googol","connected_user_id":-572,"avatar_filename":"http://api.panda.0d0f.com/v2/avatar/default?name=Googol","provider":"email","external_id":"googollee@163.com","external_username":"googollee@163.com"},"content":"@googollee@twitter blablabla","via":"web","created_at":"2010-12-27 19:52:24 +0000","relationship":[{"uri":"identity://573","relation":"mention"}],"exfee_id":110220,"ref_uri":"cross://100354"}
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

// 取得用户user_id未读的post条数
//
// 例子：
//
//     > curl "http://panda.0d0f.com:23333/cross/100354/user/-572/unread_count?method=Unread" -d '""'
//
// 返回：
//
//     1
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
