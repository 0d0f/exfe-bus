package routex

import (
	"broker"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/googollee/go-multiplexer"
	"logger"
	"model"
	"strconv"
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

func (l Location) GetGeo() (float64, float64, float64, error) {
	lat, err := strconv.ParseFloat(l.Latitude, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	lng, err := strconv.ParseFloat(l.Longitude, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	acc, err := strconv.ParseFloat(l.Accuracy, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	return lat, lng, acc, nil
}

type Token struct {
	TokenType  string `json:"token_type"`
	UserId     int64  `json:"user_id"`
	IdentityId int64  `json:"identity_id"`

	Cross    model.Cross    `json:"-"`
	Identity model.Identity `json:"-"`
}

type BreadcrumbsRepo interface {
	Save(id string, crossId uint64, l Location) error
	Load(id string, crossId uint64) ([]Location, error)
}

type GeomarksRepo interface {
	Save(crossId uint64, content []map[string]interface{}) error
	Load(crossId uint64) ([]map[string]interface{}, error)
}

type BreadcrumbsSaver struct {
	Redis *broker.RedisPool
}

func (s *BreadcrumbsSaver) Save(id string, crossId uint64, l Location) error {
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	key := s.key(id, crossId)
	e := s.Redis.Do(func(i multiplexer.Instance) {
		r := i.(*broker.RedisInstance_).Redis

		err = r.Send("LPUSH", key, b)
		if err != nil {
			return
		}
		err = r.Send("EXPIRE", key, int(time.Hour*24/time.Second))
		if err != nil {
			return
		}
		err = r.Flush()
		if err != nil {
			return
		}
	})
	if e != nil {
		return e
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) Load(id string, crossId uint64) ([]Location, error) {
	key := s.key(id, crossId)
	var lrange interface{}
	var err error
	e := s.Redis.Do(func(i multiplexer.Instance) {
		r := i.(*broker.RedisInstance_).Redis

		lrange, err = r.Do("LRANGE", key, 0, 100)
	})
	if e != nil {
		return nil, e
	}
	if err != nil {
		return nil, err
	}

	if lrange == nil {
		return nil, nil
	}
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

func (s *BreadcrumbsSaver) key(id string, crossId uint64) string {
	return fmt.Sprintf("exfe:v3:routex:cross_%d:location:%s", crossId, id)
}

const (
	GEOMARKS_INSERT = "INSERT IGNORE INTO `routex` (`cross_id`, `route`, `touched_at`) VALUES (?, ?, NOW())"
	GEOMARKS_UPDATE = "UPDATE `routex` SET `route`=?, `touched_at`=NOW() WHERE `cross_id`=?"
	GEOMARKS_GET    = "SELECT `route` FROM `routex` WHERE `cross_id`=?"
)

type GeomarksSaver struct {
	Db *sql.DB
}

func (s *GeomarksSaver) Save(crossId uint64, data []map[string]interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	n, err := s.Db.Exec(GEOMARKS_INSERT, crossId, string(b))
	if err != nil {
		return err
	}
	rows, err := n.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		_, err := s.Db.Exec(GEOMARKS_UPDATE, string(b), crossId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GeomarksSaver) Load(crossId uint64) ([]map[string]interface{}, error) {
	var row *sql.Rows
	row, err := s.Db.Query(GEOMARKS_GET, crossId)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	value, exist := "", false
	for row.Next() {
		row.Scan(&value)
		exist = true
	}
	if !exist {
		return nil, nil
	}
	var ret []map[string]interface{}
	err = json.Unmarshal([]byte(value), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
