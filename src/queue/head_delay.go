package main

import (
	"broker"
	"delayrepo"
	"encoding/json"
	"fmt"
	"gobus"
	"model"
)

type Head struct {
	services map[string]*gobus.Client
	repo     *delayrepo.Head
	config   *model.Config
}

func NewHead(services map[string]*gobus.Client, delayInMinute int, config *model.Config, quit chan int) *Head {
	name := fmt.Sprintf("delayrepo:head_%sm", delayInMinute)
	delay := delayInMinute * 60
	redis := broker.NewRedisImp()
	repo := delayrepo.NewHead(name, delay, redis)
	log := config.Log.SubPrefix(name)
	go delayrepo.ServRepository(log, repo, quit, getCallback(log, services))

	return &Head{
		services: services,
		repo:     repo,
		config:   config,
	}
}

func (i *Head) Push(meta gobus.HTTPMeta, arg PushArg, count *int) error {
	data, err := json.Marshal(arg.Data)
	if err != nil {
		return fmt.Errorf("can't marshal input data: %s", err)
	}
	err = i.repo.Push(arg.Key, data)
	if err != nil {
		return fmt.Errorf("push to repo failed: %s", err)
	}
	*count = 1
	return nil
}

type Head10m struct {
	*Head
}

func NewHead10m(services map[string]*gobus.Client, config *model.Config, quit chan int) *Head10m {
	return &Head10m{
		NewHead(services, 10, config, quit),
	}
}
