package routex

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"logger"
	"math"
	"model"
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

func (l *SimpleLocation) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"ts":%d,"acc":%.0f,"lat":%.6f,"lng":%.6f}`, l.Timestamp, l.Accuracy, l.Latitude, l.Longitude)), nil
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

type BreadcrumbCrossControl interface {
	EnableCross(userId, crossId int64, afterInSecond int) error
	DisableCross(userId, crossId int64) error
}

type BreadcrumbCache interface {
	BreadcrumbCrossControl
	Save(userId int64, l SimpleLocation) (cross_ids []int64, err error)
	Load(userId int64) (SimpleLocation, error)
	LoadCross(userId, crossId int64) (SimpleLocation, error)
}

type BreadcrumbsRepo interface {
	BreadcrumbCrossControl
	Save(userId int64, l SimpleLocation) error
	Load(userId, crossId int64) ([]SimpleLocation, error)
}

type GeomarksRepo interface {
	Save(crossId uint64, content []Geomark) error
	Load(crossId uint64) ([]Geomark, error)
}

type GeoConversionRepo interface {
	EarthToMars(lat, long float64) (float64, float64)
	MarsToEarth(lat, long float64) (float64, float64)
}

const (
	BREADCRUMBS_UPDATE_START = "UPDATE `breadcrumbs_windows` (`end_at`) VALUES (UNIX_TIMESTAMP()+?) WHERE `user_id`=? AND `cross_id`=? AND `end_at`>UNIX_TIMESTAMP()"
	BREADCRUMBS_INSERT_START = "INSERT INTO `breadcrumbs_windows` (`user_id`, `cross_id`, `start_at`, `end_at`) VALUES(?, ?, ?, UNIX_TIMESTAMP()+?)"
	BREADCRUMBS_UPDATE_END   = "UPDATE `breadcrumbs_windows` SET `end_at`=UNIX_TIMESTAMP()-1 WHERE user_id=? AND cross_id=? AND end_at>=UNIX_TIMESTAMP()"
	BREADCRUMBS_SAVE         = "INSERT INTO `breadcrumbs` (`user_id`, `lat`, `lng`, `acc`, `timestamp`) values(?, ?, ?, ?, ?);"
	BREADCRUMBS_GET          = "SELECT b.lat, b.lng, b.acc, b.timestamp FROM breadcrumbs AS b, breadcrumbs_windows AS w WHERE b.user_id=w.user_id AND b.timestamp BETWEEN w.start_at AND w.end_at AND w.user_id=? AND w.cross_id=? LIMIT 100"
)

type BreadcrumbsSaver struct {
	db *sql.DB
}

func NewBreadcrumbsSaver(db *sql.DB) *BreadcrumbsSaver {
	ret := &BreadcrumbsSaver{
		db: db,
	}
	return ret
}

func (s *BreadcrumbsSaver) EnableCross(userId, crossId int64, afterInSecond int) error {
	res, err := s.db.Exec(BREADCRUMBS_UPDATE_START, afterInSecond, userId, crossId)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r > 0 {
		return nil
	}
	if _, err := s.db.Exec(BREADCRUMBS_INSERT_START, userId, crossId, afterInSecond); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) DisableCross(userId, crossId int64) error {
	if _, err := s.db.Exec(BREADCRUMBS_UPDATE_END, userId, crossId); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) Save(userId int64, l SimpleLocation) error {
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	if _, err := s.db.Exec(BREADCRUMBS_SAVE, userId, l.Latitude, l.Longitude, l.Accuracy, l.Timestamp); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) Load(userId int64, crossId int64) ([]SimpleLocation, error) {
	rows, err := s.db.Query(BREADCRUMBS_GET, userId, crossId)
	if err != nil {
		return nil, err
	}
	var ret []SimpleLocation
	for rows.Next() {
		var l SimpleLocation
		err := rows.Scan(&l.Latitude, &l.Longitude, &l.Accuracy, &l.Timestamp)
		if err != nil {
			return nil, err
		}
		ret = append(ret, l)
	}
	return ret, nil
}

type BreadcrumbsCacheSaver struct {
	r          *redis.Pool
	saveScript *redis.Script
}

func NewBreadcrumbsCacheSaver(r *redis.Pool) *BreadcrumbsCacheSaver {
	ret := &BreadcrumbsCacheSaver{
		r: r,
	}
	ret.saveScript = redis.NewScript(1, `
		local user_id = KEYS[1]
		local data = ARGV[1]
		local keys = redis.Do("KEYS", "exfe:v3:routex:user_"..user_id..":cross:*")
		local ret = {}
		local keyprefix = string.len(keys)
		redis.Do("SET", "exfe:v3:routex:user_"..user_id, data)
		for i = 1, #keys do
			local k = keys[i]
			redis.Do("SET", k, data)
			table.insert(ret, string.sub(k, keyprefix))
		end
		return ret
	`)
	return ret
}

func (s *BreadcrumbsCacheSaver) ukey(userId int64) string {
	return fmt.Sprintf("exfe:v3:routex:user_%d", userId)
}

func (s *BreadcrumbsCacheSaver) ckey(crossId, userId int64) string {
	return fmt.Sprintf("exfe:v3:routex:user_%d:cross:%d", userId, crossId)
}

func (s *BreadcrumbsCacheSaver) EnableCross(userId, crossId int64, afterInSecond int) error {
	key, conn := s.ckey(crossId, userId), s.r.Get()
	defer conn.Close()

	if err := conn.Send("SET", key, ""); err != nil {
		return err
	}
	if err := conn.Send("EXPIRE", key, afterInSecond); err != nil {
		return err
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsCacheSaver) DisableCross(userId, crossId int64) error {
	key, conn := s.ckey(crossId, userId), s.r.Get()
	defer conn.Close()

	if _, err := conn.Do("DELETE", key); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsCacheSaver) Save(userId int64, l SimpleLocation) ([]int64, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	conn := s.r.Get()
	defer conn.Close()

	reply, err := redis.Values(s.saveScript.Do(conn, userId, b))
	if err != nil {
		return nil, err
	}
	var ret []int64
	for len(reply) > 0 {
		var crossId int64
		reply, err = redis.Scan(reply, &crossId)
		if err != nil {
			return nil, err
		}
		ret = append(ret, crossId)
	}
	return ret, nil
}

func (s *BreadcrumbsCacheSaver) Load(userId int64) (SimpleLocation, error) {
	key, conn := s.ukey(userId), s.r.Get()
	defer conn.Close()

	var ret SimpleLocation
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return ret, err
	}
	if err := json.Unmarshal(reply, &ret); err != nil {
		logger.ERROR("can't unmashal location value: %s with %s", err, string(reply))
		return ret, err
	}
	return ret, nil
}

func (s *BreadcrumbsCacheSaver) LoadCross(userId, crossId int64) (SimpleLocation, error) {
	key, conn := s.ckey(crossId, userId), s.r.Get()
	defer conn.Close()

	var ret SimpleLocation
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return ret, err
	}
	if err := json.Unmarshal(reply, &ret); err != nil {
		logger.ERROR("can't unmashal location value: %s with %s", err, string(reply))
		return ret, err
	}
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

type Offset struct {
	latOffset int
	lngOffset int
}

type GeoConversion struct {
	db    *sql.DB
	cache map[string]Offset
}

func NewGeoConversion(db *sql.DB) *GeoConversion {
	return &GeoConversion{
		db:    db,
		cache: make(map[string]Offset),
	}
}

func (c *GeoConversion) Offset(lat, long float64) (float64, float64) {
	latI := int(lat * 10)
	longI := int(long * 10)
	key := fmt.Sprintf("%d,%d", latI, longI)
	var offsetLat, offsetLong int
	if offset, ok := c.cache[key]; ok {
		offsetLat, offsetLong = offset.latOffset, offset.lngOffset
	} else {
		row, err := c.db.Query(GEOCONVERSION_GET, latI, longI)
		if err != nil {
			return 0, 0
		}
		defer row.Close()

		if !row.Next() {
			return 0, 0
		}
		err = row.Scan(&offsetLat, &offsetLong)
		if err != nil {
			logger.ERROR("geo_conversion offset for lat=%s, long=%s is not int", lat, long)
			return 0, 0
		}
		c.cache[key] = Offset{offsetLat, offsetLong}
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
