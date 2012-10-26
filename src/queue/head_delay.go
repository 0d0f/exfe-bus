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
	services map[string]*gobus.Client
	repo     *delayrepo.Head
	config   *model.Config
}

func NewHead(services map[string]*gobus.Client, delayInMinute int, config *model.Config) (*Head, *tomb.Tomb) {
	name := fmt.Sprintf("delayrepo:head_%dm", delayInMinute)
	delay := delayInMinute * 60
	redis := broker.NewRedisImp()
	repo := delayrepo.NewHead(name, delay, redis)
	log := config.Log.SubPrefix(name)
	tomb := delayrepo.ServRepository(log, repo, getCallback(log, services))

	return &Head{
		services: services,
		repo:     repo,
		config:   config,
	}, tomb
}

func (i *Head) Push(meta *gobus.HTTPMeta, arg PushArg, count *int) error {
	data, err := json.Marshal(arg.Data)
	if err != nil {
		return fmt.Errorf("can't marshal input data: %s", err)
	}
	err = i.repo.Push(fmt.Sprintf("%s,%s,%s", arg.Service, arg.Method, arg.Key), data)
	if err != nil {
		return fmt.Errorf("push to repo failed: %s", err)
	}
	*count = 1
	return nil
}

type Head10m struct {
	*Head
}

func NewHead10m(services map[string]*gobus.Client, config *model.Config) (*Head10m, *tomb.Tomb) {
	repo, tomb := NewHead(services, 10, config)
	return &Head10m{repo}, tomb
}
