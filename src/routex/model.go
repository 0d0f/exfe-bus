package routex

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"logger"
	"math"
	"model"
	"strconv"
	"strings"
)

func redisScript(r redis.Conn, hash string, script string, args ...interface{}) (interface{}, error) {
	reply, err := r.Do("EVALSHA", append([]interface{}{hash}, args...)...)
	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT") {
		reply, err = r.Do("EVAL", append([]interface{}{script}, args...)...)
	}
	return reply, err
}

func Distance(latA, lngA, latB, lngB float64) float64 {
	x := math.Cos(latA*math.Pi/180) * math.Cos(latB*math.Pi/180) * math.Cos((lngA-lngB)*math.Pi/180)
	y := math.Sin(latA*math.Pi/180) * math.Sin(latB*math.Pi/180)
	alpha := math.Acos(x + y)
	distance := alpha * 6371000
	return distance
}

type SimpleLocation struct {
	Timestamp int64   `json:"ts,omitempty"`
	Accuracy  float64 `json:"acc,omitempty"`
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lng,omitempty"`
}

func (l *SimpleLocation) ToMars(c GeoConversionRepo) {
	l.Latitude, l.Longitude = c.EarthToMars(l.Latitude, l.Longitude)
}

func (l *SimpleLocation) ToEarth(c GeoConversionRepo) {
	l.Latitude, l.Longitude = c.MarsToEarth(l.Latitude, l.Longitude)
}

