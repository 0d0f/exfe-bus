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
	*count = 0
	if len(arg.Tos) == 0 {
		var r int
		err := client.Do(arg.Method, arg.Data, &r)
		if err == nil {
			*count = 1
		}
	} else {
		for _, to := range arg.Tos {
			data, ok := arg.Data.(map[string]interface{})
			if ok {
				data["to"] = to
			}
			var r int
			err := client.Do(arg.Method, data, &r)
			if err == nil {
				*count++
			}
		}
	}
	return nil
}
