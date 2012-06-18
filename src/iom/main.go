package main

import (
	"fmt"
	"regexp"
	"net/http"
	"exfe/service"
	"log"
	"flag"
	"config"
	"os"
)

const HashUserPattern = "^/iom/(.*?)$"
const HashTagPattern = "^/iom/(.*?)/(.*?)$"

type HashHTTP struct {
	handler *HashHandler
	log *log.Logger
	userReg *regexp.Regexp
	tagReg *regexp.Regexp
}

func (h *HashHTTP) Get(userid, hash string) (string, error) {
	return h.handler.Get(userid, hash)
}

func (h *HashHTTP) Post(userid, data string) (string, error) {
	hash, err := h.handler.FindByData(userid, data)
	if err != nil {
		hash, err = h.handler.Create(userid, data)
	}
	return hash, err
}

func (h *HashHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var e error
	switch r.Method {
	case "GET":
		args := h.tagReg.FindStringSubmatch(r.URL.Path)
		if len(args) == 3 {
			url, err := h.Get(args[1], args[2])
			if err == nil {
				fmt.Fprintf(w, "%s", url)
				return
			}
			e = err
		} else {
			e = fmt.Errorf("can't parse url")
		}
	case "POST":
		r.ParseForm()
		args := h.userReg.FindStringSubmatch(r.URL.Path)
		if len(args) == 2 {
			var err error
			var hash string
			if data, ok := r.Form["data"]; ok {
				hash, err = h.Post(args[1], data[0])
			} else {
				err = fmt.Errorf("can't find data params in post data")
			}
			if err == nil {
				fmt.Fprintf(w, "%s", hash)
				return
			}
			e = err
		} else {
			e = fmt.Errorf("can't parse url")
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", e)
}

func main() {
	log := log.New(os.Stderr, "exfe.hash", log.LstdFlags)
	log.Print("Service start")

	var c exfe_service.Config

	var pidfile string
	var configFile string

	flag.StringVar(&pidfile, "pid", "", "Specify the pid file")
	flag.StringVar(&configFile, "config", "exfe.json", "Specify the configuration file")
	flag.Parse()

	config.LoadFile(configFile, &c)

	flag.Parse()
	if pidfile != "" {
		pid, err := os.Create(pidfile)
		if err != nil {
			log.Fatalf("Can't create pid(%s): %s", pidfile, err)
		}
		pid.WriteString(fmt.Sprintf("%d", os.Getpid()))
	}

	handler := &HashHTTP{
		handler: NewHashHandler(c.Redis.Netaddr, c.Redis.Db, c.Redis.Password),
		log: log,
		userReg: regexp.MustCompile(HashUserPattern),
		tagReg: regexp.MustCompile(HashTagPattern),
	}
	http.Handle("/puh/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Error: %s", err)
	}
}