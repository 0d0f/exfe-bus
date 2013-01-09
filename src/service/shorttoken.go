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

func (s *ShortToken) SetRoute(route gobus.RouteCreater) error {
	json := new(gobus.JSON)
	route().Methods("POST").Path("/shorttoken").HandlerMethod(json, s, "Create")
	route().Methods("GET").Path("/shorttoken").HandlerMethod(json, s, "Get")
	route().Methods("POST", "PUT").Path("/shorttoken/{key}").HandlerMethod(json, s, "Update")
	route().Methods("POST").Path("/shorttoken/resource").HandlerMethod(json, s, "UpdateResource")

	// old
	route().Methods("POST").Path("/shorttoken").Queries("method", "GET").HandlerMethod(json, s, "Get")

	return nil
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
//     > curl "http://127.0.0.1:23333/shorttoken" -d '{"data":"abc","resource":"123","expire_after_seconds":300}'
//
// 返回：
//
//     {"key":"0303","data":"abc"}
func (s *ShortToken) Create(params map[string]string, arg CreateArg) (model.Token, error) {
	after := time.Duration(arg.ExpireAfterSeconds) * time.Second
	ret, err := s.short.Create(arg.Resource, arg.Data, after)
	return ret, err
}

// 根据key或者resource获得一个token，如果token不存在，返回错误
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/shorttoken?key=0303&resource=123"
//
// 返回：
//
//     [{"key":"0303","data":"abc"}]
func (s *ShortToken) Get(params map[string]string) ([]model.Token, error) {
	key := params["key"]
	resource := params["resource"]
	ret, err := s.short.Get(key, resource)
	return ret, err
}

type UpdateArg struct {
	Data               *string `json:"data"`
	ExpireAfterSeconds *int    `json:"expire_after_seconds"`
}

// 更新key对应的token的data信息或者expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/shorttoken/0303" -d '{"data":"xyz","expire_after_seconds":13}'
//
// 返回：
//
//     0
func (s *ShortToken) Update(params map[string]string, arg UpdateArg) (int, error) {
	key := params["key"]
	if arg.Data != nil {
		err := s.short.UpdateData(key, *arg.Data)
		if err != nil {
			return 0, err
		}
	}
	if arg.ExpireAfterSeconds != nil {
		after := time.Duration(*arg.ExpireAfterSeconds) * time.Second
		err := s.short.Refresh(key, "", after)
		if err != nil {
			return 0, err
		}
	}
	return 0, nil
}

// 更新resource对应的token的data信息或者expire after seconds
//
// 例子：
//
//     > curl "http://127.0.0.1:23333/shorttoken/resource?resouce=123" -d '13'
//
// 返回：
//
//     0
func (s *ShortToken) UpdateResource(params map[string]string, expire int) (int, error) {
	resource := params["resource"]
	after := time.Duration(expire) * time.Second
	err := s.short.Refresh("", resource, after)
	return 0, err
}
