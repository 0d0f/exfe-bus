package main

import (
	"github.com/googollee/go-rest"
	"github.com/gorilla/mux"
	"here"
	"model"
	"net/http"
	"time"
)

type HereService struct {
	rest.Service `prefix:"/v3/here"`

	Users rest.Processor `path:"/users" method:"POST"`

	here *here.Here
}

func (h HereService) Users_(user here.User) {
	h.here.Add(user)
}

type HereStreaming struct {
	ids map[string][]chan string
}

func (h *HereStreaming) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token invalid", http.StatusBadRequest)
		return
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "doesn't support streaming", http.StatusInternalServerError)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c := make(chan string)
	h.ids[token] = append(h.ids[token], c)
	defer func() {
		conn.Close()
		for i, ch := range h.ids[token] {
			if ch == c {
				h.ids[token] = append(h.ids[token][:i], h.ids[token][i+1:]...)
				return
			}
		}
	}()

	for {
		select {
		case <-time.After(time.Second):
			_, err = bufrw.Write([]byte("\n"))
			if err != nil {
				return
			}
			err = bufrw.Flush()
			if err != nil {
				return
			}
		case data := <-c:
			_, err = bufrw.Write([]byte(data + "\n"))
			if err != nil {
				return
			}
			err = bufrw.Flush()
			if err != nil {
				return
			}
		}
	}
}

func NewHere(config *model.Config) (http.Handler, error) {
	ret := mux.NewRouter()
	service := new(HereService)
	service.here = here.New(config.Here.Threshold, config.Here.SignThreshold, time.Duration(config.Here.TimeoutInSecond)*time.Second)
	handler, err := rest.New(service)
	if err != nil {
		return nil, err
	}
	streaming := &HereStreaming{
		ids: make(map[string][]chan string),
	}
	ret.PathPrefix(handler.Prefix()).Handler(handler)
	ret.Path("/here/streaming").Handler(streaming)
	return ret, nil
}
