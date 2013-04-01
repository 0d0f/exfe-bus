package delayrepo

import (
	"errors"
	"fmt"
	"github.com/googollee/go-logger"
	"launchpad.net/tomb"
	"time"
)

type Repo interface {
	Push(key string, data []byte) error
	Pop() (key string, datas [][]byte, err error)
	NextWakeup() (time.Duration, error)
}

var EmptyError = errors.New("Empty.")
var ChangedError = errors.New("Repository changed while poping.")

type Callback func(key string, datas [][]byte)

func ServRepository(log *logger.SubLogger, repo Repo, f Callback) *tomb.Tomb {
	var tomb tomb.Tomb

	go func() {
		defer tomb.Done()

		for {
			next, err := repo.NextWakeup()
			if err != nil {
				log.Crit("next wake up failed: %s", err)
			}
			select {
			case <-tomb.Dying():
				log.Info("quit")
				return
			case <-time.After(next):
				key, datas, err := repo.Pop()
				if err != nil {
					log.Crit("pop failed: %s", err)
					continue
				}
				if len(datas) > 0 {
					f(key, datas)
				}
			}
		}
	}()

	return &tomb
}

type DelayStrategy interface {
	Push(ontime int64, key string, data []byte) error
	Pop() (key string, datas [][]byte, err error)
	NextWakeup() (time.Duration, error)
}

type Handler interface {
	Do(key string, data [][]byte)
	OnError(err error)
}

type Repository struct {
	tomb     tomb.Tomb
	timeout  time.Duration
	strategy DelayStrategy
	handler  Handler
	push     chan pushArg
}

func New(strategy DelayStrategy, handler Handler, timeout time.Duration) *Repository {
	ret := &Repository{
		timeout:  timeout,
		strategy: strategy,
		handler:  handler,
		push:     make(chan pushArg),
	}
	return ret
}

func (r *Repository) Serve() {
	defer r.tomb.Done()

	for {
		next, err := r.strategy.NextWakeup()
		if err != nil {
			r.handler.OnError(fmt.Errorf("next wake up failed: %s", err))
		}
		if next < 0 {
			next = r.timeout
		}
		select {
		case <-r.tomb.Dying():
			return
		case <-time.After(next):
			key, data, err := r.strategy.Pop()
			if err != nil {
				r.handler.OnError(fmt.Errorf("pop failed: %s", err))
				continue
			}
			if len(data) > 0 {
				r.handler.Do(key, data)
			}
		case p := <-r.push:
			err = r.strategy.Push(p.ontime, p.key, p.data)
			p.err <- err
		}
	}
}

func (r *Repository) Push(ontime int64, key string, data []byte) error {
	push := pushArg{
		ontime: ontime,
		key:    key,
		data:   data,
		err:    make(chan error),
	}
	r.push <- push
	return <-push.err
}

func (r *Repository) Quit() {
	r.tomb.Kill(nil)
	r.tomb.Wait()
}

type pushArg struct {
	ontime int64
	key    string
	data   []byte
	err    chan error
}
