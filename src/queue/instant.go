package main

import (
	"fmt"
	"gobus"
)

type Instant struct {
	services map[string]*gobus.Client
}

func NewInstant(services map[string]*gobus.Client) *Instant {
	return &Instant{
		services: services,
	}
}

func (i *Instant) Push(meta gobus.HTTPMeta, arg PushArg, count *int) error {
	client, ok := i.services[arg.Service]
	if !ok {
		return fmt.Errorf("can't find service %s", arg.Service)
	}
	datas := []interface{}{arg.Data}
	var r int
	err := client.Do(arg.Method, datas, &r)
	if err != nil {
		return err
	}
	*count = 1
	return nil
}
