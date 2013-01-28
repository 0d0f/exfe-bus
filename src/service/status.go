package main

import (
	"fmt"
	"gobus"
)

type Status struct {
	data map[string]int
}

func NewStatus() *Status {
	return &Status{
		data: make(map[string]int),
	}
}

func (s *Status) Set(key string, data int) {
	s.data[key] = data
}

func (s *Status) Get(key string) int {
	return s.data[key]
}

func (s *Status) SetRoute(r gobus.RouteCreater) (err error) {
	json := new(gobus.JSON)
	err = r().Methods("GET", "POST").Path("/_status/live").HandlerMethod(json, s, "Live")
	err = r().Methods("GET").Path("/_status/{key}").HandlerMethod(json, s, "Show")
	return
}

func (s *Status) Live(params map[string]string) (string, error) {
	return "OK", nil
}

func (s *Status) Show(params map[string]string) (int, error) {
	key := params["key"]
	ret, ok := s.data[key]
	if !ok {
		return 0, fmt.Errorf("can't find key(%s)", key)
	}
	return ret, nil
}
