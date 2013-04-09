package here

import (
	"launchpad.net/tomb"
	"log/syslog"
	"sync"
	"time"
)

type Here struct {
	cluster *Cluster
	tomb    tomb.Tomb
	timeout time.Duration
	update  chan string
	locker  sync.Mutex
	log     *syslog.Writer
}

func New(threshold, signThreshold float64, timeout time.Duration) *Here {
	l, _ := syslog.New(syslog.LOG_DEBUG, "exfe_service")

	ret := &Here{
		cluster: NewCluster(threshold, signThreshold, timeout),
		timeout: timeout,
		update:  make(chan string),
		log:     l,
	}
	go func() {
		defer ret.tomb.Done()

		for {
			select {
			case <-ret.tomb.Dying():
				return
			case <-time.After(timeout):
				l.Debug("loop clear lock")
				ret.locker.Lock()
				ids := ret.cluster.Clear()
				ret.locker.Unlock()
				l.Debug("loop clear unlock")
				for _, id := range ids {
					ret.update <- id
				}
			}
		}
	}()
	return ret
}

func (h *Here) UpdateChannel() chan string {
	return h.update
}

func (h *Here) Add(data *Data) error {
	h.log.Debug("add lock")
	h.locker.Lock()
	err := h.cluster.Add(data)
	h.locker.Unlock()
	h.log.Debug("add lock")

	if err != nil {
		return err
	}
	group := h.TokenInGroup(data.Token)
	if group == nil {
		h.update <- data.Token
	} else {
		for _, u := range group.Data {
			h.update <- u.Token
		}
	}
	return nil
}

func (h *Here) TokenInGroup(token string) *Group {
	h.log.Debug("tokeningroup lock")
	h.locker.Lock()
	defer h.locker.Unlock()
	defer h.log.Debug("tokeningroup lock")
	id, ok := h.cluster.TokenGroup[token]
	if !ok {
		return nil
	}
	return h.cluster.Groups[id]
}
