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

type Head struct {
	name   string
	repo   *delayrepo.Head
	config *model.Config
}

func NewHead(delayInSecond uint, config *model.Config) (*Head, *tomb.Tomb) {
	name := fmt.Sprintf("delayrepo:head_%ds", delayInSecond)
	delay := delayInSecond
	redis := broker.NewRedisMultiplexer(config)
	repo := delayrepo.NewHead(name, delay, redis)
	log := config.Log.SubPrefix(name)
	tomb := delayrepo.ServRepository(log, repo, getCallback(log, config))

	return &Head{
		name:   name,
		repo:   repo,
		config: config,
	}, tomb
}

// 首延迟发送队列
//
// 例子：
//
// > curl 'http://127.0.0.1:23334/Head10?method=Push' -d '{"service":"Conversation",
//       "method":"Update",
//       "merge_key":"email_cross123",
//       "tos":[{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"}],
//       "data":{"to":{"identity_id":12,"user_id":3,"name":"email1 name","auth_data":"","timezone":"+0800","token":"recipient_email1_token","language":"en_US","provider":"email","external_id":"sender1@gmail.com","external_username":"sender1@gmail.com"},"cross":{"id":123,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"},"title":"Test Cross","description":"test cross description","time":{"begin_at":{"date_word":"","date":"","time_word":"","time":"","timezone":""},"origin":"","output_format":0},"place":{"id":0,"title":"","description":"","lng":"","lat":"","provider":"","external_id":""},"exfee":{"id":0,"name":"","invitations":null}},"post":{"id":1,"by_identity":{"id":11,"name":"email1 name","nickname":"email1 nick","bio":"email1 bio","timezone":"+0800","connected_user_id":1,"avatar_filename":"http://path/to/email1.avatar","provider":"email","external_id":"sender2@hotmail.com","external_username":"sender2@hotmail.com"},"content":"email1 post sth","via":"abc","created_at":"2012-10-24 16:31:00"}}}'
func (i *Head) Push(meta *gobus.HTTPMeta, arg model.QueuePush, count *int) error {
	*count = 0
	if len(arg.Tos) == 0 {
		*count = 1
		data, err := json.Marshal(arg.Data)
		if err != nil {
			return fmt.Errorf("can't marshal input data: %s", err)
		}
		err = i.repo.Push(fmt.Sprintf("%s,%s,%s,-", arg.Service, arg.Method, arg.MergeKey), data)
		if err != nil {
			return fmt.Errorf("push to repo failed: %s", err)
		}
		return nil
	}

	data, ok := arg.Data.(map[string]interface{})
	for _, to := range arg.Tos {
		if ok {
			data["to"] = to
		}
		d, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("can't marshal input data: %s", err)
		}
		err = i.repo.Push(fmt.Sprintf("%s,%s,%s,%s(%s)@%s", arg.Service, arg.Method, arg.MergeKey, to.ExternalID, to.ExternalUsername, to.Provider), d)
		if err != nil {
			return fmt.Errorf("push to repo failed: %s", err)
		}
		*count++
	}

	return nil
}

func (i Head) String() string {
	return i.name
}
