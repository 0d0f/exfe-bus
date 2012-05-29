package main

import (
	"fmt"
	"regexp"
	"net/http"
	"exfe/service"
	"log/syslog"
	"flag"
	"config"
	"os"
)

const HashUserPattern = "^/puh/(.*?)$"
const HashTagPattern = "^/puh/(.*?)/(.*?)$"

type HashHTTP struct {
	handler *HashHandler
	log *syslog.Writer
	userReg *regexp.Regexp
	tagReg *regexp.Regexp
}

func (h *HashHTTP) Get(userid, hash string) (string, error) {
	return h.handler.Get(userid, hash)
}

func (h *HashHTTP) Post(userid, url string) (string, error) {
	hash, err := h.handler.FindByUrl(userid, url)
	if err != nil {
		hash, err = h.handler.Create(userid, url)
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
			hash, err := h.Post(args[1], r.Form["url"][0])
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
	log, err := syslog.New(syslog.LOG_INFO, "exfe.hash")
	if err != nil {
		panic(err)
	}
	log.Info("Service start")

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
			log.Crit(fmt.Sprintf("Can't create pid(%s): %s", pidfile, err))
			return
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
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
