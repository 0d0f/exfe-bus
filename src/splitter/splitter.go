package splitter

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-logger"
	"github.com/googollee/go-rest"
	"model"
	"net/http"
	"sort"
)

type Splitter struct {
	rest.Service `prefix:"/v3/splitter"`

	Split  rest.Processor `path:"" method:"POST"`
	Delete rest.Processor `path:"" method:"DELETE"`

	queueSite      string
	config         *model.Config
	log            *logger.SubLogger
	speedDurations Int64Slice
}

func NewSplitter(config *model.Config) *Splitter {
	site := config.ExfeQueue.Addr
	if site == "0.0.0.0" || site == "" {
		site = "127.0.0.1"
	}
	var speedDurations Int64Slice
	if len(config.Splitter.SpeedOn) > 0 {
		for k := range config.Splitter.SpeedOn {
			speedDurations = append(speedDurations, k)
			sort.Sort(speedDurations)
		}
	}
	return &Splitter{
		queueSite:      site,
		config:         config,
		log:            config.Log.Sub("splitter"),
		speedDurations: speedDurations,
	}
}

func (s Splitter) HandleSplit(pack BigPack) {
	log := s.log.SubCode()
	log.Debug("|post|%s|%s|%s|%s|%d|%v", pack.Method, pack.Service, pack.MergeKey, pack.Type, pack.Ontime, pack.Recipients)
	defer log.Debug("posted")

	for _, to := range pack.Recipients {
		mergeKey := fmt.Sprintf("%s_i%d", pack.MergeKey, to.IdentityID)
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
	log.Debug("|delete|%s|%s|%s|%s|%d|%v", pack.Method, pack.Service, pack.MergeKey, pack.Type, pack.Ontime, pack.Recipients)
	defer log.Debug("deleted")

	for _, to := range pack.Recipients {
		mergeKey := fmt.Sprintf("%s_i%d", pack.MergeKey, to.IdentityID)

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

func (s Splitter) speedon(ontime int64) int64 {
	if len(s.speedDurations) == 0 {
		return ontime
	}
	fmt.Println(s.speedDurations)
	for i := len(s.speedDurations) - 1; i >= 0; i-- {
		if ontime > s.speedDurations[i] {
			return s.config.Splitter.SpeedOn[s.speedDurations[i]]
		}
	}
	return ontime
}

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
