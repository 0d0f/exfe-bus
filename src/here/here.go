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
	group := h.GetGroup(h.cluster.UserGroup[user.Id])
	if group == nil {
		h.update <- user.Id
	} else {
		for _, u := range group.Users {
			h.update <- u.Id
		}
	}
}

func (h *Here) GetGroup(id string) *Group {
	return h.cluster.Groups[id]
}

func (h *Here) UserInGroup(userId string) *Group {
	id, ok := h.cluster.UserGroup[userId]
	if !ok {
		return nil
	}
	return h.cluster.Groups[id]
}
