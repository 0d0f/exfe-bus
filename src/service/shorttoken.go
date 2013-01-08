package main

import (
	"broker"
	"gobus"
	"model"
	"shorttoken"
	"time"
)

type ShortToken struct {
	short *shorttoken.ShortToken
}

func NewShortToken(config *model.Config, db *broker.DBMultiplexer) (*ShortToken, error) {
	repo, err := NewShortTokenRepository(config, db)
	if err != nil {
		return nil, err
	}
	return &ShortToken{
		short: shorttoken.New(repo, 4),
	}, nil
}

type PostArg struct {
	Resource          string `json:"resource"`
	Data              string `json:"data"`
	ExpireAfterSecond int    `json:"expire_after_second"`
}

// 根据resource，data和expire after second创建一个token
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/shorttoken" -d '{"data":"abc","resource":"123","expire_after_second":300}'
//
// 返回：
//
//     {"key":"0303","data":"abc"}
func (s *ShortToken) POST(meta *gobus.HTTPMeta, arg PostArg, reply *model.Token) error {
	after := time.Duration(arg.ExpireAfterSecond) * time.Second
	var err error
	*reply, err = s.short.Create(arg.Resource, arg.Data, after)
	return err
}

// 根据key或者resource获得一个token，如果token不存在，返回错误
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/shorttoken?method=GET&key=0303&resource=123" -d '""'
//
// 返回：
//
//     {"key":"0303","data":"abc"}
func (s *ShortToken) GET(meta *gobus.HTTPMeta, arg string, reply *[]model.Token) error {
	params := meta.Request.URL.Query()
	key := params.Get("key")
	resource := params.Get("resource")
	var err error
	*reply, err = s.short.Get(key, resource)
	return err
}

type UpdateArg struct {
	Data               *string `json:"data"`
	ExpireAfterSeconds *int    `json:"expire_after_seconds"`
}

// 更新key对应的token的data信息或者expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/shorttoken/0303?method=PUT" -d '{"data":"xyz","expire_after_seconds":13}'
//
// 返回：
//
//     0
func (s *ShortToken) PUT(meta *gobus.HTTPMeta, arg UpdateArg, reply *int) error {
	key := meta.Vars["key"]
	if arg.Data != nil {
		err := s.short.UpdateData(key, *arg.Data)
		if err != nil {
			return err
		}
	}
	if arg.ExpireAfterSeconds != nil {
		after := time.Duration(*arg.ExpireAfterSeconds) * time.Second
		err := s.short.Refresh(key, "", after)
		if err != nil {
			return err
		}
	}
	return nil
}

type ExpireArg struct {
	Resource           string `json:"resource"`
	ExpireAfterSeconds int    `json:"expire_after_seconds"`
}

// 更新resource对应的token的expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/shorttoken?method=Expire" -d '{"resource":"123","expire_after_seconds":13}'
//
// 返回：
//
//     0
func (s *ShortToken) Expire(meta *gobus.HTTPMeta, arg ExpireArg, reply *int) error {
	after := time.Duration(arg.ExpireAfterSeconds) * time.Second
	return s.short.Refresh("", arg.Resource, after)
}
