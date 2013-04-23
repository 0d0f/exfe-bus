package splitter

import (
	"broker"
	"encoding/json"
	"fmt"
	"github.com/googollee/go-rest"
	"gobus"
	"model"
	"net/http"
)

type Splitter struct {
	rest.Service `prefix:"/v3/splitter"`

	Split rest.Processor `path:"" method:"POST"`

	dispatcher *gobus.Dispatcher
	config     *model.Config
}

func NewSplitter(config *model.Config, dispatcher *gobus.Dispatcher) *Splitter {
	return &Splitter{
		dispatcher: dispatcher,
		config:     config,
	}
}

func (s Splitter) HandleSplit(pack BigPack) {
	for _, to := range pack.Recipients {
		mergeKey := fmt.Sprintf("%s_%d", pack.MergeKey, to.IdentityID)
		pack.Data["to"] = to

		url := fmt.Sprintf("http://%s:%d/v3/queue/%s/%s/%s?ontime=%d&update=%s", s.config.ExfeQueue.Addr, s.config.ExfeQueue.Port, mergeKey, pack.Method, pack.Service, pack.Ontime, pack.Type)
		b, err := json.Marshal(pack.Data)
		if err != nil {
			s.Error(http.StatusBadRequest, err)
			return
		}

		go func() {
			resp, err := broker.Http("POST", url, "plain/text", b)
			if err != nil {
				s.config.Log.Err("|splitter|%s|%s|%s|", url, err, string(b))
			} else {
				resp.Body.Close()
			}
		}()
	}
}
