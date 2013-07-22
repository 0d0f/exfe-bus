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
	"time"
)

type Location struct {
	Id          string     `json:"id,omitempty"`
	Type        string     `json:"type,omitempty"`
	CreatedAt   int64      `json:"created_at,omitempty"`
	CreatedBy   string     `json:"created_by,omitempty"`
	UpdatedAt   int64      `json:"updated_at,omitempty"`
	UpdatedBy   string     `json:"updated_by,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Icon        string     `json:"icon,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Color       string     `json:"color,omitempty"`
	Timestamp   int64      `json:"timestamp,omitempty"`
	Accuracy    string     `json:"accuracy",omitempty`
	Longitude   string     `json:"longitude,omitempty"`
	Latitude    string     `json:"latitude,omitempty"`
	Positions   []Location `json:"positions,omitempty"`
}

func (a Location) Distance(b Location) (float64, error) {
	latA, lngA, _, err := a.GetGeo()
	if err != nil {
		return 0, err
	}
	latB, lngB, _, err := b.GetGeo()
	if err != nil {
		return 0, err
	}
	x := math.Cos(latA*math.Pi/180) * math.Cos(latB*math.Pi/180) * math.Cos((lngA-lngB)*math.Pi/180)
	y := math.Sin(latA*math.Pi/180) * math.Sin(latB*math.Pi/180)
	alpha := math.Acos(x + y)
	distance := alpha * 6371000
	return distance, nil
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

func (l *Location) ToMars(c GeoConversionRepo) {
	if l.Longitude != "" && l.Latitude != "" {
		l.Latitude, l.Longitude = c.EarthToMars(l.Latitude, l.Longitude)
	}
	for i := range l.Positions {
		l.Positions[i].ToMars(c)
	}
}

func (l *Location) ToEarth(c GeoConversionRepo) {
	if l.Longitude != "" && l.Latitude != "" {
		l.Latitude, l.Longitude = c.MarsToEarth(l.Latitude, l.Longitude)
	}
	for i := range l.Positions {
		l.Positions[i].ToEarth(c)
	}
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
	Save(crossId uint64, content []Location) error
	Load(crossId uint64) ([]Location, error)
}

type GeoConversionRepo interface {
	EarthToMars(lat, long string) (string, string)
	MarsToEarth(lat, long string) (string, string)
}

type BreadcrumbsSaver struct {
	Redis *redis.Pool
}

func (s *BreadcrumbsSaver) Save(id string, crossId uint64, l Location) error {
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	key := s.key(id, crossId)
	conn := s.Redis.Get()
	defer conn.Close()

	err = conn.Send("LPUSH", key, b)
	if err != nil {
		return err
	}
	err = conn.Send("EXPIRE", key, int(time.Hour*24/time.Second))
	if err != nil {
		return err
	}
	err = conn.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (s *BreadcrumbsSaver) Load(id string, crossId uint64) ([]Location, error) {
	key := s.key(id, crossId)
	conn := s.Redis.Get()
	defer conn.Close()

	lrange, err := conn.Do("LRANGE", key, 0, 100)
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

func (s *GeomarksSaver) Save(crossId uint64, data []Location) error {
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

func (s *GeomarksSaver) Load(crossId uint64) ([]Location, error) {
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
	var ret []Location
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

func (c *GeoConversion) MarsToEarth(lat, long string) (string, string) {
	latf, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return lat, long
	}
	longf, err := strconv.ParseFloat(long, 64)
	if err != nil {
		return lat, long
	}
	offsetLat, offsetLong := c.Offset(latf, longf)
	latf = latf - offsetLat
	longf = longf - offsetLong
	return fmt.Sprintf("%f", latf), fmt.Sprintf("%f", longf)
}

func (c *GeoConversion) EarthToMars(lat, long string) (string, string) {
	latf, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return lat, long
	}
	longf, err := strconv.ParseFloat(long, 64)
	if err != nil {
		return lat, long
	}
	offsetLat, offsetLong := c.Offset(latf, longf)
	latf = latf + offsetLat
	longf = longf + offsetLong
	return fmt.Sprintf("%f", latf), fmt.Sprintf("%f", longf)
}
