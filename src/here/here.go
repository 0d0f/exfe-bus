package here

import (
	"launchpad.net/tomb"
	"log/syslog"
	"time"
)

type findArg struct {
	token string
	ret   chan bool
}

type Here struct {
	cluster *Cluster
	tomb    tomb.Tomb
	timeout time.Duration
	update  chan Group
	add     chan *Data
	find    chan findArg
	log     *syslog.Writer
}

func New(threshold, signThreshold float64, timeout time.Duration) *Here {
	l, _ := syslog.New(syslog.LOG_DEBUG, "exfe_service")
	return &Here{
		cluster: NewCluster(threshold, signThreshold, timeout),
		timeout: timeout,
		update:  make(chan Group),
		add:     make(chan *Data),
		find:    make(chan findArg),
		log:     l,
	}
}

func (h *Here) Serve() {
	defer h.tomb.Done()

	for {
		select {
		case <-h.tomb.Dying():
			return
		case data := <-h.add:
			h.log.Debug("add")
			group := h.cluster.Add(data)
			if group != nil {
				h.log.Debug("updating")
				h.update <- *group
				h.log.Debug("updated")
			}
		case arg := <-h.find:
			h.log.Debug("find")
			_, ok := h.cluster.TokenGroup[arg.token]
			arg.ret <- ok
		case <-time.After(h.timeout):
		}
		h.log.Debug("clear")
		groups := h.cluster.Clear()
		for _, group := range groups {
			h.log.Debug("updating")
			h.update <- group
			h.log.Debug("updated")
		}
	}
}

func (h *Here) UpdateChannel() chan Group {
	return h.update
}

func (h *Here) Add(data *Data) error {
	err := data.Init()
	data.UpdatedAt = time.Now()
	if err != nil {
		return err
	}
	h.add <- data

	return nil
}

func (h *Here) Exist(token string) bool {
	arg := findArg{
		token: token,
		ret:   make(chan bool),
	}
	h.find <- arg
	return <-arg.ret
}
