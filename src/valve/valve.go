package valve

import (
	"errors"
	"launchpad.net/tomb"
	"time"
)

var QueueFull = errors.New("queue full")

type Worker interface {
	Do() (interface{}, error)
}

type Valve struct {
	push   chan request
	period time.Duration
	tomb   tomb.Tomb
}

type response struct {
	ret interface{}
	err error
}

type request struct {
	worker Worker
	ret    chan response
}

func New(depth int, period time.Duration) *Valve {
	return &Valve{
		push:   make(chan request, depth),
		period: period,
	}
}

func (v *Valve) Serve() {
	defer v.tomb.Done()

	for {
		select {
		case request := <-v.push:
			begin := time.Now()

			ret, err := request.worker.Do()
			request.ret <- response{ret, err}

			end := time.Now()
			if d := end.Sub(begin); d < v.period {
				time.Sleep(v.period - d)
			}
		case <-v.tomb.Dying():
			return
		}
	}
}

func (v *Valve) Quit() {
	v.tomb.Kill(nil)
	v.tomb.Wait()
}

func (v *Valve) Do(worker Worker) (interface{}, error) {
	req := request{worker, make(chan response)}
	select {
	case v.push <- req:
		resp := <-req.ret
		return resp.ret, resp.err
	default:
	}
	return nil, QueueFull
}
