package notifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type Saver struct {
	pool *redis.Pool
}

func NewResponseSaver(pool *redis.Pool) *Saver {
	return &Saver{
		pool: pool,
	}
}

func (s *Saver) Save(id string, item ResponseItem, ontime int64) error {
	conn := s.pool.Get()
	defer conn.Close()

	b, err := json.Marshal(item)
	if err != nil {
		return err
	}
	err = conn.Send("SET", s.key(id), b)
	if err != nil {
		return err
	}
	err = conn.Send("EXPIREAT", s.key(id), ontime)
	if err != nil {
		return err
	}
	err = conn.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (s *Saver) Load(id string) (ResponseItem, error) {
	conn := s.pool.Get()
	defer conn.Close()

	var ret ResponseItem
	reply, err := redis.Bytes(conn.Do("GET", s.key(id)))
	if err != nil {
		return ret, err
	}
	if len(reply) == 0 {
		return ret, errors.New("not exist")
	}
	err = json.Unmarshal(reply, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (s *Saver) key(id string) string {
	return fmt.Sprintf("exfe:v3:poster:response:%s", id)
}
