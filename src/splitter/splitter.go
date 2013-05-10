package splitter

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/googollee/go-rest"
	"model"
	"net/http"
)

type Splitter struct {
	rest.Service `prefix:"/v3/splitter"`

	Split  rest.Processor `path:"" method:"POST"`
	Delete rest.Processor `path:"" method:"DELETE"`

	queueSite string
	config    *model.Config
	log       *logger.SubLogger
}

func NewSplitter(config *model.Config) *Splitter {
	site := config.ExfeQueue.Addr
	if site == "0.0.0.0" || site == "" {
		site = "127.0.0.1"
	}
	return &Splitter{
		queueSite: site,
		config:    config,
		log:       config.Log.Sub("splitter"),
	}
}

func (s Splitter) HandleSplit(pack BigPack) {
	log := s.log.SubCode()
	log.Debug("post to %s,%s,%s,%v", pack.Method, pack.Service, pack.MergeKey, pack.Recipients)
	defer log.Debug("posted")

	for _, to := range pack.Recipients {
		mergeKey := fmt.Sprintf("%s_%d", pack.MergeKey, to.IdentityID)
		pack.Data["to"] = to

		url := fmt.Sprintf("http://%s:%d/v3/queue/%s/%s/%s?ontime=%d&update=%s", s.queueSite, s.config.ExfeQueue.Port, mergeKey, pack.Method, pack.Service, pack.Ontime, pack.Type)
		b, err := json.Marshal(pack.Data)
		if err != nil {
			s.Error(http.StatusBadRequest, err)
			return
		}

		go func() {
			resp, err := broker.Http("POST", url, "plain/text", b)
			if err != nil {
				s.config.Log.Err("|splitter|POST|%s|%s|%s|", url, err, string(b))
			} else {
				resp.Body.Close()
			}
		}()
	}
}

func (s Splitter) HandleDelete(pack BigPack) {
	log := s.log.SubCode()
	log.Debug("delete to %s,%s,%s,%v", pack.Method, pack.Service, pack.MergeKey, pack.Recipients)
	defer log.Debug("deleted")

	for _, to := range pack.Recipients {
		mergeKey := fmt.Sprintf("%s_%d", pack.MergeKey, to.IdentityID)

		url := fmt.Sprintf("http://%s:%d/v3/queue/%s/%s/%s", s.config.ExfeQueue.Addr, s.config.ExfeQueue.Port, mergeKey, pack.Method, pack.Service)

		go func() {
			resp, err := broker.Http("DELETE", url, "plain/text", nil)
			if err != nil {
				s.config.Log.Err("|splitter|DELETE|%s|%s|", url, err)
			} else {
				resp.Body.Close()
			}
		}()
	}
}
