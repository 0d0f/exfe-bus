package delayrepo

import (
	"errors"
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
