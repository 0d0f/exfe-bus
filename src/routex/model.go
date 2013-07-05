package routex

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"logger"
	"time"
)

type Location struct {
	Timestamp int64  `json:"timestamp"`
	Lng       string `json:"lng"`
	Lat       string `json:"lat"`
}

type CrossToken struct {
	TokenType  string `json:"token_type"`
	CrossId    uint64 `json:"cross_id"`
	IdentityId uint64 `json:"identity_id"`
	UserId     int64  `json:"user_id"`
	CreatedAt  int64  `json:"created_time"`
	UpdatedAt  int64  `json:"updated_time"`
}

type LocationRepo interface {
	Save(id string, crossId uint64, l Location) error
	Load(id string, crossId uint64) ([]Location, error)
}

type RouteRepo interface {
	Save(crossId uint64, content string) error
	Load(crossId uint64) (string, error)
}

type LocationSaver struct {
	Redis redis.Conn
}

func (s *LocationSaver) Save(id string, crossId uint64, l Location) error {
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	key := s.key(id, crossId)
	err = s.Redis.Send("LPUSH", key, b)
	if err != nil {
		return err
	}
	err = s.Redis.Send("EXPIRE", key, int(time.Hour*24/time.Second))
	if err != nil {
		return err
	}
	err = s.Redis.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (s *LocationSaver) Load(id string, crossId uint64) ([]Location, error) {
	key := s.key(id, crossId)
	values, err := redis.Values(s.Redis.Do("LRANGE", key, 0, 100))
	if err != nil {
		return nil, err
	}
	var ret []Location
	for len(values) > 0 {
		var b []byte
		values, err = redis.Scan(values, &b)
		if err != nil {
			return nil, err
		}
		var location Location
		err := json.Unmarshal(b, &location)
		if err != nil {
			logger.ERROR("can't unmashal location value: %s with %s", err, string(b))
			continue
		}
		ret = append(ret, location)
	}
	return ret, nil
}

func (s *LocationSaver) key(id string, crossId uint64) string {
	return fmt.Sprintf("exfe:v3:routex:cross_%d:location:%s", id, crossId)
}

type RouteSaver struct {
	Redis redis.Conn
}

func (s *RouteSaver) Save(crossId uint64, content string) error {
	key := s.key(crossId)
	err := s.Redis.Send("SET", key, content)
	if err != nil {
		return err
	}
	err = s.Redis.Send("EXPIRE", key, int(time.Hour*24*7/time.Second))
	if err != nil {
		return err
	}
	err = s.Redis.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (s *RouteSaver) Load(crossId uint64) (string, error) {
	key := s.key(crossId)
	content, err := redis.String(s.Redis.Do("GET", key))
	if err != nil {
		return "", err
	}
	return content, nil
}

func (s *RouteSaver) key(crossId uint64) string {
	return fmt.Sprintf("exfe:v3:routex:cross_%d:route", crossId)
}
