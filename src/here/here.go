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
		update:  make(chan string),
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

func (h *Here) Add(user User) {
	h.locker.Lock()
	h.cluster.AddUser(&user)
	h.locker.Unlock()
	h.update <- h.cluster.UserGroup[user.Id]
}

func (h *Here) GetGroup(id string) *Group {
	ret, ok := h.cluster.Groups[id]
	if !ok {
		ret = NewGroup()
	}
	return ret
}

func (h *Here) UserInGroupId(userId string) string {
	return h.cluster.UserGroup[userId]
}
