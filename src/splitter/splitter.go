package splitter

import (
	"broker"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-rest"
	"logger"
	"model"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Splitter struct {
	rest.Service `prefix:"/v3/splitter"`

	Split  rest.Processor `path:"" method:"POST"`
	Delete rest.Processor `path:"" method:"DELETE"`

	queueSite      string
	config         *model.Config
	speedDurations []string
}

func NewSplitter(config *model.Config) *Splitter {
	site := config.ExfeQueue.Addr
	if site == "0.0.0.0" || site == "" {
		site = "127.0.0.1"
	}
	var speedDurations []string
	if len(config.Splitter.SpeedOn) > 0 {
		for k := range config.Splitter.SpeedOn {
			speedDurations = append(speedDurations, k)
			sort.Strings(speedDurations)
		}
	}
	return &Splitter{
		queueSite:      site,
		config:         config,
		speedDurations: speedDurations,
	}
}

func (s Splitter) HandleSplit(pack BigPack) {
	b, err := base64.URLEncoding.DecodeString(pack.Service)
	if err != nil {
		s.Error(http.StatusBadRequest, s.GetError(4, fmt.Sprintf("service(%s) invalid: %s", pack.Service, err)))
		return
	}

	fl := logger.FUNC(pack.Method, string(b), pack.MergeKey, pack.Update, pack.Ontime, pack.Recipients)
	defer fl.Quit()

	pack.Ontime = s.speedon(pack.Ontime)

	for _, to := range pack.Recipients {
		mergeKey := fmt.Sprintf("%s_i%d", pack.MergeKey, to.IdentityID)
		pack.Data["to"] = to

		url := fmt.Sprintf("http://%s:%d/v3/queue/%s/%s/%s?ontime=%d&update=%s", s.queueSite, s.config.ExfeQueue.Port, mergeKey, pack.Method, pack.Service, pack.Ontime, pack.Update)
		b, err := json.Marshal(pack.Data)
		if err != nil {
			s.Error(http.StatusBadRequest, err)
			return
		}

		go func() {
			resp, err := broker.Http("POST", url, "plain/text", b)
			if err != nil {
				logger.ERROR("post %s error: %s, with %s", url, err, string(b))
			} else {
				resp.Body.Close()
			}
		}()
	}
}

func (s Splitter) HandleDelete(pack BigPack) {
	b, err := base64.URLEncoding.DecodeString(pack.Service)
	if err != nil {
		s.Error(http.StatusBadRequest, s.GetError(4, fmt.Sprintf("service(%s) invalid: %s", pack.Service, err)))
		return
	}

	fl := logger.FUNC(pack.Method, string(b), pack.MergeKey, pack.Update, pack.Ontime, pack.Recipients)
	defer fl.Quit()

	pack.Ontime = s.speedon(pack.Ontime)

	for _, to := range pack.Recipients {
		mergeKey := fmt.Sprintf("%s_i%d", pack.MergeKey, to.IdentityID)

		url := fmt.Sprintf("http://%s:%d/v3/queue/%s/%s/%s", s.config.ExfeQueue.Addr, s.config.ExfeQueue.Port, mergeKey, pack.Method, pack.Service)

		go func() {
			resp, err := broker.Http("DELETE", url, "plain/text", nil)
			if err != nil {
				logger.ERROR("delete %s error: %s", url, err)
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
	now := time.Now().Unix()
	d := ontime - now
	for i := len(s.speedDurations) - 1; i >= 0; i-- {
		k, err := strconv.ParseInt(s.speedDurations[i], 10, 64)
		if err != nil {
			continue
		}
		if d > k {
			return now + s.config.Splitter.SpeedOn[s.speedDurations[i]]
		}
	}
	return ontime
}
