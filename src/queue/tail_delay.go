package main

import (
	"broker"
	"delayrepo"
	"encoding/json"
	"fmt"
	"gobus"
	"launchpad.net/tomb"
	"model"
)

type Tail struct {
	name   string
	repo   *delayrepo.Tail
	config *model.Config
	url    string
}

func NewTail(delayInSecond uint, url string, config *model.Config) (*Tail, *tomb.Tomb) {
	name := fmt.Sprintf("delayrepo:tail_%ds", delayInSecond)
	delay := delayInSecond
	redis := broker.NewRedisMultiplexer(config)
	repo := delayrepo.NewTail(name, delay, redis)
	log := config.Log.SubPrefix(name)
	tomb := delayrepo.ServRepository(log, repo, getCallback(log, config))

	return &Tail{
		name:   name,
		repo:   repo,
		config: config,
		url:    url,
	}, tomb
}

func (i *Tail) SetRoute(route gobus.RouteCreater) error {
	json := new(gobus.JSON)
	route().Methods("POST").Path("/"+i.url).HandlerMethod(json, i, "Push")
	return nil
}

// 尾延迟发送队列
//
// 例子：
//
// > curl 'http://127.0.0.1:23334/tail10' -d '{"service":"Conversation",
//       "method":"Update",
//       "merge_key":"email_cross123",
//       "tos":[{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"}],
//       "data":{"to":{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":""},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":0,"name":"","invitations":null}},"post":{"id":1,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"},"content":"email1 post sth","via":"abc","created_at":"2012-10-24 16:31:00"}}}'
func (i *Tail) Push(params map[string]string, arg model.QueuePush) (int, error) {
	ret := 0
	if len(arg.Tos) == 0 {
		ret = 1
		data, err := json.Marshal(arg.Data)
		if err != nil {
			return 0, fmt.Errorf("can't marshal input data: %s", err)
		}
		err = i.repo.Push(fmt.Sprintf("%s,%s,%s,-", arg.Service, arg.Method, arg.MergeKey), data)
		if err != nil {
			return 0, fmt.Errorf("push to repo failed: %s", err)
		}
		return ret, nil
	}

	data, ok := arg.Data.(map[string]interface{})
	for _, to := range arg.Tos {
		if ok {
			data["to"] = to
		}
		d, err := json.Marshal(data)
		if err != nil {
			return ret, fmt.Errorf("can't marshal input data: %s", err)
		}
		err = i.repo.Push(fmt.Sprintf("%s,%s,%s,%s(%s)@%s", arg.Service, arg.Method, arg.MergeKey, to.ExternalID, to.ExternalUsername, to.Provider), d)
		if err != nil {
			return ret, fmt.Errorf("push to repo failed: %s", err)
		}
		ret++
	}

	return ret, nil
}

func (i Tail) String() string {
	return i.name
}
