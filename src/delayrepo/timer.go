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
	storage   TimerStorage
	pushArg   chan pushArg
	deleteArg chan deleteArg
	timeout   time.Duration
	tomb      tomb.Tomb
}

func NewTimer(storage TimerStorage, timeout time.Duration) (*Timer, error) {
	return &Timer{
		storage:   storage,
		pushArg:   make(chan pushArg),
		deleteArg: make(chan deleteArg),
		timeout:   timeout,
	}, nil
}

func (t *Timer) Serve(handler Handler) {
	defer t.tomb.Done()

	for {
		next, err := t.NextWakeup()
		if err != nil {
			handler.OnError(fmt.Errorf("next wake up failed: %s", err))
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
				go handler.Do(key, data)
			}
		case p := <-t.pushArg:
			err = t.push(p.updateType, p.ontime, p.key, p.data)
			p.err <- err
		case p := <-t.deleteArg:
			err = t.delete(p.key)
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
	err := <-arg.err
	close(arg.err)
	return err
}

type deleteArg struct {
	key string
	err chan error
}

func (t *Timer) Delete(key string) error {
	arg := deleteArg{
		key: key,
		err: make(chan error),
	}
	t.deleteArg <- arg
	err := <-arg.err
	close(arg.err)
	return err
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

func (t *Timer) delete(key string) error {
	_, err := t.storage.Load(key)
	return err
}

func (t *Timer) NextWakeup() (time.Duration, error) {
	key, err := t.storage.Next()
	if err != nil {
		return t.timeout, err
	}
	if key == "" {
		return t.timeout, nil
	}
	ontime, err := t.storage.Ontime(key)
	if err != nil {
		return t.timeout, err
	}
	next := time.Unix(ontime, 0).Sub(time.Now())
	return next, nil
}
