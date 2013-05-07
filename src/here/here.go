package here

import (
	"launchpad.net/tomb"
	"time"
)

type findArg struct {
	token string
	ret   chan *Group
}

type Here struct {
	cluster *Cluster
	tomb    tomb.Tomb
	timeout time.Duration
	update  chan Group
	add     chan *Data
	find    chan findArg
}

func New(threshold, signThreshold float64, timeout time.Duration) *Here {
	return &Here{
		cluster: NewCluster(threshold, signThreshold, timeout),
		timeout: timeout,
		update:  make(chan Group),
		add:     make(chan *Data),
		find:    make(chan findArg),
	}
}

func (h *Here) Serve() {
	defer h.tomb.Done()

	for {
		select {
		case <-h.tomb.Dying():
			return
		case data := <-h.add:
			group := h.cluster.Add(data)
			if group != nil {
				h.update <- *group
			}
		case arg := <-h.find:
			key, ok := h.cluster.TokenGroup[arg.token]
			if !ok {
				arg.ret <- nil
				continue
			}
			group, ok := h.cluster.Groups[key]
			if !ok {
				arg.ret <- nil
				continue
			}
			arg.ret <- group
		case <-time.After(h.timeout):
		}
		groups := h.cluster.Clear()
		for _, group := range groups {
			h.update <- group
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

func (h *Here) Exist(token string) *Group {
	arg := findArg{
		token: token,
		ret:   make(chan *Group),
	}
	h.find <- arg
	ret := <-arg.ret
	close(arg.ret)
	return ret
}
