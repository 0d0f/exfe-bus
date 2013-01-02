package main

import (
	"gobus"
	"model"
	"shorttoken"
	"time"
)

type ShortToken struct {
	short *shorttoken.ShortToken
}

func NewShortToken(repo shorttoken.Repo) *ShortToken {
	return &ShortToken{
		short: shorttoken.New(repo, 4),
	}
}

type PostArg struct {
	Resource          string `json:"resource"`
	Data              string `json:"data"`
	ExpireAfterSecond int    `json:"expire_after_second"`
}

func (s *ShortToken) POST(meta *gobus.HTTPMeta, arg PostArg, reply *model.Token) error {
	after := time.Duration(arg.ExpireAfterSecond) * time.Second
	var err error
	*reply, err = s.short.Create(arg.Resource, arg.Data, after)
	return err
}

func (s *ShortToken) GET(meta *gobus.HTTPMeta, arg string, reply *model.Token) error {
	params := meta.Request.URL.Query()
	key := params.Get("key")
	resource := params.Get("resource")
	var err error
	*reply, err = s.short.Get(key, resource)
	return err
}

type VerifyReply struct {
	Token   model.Token `json:"token"`
	Matched bool        `json:"matched"`
}

type UpdateArg struct {
	Data              *string `json:"data"`
	ExpireAfterSecond *int    `json:"expire_after_second"`
}

func (s *ShortToken) Update(meta *gobus.HTTPMeta, arg UpdateArg, reply *int) error {
	key := meta.Vars["key"]
	if arg.Data != nil {
		err := s.short.UpdateData(key, *arg.Data)
		if err != nil {
			return err
		}
	}
	if arg.ExpireAfterSecond != nil {
		after := time.Duration(*arg.ExpireAfterSecond) * time.Second
		err := s.short.Refresh(key, "", after)
		if err != nil {
			return err
		}
	}
	return nil
}
