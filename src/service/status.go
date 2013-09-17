package main

import (
	"github.com/googollee/go-rest"
	"net/http"
)

type Status struct {
	rest.Service `prefix:"/_status"`

	live rest.SimpleNode `method:"POST" route:"/live"`
	show rest.SimpleNode `method:"GET" route:"/:key"`

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

func (s *Status) Live(ctx rest.Context) {
	ctx.Render("OK")
}

func (s *Status) Show(ctx rest.Context) {
	var key string
	ctx.Bind("key", &key)
	if err := ctx.BindError(); err != nil {
		ctx.Return(http.StatusBadRequest, "%s", err)
		return
	}
	ret, ok := s.data[key]
	if !ok {
		ctx.Return(http.StatusBadRequest, "can't find key(%s)", key)
		return
	}
	ctx.Render(ret)
}
