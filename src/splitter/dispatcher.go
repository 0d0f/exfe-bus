package splitter

import (
	"fmt"
	"github.com/googollee/go-rest"
	"gobus"
	"model"
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

func (d Splitter) HandleSplit(pack BigPack) {
	for p := range pack.Each() {
		url := fmt.Sprintf("bus://exfe/v3/queue/%s/%s/%s", p.MergeKey, p.Method, p.Service)
		var reply interface{}
		err := d.dispatcher.Do(url, "POST", p.Data, &reply)
		if err != nil {
			d.config.Log.Err("|dispatcher|%s|%s|%v|", url, err, p.Data)
		}
	}
}
