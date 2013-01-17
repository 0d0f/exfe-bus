package streaming

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type BufWriter interface {
	io.Writer
	Flush() error
}

type Streaming struct {
	channels map[string][]chan string
	locker   sync.Locker
	timeout  time.Duration
}

func New(timeout time.Duration) *Streaming {
	return &Streaming{
		channels: make(map[string][]chan string),
		locker:   new(sync.Mutex),
		timeout:  timeout,
	}
}

func (s *Streaming) Connect(id string, conn net.Conn, w BufWriter) error {
	c, err := s.connecting(id)
	if err != nil {
		return err
	}
	defer s.shutdown(id, c)

	p := make([]byte, 512)
	for {
		select {
		case input := <-c:
			_, err := w.Write([]byte(input))
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("\n"))
			if err != nil {
				return err
			}
			err = w.Flush()
			if err != nil {
				return err
			}
		case <-time.After(s.timeout):
			conn.SetReadDeadline(time.Now())
			_, err := conn.Read(p)
			if err != nil {
				if netErr, ok := err.(net.Error); !(ok && netErr.Timeout()) {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Streaming) Feed(id, content string) (err error) {
	conns := s.getConns(id)
	if conns == nil {
		return fmt.Errorf("%s not connected", id)
	}

	for _, c := range conns {
		s.send(c, content)
	}
	return
}

func (s *Streaming) connecting(id string) (chan string, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	ret := make(chan string)
	if _, ok := s.channels[id]; ok {
		s.channels[id] = append(s.channels[id], ret)
	} else {
		s.channels[id] = []chan string{ret}
	}
	return ret, nil
}

func (s *Streaming) shutdown(id string, c chan string) error {
	s.locker.Lock()
	defer s.locker.Unlock()

	conns, ok := s.channels[id]
	if !ok {
		return fmt.Errorf("no connection")
	}

	var i int = -1
	for i = range conns {
		if conns[i] == c {
			break
		}
	}
	close(c)

	if i == len(conns) {
		return fmt.Errorf("no connection")
	}
	conns = append(conns[:i], conns[i+1:]...)
	if len(conns) == 0 {
		delete(s.channels, id)
	} else {
		s.channels[id] = conns
	}
	return nil
}

func (s *Streaming) getConns(id string) []chan string {
	s.locker.Lock()
	defer s.locker.Unlock()

	conns, ok := s.channels[id]
	if !ok {
		return nil
	}
	ret := make([]chan string, len(conns))
	for i := range conns {
		ret[i] = conns[i]
	}
	return ret
}

func (s *Streaming) send(c chan string, content string) {
	defer func() {
		recover()
	}()

	c <- content
}
