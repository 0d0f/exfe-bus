package timer

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type Server struct {
	pool            *redis.Pool
	prefix          string
	setName         string
	notifyName      string
	dataName        string
	timeoutInSecond int64

	sendChan chan int
	quitChan chan error
	conn     redis.Conn
}

func NewServer(pool *redis.Pool, prefix string, timeoutInSecond int64) *Server {
	return &Server{
		pool:            pool,
		prefix:          prefix,
		setName:         sortedSetname(prefix),
		notifyName:      notifyName(prefix),
		dataName:        dataName(prefix),
		timeoutInSecond: timeoutInSecond,
	}
}

func (s *Server) Serve() error {
	s.sendChan, s.quitChan = make(chan int), make(chan error)
	var err error
	if s.conn, err = s.pool.Dial(); err != nil {
		return err
	}
	go s.listenSend()

	for {
		wait, err := s.check()
		if err != nil {
			return err
		}
		select {
		case <-time.After(time.Second * time.Duration(wait)):
		case err := <-s.quitChan:
			return err
		case <-s.sendChan:
		}
	}
}

func (s *Server) Close() {
	s.conn.Close()
}

func (s *Server) listenSend() {
	defer s.conn.Close()

	psc := redis.PubSubConn{s.conn}
	if err := psc.Subscribe(s.notifyName); err != nil {
		s.quitChan <- err
		return
	}
	for {
		switch n := psc.Receive().(type) {
		case redis.Message:
			s.sendChan <- 1
		case redis.PMessage:
			s.sendChan <- 1
		case redis.Subscription:
			if n.Count == 0 {
				s.quitChan <- nil
				return
			}
		case error:
			s.quitChan <- n
			return
		}
	}
}

func (s *Server) check() (int64, error) {
	conn := s.pool.Get()
	defer conn.Close()

	for {
		values, err := redis.Values(conn.Do("ZRANGEBYSCORE", s.setName, "-INF", "+INF", "WITHSCORES", "LIMIT", "0", "1"))
		if err != nil {
			return 0, err
		}
		var name string
		var nextTime int64
		wait := s.timeoutInSecond
		if len(values) > 0 {
			if values, err = redis.Scan(values, &name, &nextTime); err != nil {
				return 0, err
			}
			now := time.Now().Unix()
			wait = nextTime - now
		}
		if wait > 0 {
			return wait, nil
		}
		if _, err := conn.Do("ZREM", s.setName, name); err != nil {
			return 0, err
		}
		data, err := redis.String(conn.Do("HGET", s.dataName, name))
		if err != nil {
			return 0, err
		}
		conn.Do("PUBLISH", s.prefix+name, data)
	}
}

func notifyName(prefix string) string {
	return prefix + ":notify"
}

func sortedSetname(prefix string) string {
	return prefix + ":timer"
}

func dataName(prefix string) string {
	return prefix + ":data"
}
