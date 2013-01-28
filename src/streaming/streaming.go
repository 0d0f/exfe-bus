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
	channels map[string]map[string]chan string
	locker   sync.Locker
	timeout  time.Duration
}

func New(timeout time.Duration) *Streaming {
	return &Streaming{
		channels: make(map[string]map[string]chan string),
		locker:   new(sync.Mutex),
		timeout:  timeout,
	}
}

func (s *Streaming) Connect(id, key string, conn net.Conn, w BufWriter) error {
	c, err := s.connecting(id, key)
	if err != nil {
		return err
	}
	defer s.shutdown(id, key)

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

func (s *Streaming) Feed(id, content string) error {
	conns := s.getConns(id)
	if conns == nil {
		return fmt.Errorf("id(%s) not connected", id)
	}

	for _, c := range conns {
		s.send(c, content)
	}
	return nil
}

func (s *Streaming) Send(id, key, content string) error {
	c := s.getConn(id, key)
	if c == nil {
		return fmt.Errorf("id(%s)/key(%s) not connected", id, key)
	}
	s.send(c, content)
	return nil
}

func (s *Streaming) connecting(id, key string) (chan string, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	ret := make(chan string)
	if _, ok := s.channels[id]; !ok {
		s.channels[id] = make(map[string]chan string)
	}
	if _, ok := s.channels[id][key]; ok {
		close(ret)
		return nil, fmt.Errorf("connection id(%s)/key(%s) has connected", id, key)
	}
	s.channels[id][key] = ret
	return ret, nil
}

func (s *Streaming) shutdown(id, key string) error {
	s.locker.Lock()
	defer s.locker.Unlock()

	_, ok := s.channels[id]
	if !ok {
		return fmt.Errorf("no connection")
	}
	c, ok := s.channels[id][key]
	if !ok {
		return fmt.Errorf("no connection")
	}

	close(c)

	delete(s.channels[id], key)
	if len(s.channels[id]) == 0 {
		delete(s.channels, id)
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
	i := 0
	for _, c := range conns {
		ret[i] = c
		i++
	}
	return ret
}

func (s *Streaming) getConn(id, key string) chan string {
	s.locker.Lock()
	defer s.locker.Unlock()

	conns, ok := s.channels[id]
	if !ok {
		return nil
	}
	c, ok := conns[key]
	if !ok {
		return nil
	}

	return c
}

func (s *Streaming) send(c chan string, content string) {
	defer func() {
		recover()
	}()

	c <- content
}
