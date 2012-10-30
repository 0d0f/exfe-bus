package main

import (
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

func (i *Instant) Push(meta *gobus.HTTPMeta, arg PushArg, count *int) error {
	client, err := arg.FindService(i.services)
	if err != nil {
		return err
	}
	datas, _ := arg.Expand()
	*count = 0
	for _, data := range datas {
		var r int
		err := client.Do(arg.Method, data, &r)
		if err == nil {
			*count++
		}
	}
	return nil
}
