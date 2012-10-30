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
	name     string
	repo     *delayrepo.Head
	config   *model.Config
}

func NewHead(services map[string]*gobus.Client, delayInSecond int, config *model.Config) (*Head, *tomb.Tomb) {
	name := fmt.Sprintf("delayrepo:head_%ds", delayInSecond)
	delay := delayInSecond
	redis := broker.NewRedisImp()
	repo := delayrepo.NewHead(name, delay, redis)
	log := config.Log.SubPrefix(name)
	tomb := delayrepo.ServRepository(log, repo, getCallback(log, services))

	return &Head{
		services: services,
		name:     name,
		repo:     repo,
		config:   config,
	}, tomb
}

func (i *Head) Push(meta *gobus.HTTPMeta, arg PushArg, count *int) error {
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
		err = i.repo.Push(fmt.Sprintf("%s,%s,%s,%s", arg.Service, arg.Method, arg.MergeKey, to), d)
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
