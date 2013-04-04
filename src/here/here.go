package here

import (
	"launchpad.net/tomb"
	"sync"
	"time"
)

type Here struct {
	cluster *Cluster
	tomb    tomb.Tomb
	timeout time.Duration
	update  chan string
	locker  sync.Mutex
}

func New(threshold, signThreshold float64, timeout time.Duration) *Here {
	ret := &Here{
		cluster: NewCluster(threshold, signThreshold, timeout),
		timeout: timeout,
		update:  make(chan string, 1),
	}
	go func() {
		defer ret.tomb.Done()

		for {
			select {
			case <-ret.tomb.Dying():
				return
			case <-time.After(timeout):
			}
			ret.locker.Lock()
			ids := ret.cluster.Clear()
			ret.locker.Unlock()
			for _, id := range ids {
				ret.update <- id
			}
		}
	}()
	return ret
}

func (h *Here) UpdateChannel() chan string {
	return h.update
}

func (h *Here) Add(data Data) error {
	h.locker.Lock()
	err := h.cluster.AddUser(&data)
	h.locker.Unlock()

	if err != nil {
		return err
	}
	group := h.UserInGroup(data.Token)
	if group == nil {
		h.update <- data.Token
	} else {
		for _, u := range group.Data {
			h.update <- u.Token
		}
	}
	return nil
}

func (h *Here) UserInGroup(userId string) *Group {
	h.locker.Lock()
	defer h.locker.Unlock()
	id, ok := h.cluster.UserGroup[userId]
	if !ok {
		return nil
	}
	return h.cluster.Groups[id]
}