type Geomark struct {
	Id          string           `json:"id,omitempty"`
	Type        string           `json:"type,omitempty"`
	CreatedAt   int64            `json:"created_at,omitempty"`
	CreatedBy   string           `json:"created_by,omitempty"`
	UpdatedAt   int64            `json:"updated_at,omitempty"`
	UpdatedBy   string           `json:"updated_by,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Icon        string           `json:"icon,omitempty"`
	Title       string           `json:"title,omitempty"`
	Description string           `json:"description,omitempty"`
	Color       string           `json:"color,omitempty"`
	Accuracy    float64          `json:"acc,omitempty"`
	Latitude    float64          `json:"lat,omitempty"`
	Longitude   float64          `json:"lng,omitempty"`
	Positions   []SimpleLocation `json:"positions,omitempty"`
}

func (g *Geomark) ToMars(c GeoConversionRepo) {
	switch g.Type {
	case "location":
		g.Latitude, g.Longitude = c.EarthToMars(g.Latitude, g.Longitude)
	case "route":
		for i, p := range g.Positions {
			p.Latitude, p.Longitude = c.EarthToMars(p.Latitude, p.Longitude)
			g.Positions[i] = p
		}
	}
}

func (g *Geomark) ToEarth(c GeoConversionRepo) {
	switch g.Type {
	case "location":
		g.Latitude, g.Longitude = c.MarsToEarth(g.Latitude, g.Longitude)
	case "route":
		for i, p := range g.Positions {
			p.Latitude, p.Longitude = c.MarsToEarth(p.Latitude, p.Longitude)
			g.Positions[i] = p
		}
	}
}

type Token struct {
	TokenType  string `json:"token_type"`
	UserId     int64  `json:"user_id"`
	CrossId    uint64 `json:"cross_id"`
	IdentityId int64  `json:"identity_id"`

	Cross    model.Cross `json:"-"`
	Readonly bool        `json:"-"`
}

type BreadcrumbsRepo interface {
	EnableCross(userId int64, crossId uint64) error
	DisableCross(userId int64, crossId uint64) error
	Crosses(userId int64) ([]uint64, error)
	Save(userId int64, l SimpleLocation) error
	Load(userId int64) (Geomark, error)
}

type GeomarksRepo interface {
	Save(crossId uint64, content []Geomark) error
	Load(crossId uint64) ([]Geomark, error)
}

type GeoConversionRepo interface {
	EarthToMars(lat, long float64) (float64, float64)
	MarsToEarth(lat, long float64) (float64, float64)
}

type BreadcrumbsSaver struct {
	Redis *redis.Pool
}

func (s *BreadcrumbsSaver) bkey(userId int64) string {
	return fmt.Sprintf("exfe:v3:routex:breadcrumbs:user_%d", userId)
}

func (s *BreadcrumbsSaver) ckey(crossId interface{}, userId int64) string {
	switch crossId.(type) {
	case int64:
		return fmt.Sprintf("exfe:v3:routex:user_%d:cross:%d", userId, crossId)
	case uint64:
		return fmt.Sprintf("exfe:v3:routex:user_%d:cross:%d", userId, crossId)
	default:
		return fmt.Sprintf("exfe:v3:routex:user_%d:cross:%v", userId, crossId)
	}
}

func (s *BreadcrumbsSaver) EnableCross(userId int64, crossId uint64) error {
	key, conn := s.ckey(crossId, userId), s.Redis.Get()
	defer conn.Close()

	if err := conn.Send("SET", key, crossId); err != nil {
		return err
	}
	if err := conn.Send("EXPIRE", key, 7200); err != nil {
		return err
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) DisableCross(userId int64, crossId uint64) error {
	key, conn := s.ckey(crossId, userId), s.Redis.Get()
	defer conn.Close()

	if _, err := conn.Do("DELETE", key); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) Crosses(userId int64) ([]uint64, error) {
	key, prefixLen, conn := s.ckey("*", userId), len(s.ckey("", userId)), s.Redis.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("KEYS", key))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	var ret []uint64
	for len(values) > 0 {
		var crossKey string
		values, err = redis.Scan(values, &crossKey)
		if err != nil {
			logger.ERROR("can't parse to string: %s", err)
			continue
		}
		crossId, err := strconv.ParseUint(crossKey[prefixLen:], 10, 64)
		if err != nil {
			logger.ERROR("can't parse cross id %s: %s", crossKey[prefixLen:], err)
			continue
		}
		ret = append(ret, crossId)
	}
	return ret, nil
}

func (s *BreadcrumbsSaver) Save(userId int64, l SimpleLocation) error {
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	key, conn := s.bkey(userId), s.Redis.Get()
	defer conn.Close()

	if err := conn.Send("LPUSH", key, b); err != nil {
		return err
	}
	if err := conn.Send("EXPIRE", key, 7200); err != nil {
		return err
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) Load(userId int64) (Geomark, error) {
	key, conn := s.bkey(userId), s.Redis.Get()
	defer conn.Close()

	ret := Geomark{}
	lrange, err := conn.Do("LRANGE", key, 0, 100)
	if err != nil {
		return ret, err
	}

	if lrange == nil {
		return ret, nil
	}
	values, err := redis.Values(lrange, err)
	if err != nil {
		return ret, err
	}
	for len(values) > 0 {
		var b []byte
		values, err = redis.Scan(values, &b)
		if err != nil {
			return ret, err
		}
		var location SimpleLocation
		err := json.Unmarshal(b, &location)
		if err != nil {
			logger.ERROR("can't unmashal location value: %s with %s", err, string(b))
			continue
		}
		ret.Positions = append(ret.Positions, location)
	}
	ret.Id, ret.Type, ret.Tags = fmt.Sprintf("%d", userId), "route", []string{"breadcrumbs"}
	return ret, nil
}

const (
	GEOMARKS_INSERT = "INSERT IGNORE INTO `routex` (`cross_id`, `route`, `touched_at`) VALUES (?, ?, NOW())"
	GEOMARKS_UPDATE = "UPDATE `routex` SET `route`=?, `touched_at`=NOW() WHERE `cross_id`=?"
	GEOMARKS_GET    = "SELECT `route` FROM `routex` WHERE `cross_id`=?"
)

type GeomarksSaver struct {
	Db *sql.DB
}

func (s *GeomarksSaver) Save(crossId uint64, data []Geomark) error {
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

func (s *GeomarksSaver) Load(crossId uint64) ([]Geomark, error) {
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
	var ret []Geomark
	err = json.Unmarshal([]byte(value), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

const (
	GEOCONVERSION_GET = "SELECT `offset_lat`, `offset_long` FROM `gps_conversion` WHERE `lat`=? AND `long`=?"
)

type GeoConversion struct {
	Db *sql.DB
}

func (c *GeoConversion) Offset(lat, long float64) (float64, float64) {
	latI := int(lat * 10)
	longI := int(long * 10)
	row, err := c.Db.Query(GEOCONVERSION_GET, latI, longI)
	if err != nil {
		return 0, 0
	}
	defer row.Close()

	if !row.Next() {
		return 0, 0
	}
	var offsetLat, offsetLong int
	err = row.Scan(&offsetLat, &offsetLong)
	if err != nil {
		logger.ERROR("geo_conversion offset for lat=%s, long=%s is not int", lat, long)
		return 0, 0
	}
	return float64(offsetLat) * 0.0001, float64(offsetLong) * 0.0001
}

func (c *GeoConversion) MarsToEarth(lat, long float64) (float64, float64) {
	offsetLat, offsetLong := c.Offset(lat, long)
	lat = lat - offsetLat
	long = long - offsetLong
	return lat, long
}

func (c *GeoConversion) EarthToMars(lat, long float64) (float64, float64) {
	offsetLat, offsetLong := c.Offset(lat, long)
	lat = lat + offsetLat
	long = long + offsetLong
	return lat, long
}
