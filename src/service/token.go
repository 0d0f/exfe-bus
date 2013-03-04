package main

import (
	"broker"
	"github.com/googollee/go-rest"
	"model"
	"net/http"
	"time"
	"token"
)

type Token struct {
	rest.Service `root:"/v3/token"`

	Create         rest.Processor `method:"POST" path:"/tokens/(short|long)"`
	KeyGet         rest.Processor `method:"GET" path:"/tokens/key/([a-zA-Z0-9]+)"`
	ResourceGet    rest.Processor `method:"GET" path:"/tokens/resource"`
	KeyUpdate      rest.Processor `method:"POST" path:"/tokens/key/([a-zA-Z0-9]+)"`
	ResourceUpdate rest.Processor `method:"POST" path:"/tokens/resource"`

	manager *token.Manager
}

func NewToken(config *model.Config, db *broker.DBMultiplexer) (*Token, error) {
	repo, err := NewTokenRepo(config, db)
	if err != nil {
		return nil, err
	}
	token := &Token{
		manager: token.New(repo),
	}
	return token, nil
}

type CreateArg struct {
	Data               string `json:"data"`
	Resource           string `json:"resource"`
	ExpireAfterSeconds int    `json:"expire_after_seconds"`
}

// 根据resource，data和expire after seconds创建一个token
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/tokens/long" -d '{"data":"abc","resource":"123","expire_after_seconds":300}'
//
// 返回：
//
//     {"key":"0303","data":"abc"}
func (s Token) Create_(genType string, arg CreateArg) (ret model.Token) {
	after := time.Duration(arg.ExpireAfterSeconds) * time.Second
	ret, err := s.manager.Create(genType, arg.Resource, arg.Data, after)
	if err != nil {
		s.Error(http.StatusNotFound, err)
		return
	}
	return ret
}

// 根据key获得一个token，如果token不存在，返回错误
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/token/key/0303"
//
// 返回：
//
//     [{"key":"0303","data":"abc"}]
func (s Token) KeyGet_(key string) []model.Token {
	ret, err := s.manager.Get(key, "")
	if err != nil {
		s.Error(http.StatusNotFound, err)
		return nil
	}
	return ret
}

// 根据resource获得一个token，如果token不存在，返回错误
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/token/resource" -d '"abc"'
//
// 返回：
//
//     [{"key":"0303","data":"abc"}]
func (s Token) ResourceGet_(resource string) []model.Token {
	ret, err := s.manager.Get("", resource)
	if err != nil {
		s.Error(http.StatusNotFound, err)
		return nil
	}
	return ret
}

type UpdateArg struct {
	Data               *string `json:"data"`
	ExpireAfterSeconds *int    `json:"expire_after_seconds"`
	Resource           string  `json:"resource"`
}

// 更新key对应的token的data信息或者expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/token/key/0303" -d '{"data":"xyz","expire_after_seconds":13}'
func (s Token) KeyUpdate_(key string, arg UpdateArg) {
	if arg.Data != nil {
		err := s.manager.UpdateData(key, *arg.Data)
		if err != nil {
			s.Error(http.StatusBadRequest, err)
			return
		}
	}
	if arg.ExpireAfterSeconds != nil {
		after := time.Duration(*arg.ExpireAfterSeconds) * time.Second
		err := s.manager.Refresh(key, "", after)
		if err != nil {
			s.Error(http.StatusBadRequest, err)
			return
		}
	}
}

// 更新resource对应的token的data信息或者expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/token/resource" -d '{"resource":"abc", "data":"xyz","expire_after_seconds":13}'
func (s Token) ResourceUpdate_(arg UpdateArg) {
	if arg.ExpireAfterSeconds != nil {
		after := time.Duration(*arg.ExpireAfterSeconds) * time.Second
		err := s.manager.Refresh("", arg.Resource, after)
		if err != nil {
			s.Error(http.StatusBadRequest, err)
			return
		}
	}
}
