package routex

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"logger"
	"model"
	"time"
)

type Location struct {
	Id          string   `json:"id,omitempty"`
	Type        string   `json:"type,omitempty"`
	CreatedAt   int64    `json:"created_at,omitempty"`
	CreatedBy   string   `json:"created_by,omitempty"`
	UpdatedAt   int64    `json:"updated_at,omitempty"`
	UpdatedBy   string   `json:"updated_by,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Timestamp   int64    `json:"timestamp,omitempty"`
	Accuracy    string   `json:"accuracy",omitempty`
	Longitude   string   `json:"longitude"`
	Latitude    string   `json:"latitude"`
}

type Route struct {
	Id          string     `json:"id,omitempty"`
	Type        string     `json:"type,omitempty"`
	CreatedAt   int64      `json:"created_at,omitempty"`
	CreatedBy   string     `json:"created_by,omitempty"`
	UpdatedAt   int64      `json:"updated_at,omitempty"`
	UpdatedBy   string     `json:"updated_by,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Timestamp   int64      `json:"timestamp,omitempty"`
	Locations   []Location `json:"locations"`
}

type Token struct {
	TokenType string `json:"token_type"`
	UserId    int64  `json:"user_id"`

	Cross    model.Cross    `json:"-"`
	Identity model.Identity `json:"-"`
}

type LocationRepo interface {
	Save(id string, crossId uint64, l Location) error
	Load(id string, crossId uint64) ([]Location, error)
}

type RouteRepo interface {
	Save(crossId uint64, content []map[string]interface{}) error
	Load(crossId uint64) ([]map[string]interface{}, error)
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
	lrange, err := s.Redis.Do("LRANGE", key, 0, 100)
	if lrange == nil {
		return nil, nil
	}
	fmt.Println("location", lrange, err)
	values, err := redis.Values(lrange, err)
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
	return fmt.Sprintf("exfe:v3:routex:cross_%d:location:%s", crossId, id)
}

type RouteSaver struct {
	Redis redis.Conn
}

func (s *RouteSaver) Save(crossId uint64, data []map[string]interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	key := s.key(crossId)
	err = s.Redis.Send("SET", key, b)
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

func (s *RouteSaver) Load(crossId uint64) ([]map[string]interface{}, error) {
	key := s.key(crossId)
	get, err := s.Redis.Do("GET", key)
	if get == nil {
		return nil, nil
	}
	fmt.Println("route", get, err)
	b, err := redis.Bytes(get, err)
	if err != nil {
		return nil, err
	}
	var ret []map[string]interface{}
	err = json.Unmarshal(b, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *RouteSaver) key(crossId uint64) string {
	return fmt.Sprintf("exfe:v3:routex:cross_%d:route", crossId)
}
