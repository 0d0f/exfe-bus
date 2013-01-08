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

type Conversation struct {
	conversation *conversation.Conversation
}

func NewConversation_(config *model.Config, db *broker.DBMultiplexer, redis *broker.RedisMultiplexer, dispatcher *gobus.Dispatcher) (*Conversation, error) {
	repo, err := NewPostRepository(config, db, redis, dispatcher)
	if err != nil {
		return nil, err
	}
	return &Conversation{
		conversation: conversation.New(repo),
	}, nil
}

func (c *Conversation) SetRoute(route gobus.RouteCreater) {
	json := new(gobus.JSON)
	route().Methods("POST").Path("/cross/{cross_id}/conversation").HandlerFunc(gobus.Must(gobus.Method(json, c, "Create")))
	route().Methods("GET").Path("/cross/{cross_id}/conversation").HandlerFunc(gobus.Must(gobus.Method(json, c, "Find")))
	route().Methods("DELETE").Path("/cross/{cross_id}/conversation/{post_id}").HandlerFunc(gobus.Must(gobus.Method(json, c, "Delete")))
	route().Methods("GET").Path("/cross/{cross_id}/user/{user_id}/unread_count").HandlerFunc(gobus.Must(gobus.Method(json, c, "Unread")))
}

// 发一条新的Post到cross_id
//
// 例子：
//
//     > curl http://panda.0d0f.com:23333/cross/100354/conversation?via=web&created_at=1293479544 -d '{"by_identity":{"id":572},"content":"@googollee@twitter blablabla"}'
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
func (c *Conversation) Create(params map[string]string, arg model.Post) (model.Post, error) {
	via := params["via"]
	createdAt_ := params["created_at"]
	createdAt := time.Now().Unix()
	if createdAt_ != "" {
		var err error
		createdAt, err = strconv.ParseInt(createdAt_, 10, 64)
		if err != nil {
			return model.Post{}, fmt.Errorf("can't pasrse created_at: %s", createdAt_)
		}
	}
	crossID, err := strconv.ParseUint(params["cross_id"], 10, 64)
	if err != nil {
		return model.Post{}, fmt.Errorf("can't parse cross_id: %s", params["cross_id"])
	}
	ret, err := c.conversation.NewPost(crossID, arg, via, createdAt)
	return ret, err
}

// 查询Posts
//
// 例子：
//
//     > curl "http://panda.0d0f.com:23333/cross/100354/conversation?clear_user=378&since=2010-12-27+19:52:24&until=2010-12-27+19:52:24&min=11&max=11"
//
// 返回：
//
//     [{"id":11,"by_identity":{"id":572,"name":"Googol","connected_user_id":-572,"avatar_filename":"http://api.panda.0d0f.com/v2/avatar/default?name=Googol","provider":"email","external_id":"googollee@163.com","external_username":"googollee@163.com"},"content":"@googollee@twitter blablabla","via":"web","created_at":"2010-12-27 19:52:24 +0000","relationship":[{"uri":"identity://573","relation":"mention"}],"exfee_id":110220,"ref_uri":"cross://100354"}]
func (c *Conversation) Find(params map[string]string) ([]model.Post, error) {
	crossID, err := strconv.ParseUint(params["cross_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse cross_id: %s", params["cross_id"])
	}
	clearUserID, err := strconv.ParseInt(params["clear_user"], 10, 64)
	if err != nil {
		if params["clear_user"] == "" {
			clearUserID = 0
		} else {
			return nil, fmt.Errorf("can't parse clear_user: %s", params["clear_user"])
		}
	}
	sinceTime := params["since"]
	untilTime := params["until"]
	minID, err := strconv.ParseUint(params["min"], 10, 64)
	if err != nil {
		if params["min"] == "" {
			minID = 0
		} else {
			return nil, fmt.Errorf("can't parse min: %s", params["min"])
		}
	}
	maxID, err := strconv.ParseUint(params["max"], 10, 64)
	if err != nil {
		if params["max"] == "" {
			maxID = 0
		} else {
			return nil, fmt.Errorf("can't parse max: %s", params["max"])
		}
	}

	ret, err := c.conversation.FindPosts(crossID, clearUserID, sinceTime, untilTime, minID, maxID)
	return ret, err
}

// 删除一条Post
//
// 例子：
//
//     > curl "http://panda.0d0f.com:23333/cross/100354/conversation/11" -X DELETE
//
// 返回：
//
//     {"id":11,"by_identity":{"id":572,"name":"Googol","connected_user_id":-572,"avatar_filename":"http://api.panda.0d0f.com/v2/avatar/default?name=Googol","provider":"email","external_id":"googollee@163.com","external_username":"googollee@163.com"},"content":"@googollee@twitter blablabla","via":"web","created_at":"2010-12-27 19:52:24 +0000","relationship":[{"uri":"identity://573","relation":"mention"}],"exfee_id":110220,"ref_uri":"cross://100354"}
func (c *Conversation) Delete(params map[string]string) (model.Post, error) {
	crossID, err := strconv.ParseUint(params["cross_id"], 10, 64)
	if err != nil {
		return model.Post{}, fmt.Errorf("can't parse cross_id: %s", params["cross_id"])
	}
	postID, err := strconv.ParseUint(params["post_id"], 10, 64)
	if err != nil {
		return model.Post{}, fmt.Errorf("can't parse post_id: %s", params["post_id"])
	}

	posts, err := c.conversation.FindPosts(uint64(crossID), 0, "", "", postID, postID)
	if err != nil {
		return model.Post{}, err
	}
	if len(posts) == 0 {
		return model.Post{}, fmt.Errorf("can't find post with id %d", postID)
	}
	err = c.conversation.DeletePost(crossID, postID)
	if err != nil {
		return model.Post{}, err
	}
	ret := posts[0]
	return ret, nil
}

// 取得用户user_id未读的post条数
//
// 例子：
//
//     > curl "http://panda.0d0f.com:23333/cross/100354/user/-572/unread_count"
//
// 返回：
//
//     1
func (c *Conversation) Unread(params map[string]string) (int, error) {
	crossID, err := strconv.ParseInt(params["cross_id"], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse cross_id: %s", params["cross_id"])
	}
	userID, err := strconv.ParseInt(params["user_id"], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse user_id: %s", params["user_id"])
	}

	ret, err := c.conversation.GetUnreadCount(crossID, userID)
	return ret, err
}
