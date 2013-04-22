package delayrepo

import (
	"broker"
	"fmt"
	"launchpad.net/tomb"
	"time"
)

type UpdateType string

const (
	Always UpdateType = "always"
	Once              = "once"
)

type Handler interface {
	Do(key string, data [][]byte)
	OnError(err error)
}

type TimerStorage interface {
	Save(updateType broker.UpdateType, ontime int64, key string, data []byte) error
	Load(key string) ([][]byte, error)
	Ontime(key string) (int64, error)
	Next() (string, error)
}

type Timer struct {
	storage TimerStorage
	pushArg chan pushArg
	tomb    tomb.Tomb
}

func NewTimer(storage TimerStorage) (*Timer, error) {
	return &Timer{
		storage: storage,
		pushArg: make(chan pushArg),
	}, nil
}

func (t *Timer) Serve(handler Handler, timeout time.Duration) {
	defer t.tomb.Done()

	for {
		next, err := t.NextWakeup()
		if err != nil {
			handler.OnError(fmt.Errorf("next wake up failed: %s", err))
		}
		if next < 0 {
			next = timeout
		}
		select {
		case <-t.tomb.Dying():
			return
		case <-time.After(next):
			key, data, err := t.pop()
			if err != nil {
				handler.OnError(fmt.Errorf("pop failed: %s", err))
				continue
			}
			if len(data) > 0 {
				handler.Do(key, data)
			}
		case p := <-t.pushArg:
			err = t.push(p.updateType, p.ontime, p.key, p.data)
			p.err <- err
		}
	}
}

func (t *Timer) Quit() {
	t.tomb.Kill(nil)
	t.tomb.Wait()
}

type pushArg struct {
	updateType broker.UpdateType
	ontime     int64
	key        string
	data       []byte
	err        chan error
}

func (t *Timer) Push(updateType UpdateType, ontime int64, key string, data []byte) error {
	switch updateType {
	case Always:
	case Once:
	default:
		return fmt.Errorf("invalid update type: %s", updateType)
	}
	arg := pushArg{
		updateType: broker.UpdateType(updateType),
		ontime:     ontime,
		key:        key,
		data:       data,
		err:        make(chan error),
	}
	t.pushArg <- arg
	return <-arg.err
}

func (t *Timer) push(updateType broker.UpdateType, ontime int64, key string, data []byte) error {
	return t.storage.Save(updateType, ontime, key, data)
}

func (t *Timer) pop() (string, [][]byte, error) {
	key, err := t.storage.Next()
	if err != nil {
		return "", nil, err
	}
	data, err := t.storage.Load(key)
	if err != nil {
		return "", nil, err
	}
	return key, data, nil
}

func (t *Timer) NextWakeup() (time.Duration, error) {
	key, err := t.storage.Next()
	if err != nil {
		return -1, err
	}
	ontime, err := t.storage.Ontime(key)
	if err != nil {
		return -1, err
	}
	next := time.Unix(ontime, 0).Sub(time.Now())
	if next < 0 {
		next = 0
	}
	return next, nil
}
