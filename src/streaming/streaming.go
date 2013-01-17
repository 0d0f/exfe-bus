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
	channels map[string]chan string
	locker   sync.Locker
	timeout  time.Duration
}

func New() *Streaming {
	return &Streaming{
		channels: make(map[string]chan string),
		locker:   new(sync.Mutex),
		timeout:  time.Second,
	}
}

func (s *Streaming) Connect(id string, conn net.Conn, w BufWriter) error {
	c, err := s.connecting(id)
	if err != nil {
		return err
	}
	defer s.shutdown(id)

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

func (s *Streaming) Feed(id, input string) (err error) {
	c, ok := s.channel(id)

	if !ok {
		return fmt.Errorf("%s not connected", id)
	}

	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	c <- input
	return
}

func (s *Streaming) connecting(id string) (chan string, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	if _, ok := s.channels[id]; ok {
		return nil, fmt.Errorf("has connected")
	}
	s.channels[id] = make(chan string)
	return s.channels[id], nil
}

func (s *Streaming) shutdown(id string) error {
	s.locker.Lock()
	defer s.locker.Unlock()

	if c, ok := s.channels[id]; !ok {
		return fmt.Errorf("no connection")
	} else {
		close(c)
	}
	delete(s.channels, id)
	return nil
}

func (s *Streaming) channel(id string) (chan string, bool) {
	s.locker.Lock()
	defer s.locker.Unlock()

	c, ok := s.channels[id]
	return c, ok
}
