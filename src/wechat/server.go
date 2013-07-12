package main

import (
	"broker"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Server struct {
	wc      *WeChat
	kvSaver *broker.KVSaver
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	if strings.HasPrefix(to, "e") && strings.HasSuffix(to, "@exfe") {
		defer r.Body.Close()
		chatroomId, exist, err := s.kvSaver.Check([]string{to})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !exist {
			http.Error(w, fmt.Sprintf("can't find chatroom for %s", to), http.StatusBadRequest)
			return
		}
		to = chatroomId
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(to, string(b))
	err = s.wc.SendMessage(to, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func runServer(addr string, port uint, wc *WeChat, kvSaver *broker.KVSaver) {
	http.Handle("/send", &Server{wc, kvSaver})
	go http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)

}
