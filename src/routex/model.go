package routex

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"logger"
	"math"
	"model"
	"time"
)

func Distance(latA, lngA, latB, lngB float64) float64 {
	x := math.Cos(latA*math.Pi/180) * math.Cos(latB*math.Pi/180) * math.Cos((lngA-lngB)*math.Pi/180)
	y := math.Sin(latA*math.Pi/180) * math.Sin(latB*math.Pi/180)
	s := x + y
	if s > 1 {
		s = 1
	}
	if s < -1 {
		s = -1
	}
	alpha := math.Acos(s)
	distance := alpha * 6371000
	return distance
}

type TutorialData struct {
	Offset    int64   `json:"offset"`
	Accuracy  float64 `json:"acc"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type SimpleLocation struct {
	Timestamp int64     `json:"t,omitempty"`
	GPS       []float64 `json:"gps,omitempty"` // latitude, longitude, accuracy
}

func (l *SimpleLocation) ToMars(c GeoConversionRepo) {
	if len(l.GPS) < 2 {
		return
	}
	gps := make([]float64, len(l.GPS))
	for i, p := range l.GPS {
		gps[i] = p
	}
	l.GPS = gps
	lat, lng := l.GPS[0], l.GPS[1]
	lat, lng = c.EarthToMars(lat, lng)
	l.GPS[0], l.GPS[1] = lat, lng
}

func (l *SimpleLocation) ToEarth(c GeoConversionRepo) {
	if len(l.GPS) < 2 {
		return
	}
	gps := make([]float64, len(l.GPS))
	for i, p := range l.GPS {
		gps[i] = p
	}
	lat, lng := l.GPS[0], l.GPS[1]
	lat, lng = c.MarsToEarth(lat, lng)
	l.GPS[0], l.GPS[1] = lat, lng
}

func (l *SimpleLocation) MarshalJSON() ([]byte, error) {
	gps := bytes.NewBuffer(nil)
	for _, p := range l.GPS {
		if _, err := gps.WriteString(fmt.Sprintf("%.7f,", p)); err != nil {
			return nil, err
		}
	}

	ret := bytes.NewBuffer(nil)
	if _, err := ret.WriteString(fmt.Sprintf(`{"t":%d`, l.Timestamp)); err != nil {
		return nil, err
	}
	if l := gps.Len(); l > 0 {
		gps.Truncate(l - 1)
		if _, err := ret.WriteString(`,"gps":[`); err != nil {
			return nil, err
		}
		if _, err := ret.Write(gps.Bytes()); err != nil {
			return nil, err
		}
		if _, err := ret.WriteString(`]`); err != nil {
			return nil, err
		}
	}
	if _, err := ret.WriteString("}"); err != nil {
		return nil, err
	}
	return ret.Bytes(), nil
}

type Geomark struct {
	Id          string           `json:"id,omitempty"`
	Type        string           `json:"type,omitempty"`
	Action      string           `json:"action,omitempty"`
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
	g.convert(c.EarthToMars)
}

func (g *Geomark) ToEarth(c GeoConversionRepo) {
	g.convert(c.MarsToEarth)
}

func (g *Geomark) convert(f func(lat, lng float64) (float64, float64)) {
	switch g.Type {
	case "location":
		g.Latitude, g.Longitude = f(g.Latitude, g.Longitude)
	case "route":
		pos := make([]SimpleLocation, len(g.Positions))
		for i, p := range g.Positions {
			if len(p.GPS) >= 2 {
				gps := make([]float64, len(p.GPS))
				for i, d := range p.GPS {
					gps[i] = d
				}
				gps[0], gps[1] = f(p.GPS[0], p.GPS[1])
				p.GPS = gps
			}
			pos[i] = p
		}
		g.Positions = pos
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

type Routex struct {
	CrossId   int64 `json:"cross_id,omitempty"`
	Enable    bool  `json:"enable,omitempty"`
	UpdatedAt int64 `json:"updated_at, omitempty"`
}

type RoutexControl interface {
	EnableCross(userId, crossId int64, afterInSecond int) error
	DisableCross(userId, crossId int64) error
}

type RoutexRepo interface {
	RoutexControl
	Search(crossIds []int64) ([]Routex, error)
	Get(userId, crossId int64) (*Routex, error)
	Update(userId, crossId int64) error
}

type BreadcrumbCache interface {
	RoutexControl
	Save(userId int64, l SimpleLocation) error
	Load(userId int64) (SimpleLocation, error)
	SaveCross(userId int64, l SimpleLocation) (cross_ids []int64, err error)
	LoadCross(userId, crossId int64) (SimpleLocation, bool, error)
}

type BreadcrumbsRepo interface {
	RoutexControl
	Save(userId int64, l SimpleLocation) error
	Load(userId, crossId, afterTimestamp int64) ([]SimpleLocation, error)
	Update(userId int64, l SimpleLocation) error
}

type GeomarksRepo interface {
	Set(crossId int64, mark Geomark) error
	Get(crossId int64) ([]Geomark, error)
	Delete(crossId int64, type_, id string) error
}

type GeoConversionRepo interface {
	EarthToMars(lat, long float64) (float64, float64)
	MarsToEarth(lat, long float64) (float64, float64)
}

const (
	ROUTEX_SETUP_INSERT      = "INSERT IGNORE INTO `routex` (`user_id`, `cross_id`, `enable`, `updated_at`) VALUES(?, ?, ?, UNIX_TIMESTAMP())"
	ROUTEX_SETUP_UPDATE      = "UPDATE `routex` SET `enable`=?, `updated_at`=UNIX_TIMESTAMP() WHERE `user_id`=? AND `cross_id`=?"
	ROUTEX_SETUP_SEARCH      = "SELECT `cross_id`, `enable`, `updated_at` FROM `routex` WHERE `cross_id` IN (%s) GROUP By `cross_id` ORDER BY `updated_at` DESC"
	ROUTEX_SETUP_GET         = "SELECT `cross_id`, `enable`, `updated_at` FROM `routex` WHERE `user_id`=? AND `cross_id`=? LIMIT 1"
	ROUTEX_SETUP_ONLY_UPDATE = "UPDATE `routex` SET `updated_at`=UNIX_TIMESTAMP() WHERE `user_id`=? AND `cross_id`=?"
)

type RoutexSaver struct {
	db *sql.DB
}

func NewRoutexSaver(db *sql.DB) *RoutexSaver {
	return &RoutexSaver{
		db: db,
	}
}

func (s *RoutexSaver) EnableCross(userId, crossId int64, afterInSecond int) error {
	res, err := s.db.Exec(ROUTEX_SETUP_UPDATE, true, userId, crossId)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		if _, err := s.db.Exec(ROUTEX_SETUP_INSERT, userId, crossId, true); err != nil {
			return err
		}
	}
	return nil
}

func (s *RoutexSaver) DisableCross(userId, crossId int64) error {
	res, err := s.db.Exec(ROUTEX_SETUP_UPDATE, false, userId, crossId)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		if _, err := s.db.Exec(ROUTEX_SETUP_INSERT, userId, crossId, false); err != nil {
			return err
		}
	}
	return nil
}

func (s *RoutexSaver) Search(crossIds []int64) ([]Routex, error) {
	ids := ""
	for _, id := range crossIds {
		ids = fmt.Sprintf("%s,%d", ids, id)
	}
	ids = ids[1:]
	sql := fmt.Sprintf(ROUTEX_SETUP_SEARCH, ids)
	row, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	var ret []Routex
	for row.Next() {
		var routex Routex
		if err := row.Scan(&routex.CrossId, &routex.Enable, &routex.UpdatedAt); err != nil {
			return nil, err
		}
		ret = append(ret, routex)
	}
	return ret, nil
}

func (s *RoutexSaver) Get(userId, crossId int64) (*Routex, error) {
	row, err := s.db.Query(ROUTEX_SETUP_GET, userId, crossId)
	if err != nil {
		return nil, err
	}
	if !row.Next() {
		return nil, nil
	}
	var ret Routex
	if err := row.Scan(&ret.CrossId, &ret.Enable, &ret.UpdatedAt); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (s *RoutexSaver) Update(userId, crossId int64) error {
	_, err := s.db.Exec(ROUTEX_SETUP_ONLY_UPDATE, userId, crossId)
	if err != nil {
		return err
	}
	return nil
}

const (
	BREADCRUMBS_UPDATE_START = "UPDATE `breadcrumbs_windows` SET `end_at`=UNIX_TIMESTAMP()+? WHERE `user_id`=? AND `cross_id`=? AND `end_at`>=UNIX_TIMESTAMP()"
	BREADCRUMBS_INSERT_START = "INSERT INTO `breadcrumbs_windows` (`user_id`, `cross_id`, `start_at`, `end_at`) VALUES(?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()+?)"
	BREADCRUMBS_UPDATE_END   = "UPDATE `breadcrumbs_windows` SET `end_at`=UNIX_TIMESTAMP()-1 WHERE `user_id`=? AND `cross_id`=? AND `end_at`>=UNIX_TIMESTAMP()"
	BREADCRUMBS_SAVE         = "INSERT INTO `breadcrumbs` (`user_id`, `lat`, `lng`, `acc`, `timestamp`) VALUES(?, ?, ?, ?, UNIX_TIMESTAMP());"
	BREADCRUMBS_GET          = "SELECT b.lat, b.lng, b.acc, b.timestamp FROM breadcrumbs AS b, breadcrumbs_windows AS w WHERE b.user_id=w.user_id AND b.timestamp BETWEEN w.start_at AND w.end_at AND w.user_id=? AND w.cross_id=? AND b.timestamp<=? ORDER BY b.timestamp DESC LIMIT 100"
	BREADCRUMBS_UPDATE       = "UPDATE `breadcrumbs` SET lat=?, lng=?, acc=?, timestamp=UNIX_TIMESTAMP() WHERE user_id=? ORDER BY timestamp DESC LIMIT 1"
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
	if len(l.GPS) < 3 {
		return fmt.Errorf("invalid simple location: %s", l)
	}
	if _, err := s.db.Exec(BREADCRUMBS_SAVE, userId, l.GPS[0], l.GPS[1], l.GPS[2]); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) Update(userId int64, l SimpleLocation) error {
	if len(l.GPS) < 3 {
		return fmt.Errorf("invalid simple location: %s", l)
	}
	if _, err := s.db.Exec(BREADCRUMBS_UPDATE, l.GPS[0], l.GPS[1], l.GPS[2], userId); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) Load(userId, crossId, afterTimestamp int64) ([]SimpleLocation, error) {
	rows, err := s.db.Query(BREADCRUMBS_GET, userId, crossId, afterTimestamp)
	if err != nil {
		return nil, err
	}
	var ret []SimpleLocation
	for rows.Next() {
		var l SimpleLocation
		var lat, lng, accuracy float64
		err := rows.Scan(&lat, &lng, &accuracy, &l.Timestamp)
		if err != nil {
			return nil, err
		}
		l.GPS = []float64{lat, lng, accuracy}
		ret = append(ret, l)
	}
	return ret, nil
}

type BreadcrumbCacheSaver struct {
	r          *redis.Pool
	saveScript *redis.Script
}

func NewBreadcrumbCacheSaver(r *redis.Pool) *BreadcrumbCacheSaver {
	ret := &BreadcrumbCacheSaver{
		r: r,
	}
	ret.saveScript = redis.NewScript(1, `
		local user_id = KEYS[1]
		local data = ARGV[1]
		local now = ARGV[2]
		local userkey = "exfe:v3:routex:user_"..user_id
		local matchkey = "exfe:v3:routex:user_"..user_id..":cross"
		local crosses = redis.call("ZRANGEBYSCORE", matchkey, now, "+INF")
		local ret = {}
		redis.call("EXPIRE", userkey, 600)
		for i = 1, #crosses do
			local c = crosses[i]
			redis.call("SET", matchkey..":"..c, data, "EX", "7200")
			table.insert(ret, c)
		end
		return ret
	`)
	return ret
}

func (s *BreadcrumbCacheSaver) ukey(userId int64) string {
	return fmt.Sprintf("exfe:v3:routex:user_%d", userId)
}

func (s *BreadcrumbCacheSaver) ckey(crossId, userId int64) string {
	return fmt.Sprintf("exfe:v3:routex:user_%d:cross:%d", userId, crossId)
}

func (s *BreadcrumbCacheSaver) EnableCross(userId, crossId int64, afterInSecond int) error {
	key, conn := s.ukey(userId)+":cross", s.r.Get()
	defer conn.Close()

	till := time.Now().Add(time.Duration(afterInSecond) * time.Second).Unix()
	if err := conn.Send("ZADD", key, till, crossId); err != nil {
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

func (s *BreadcrumbCacheSaver) DisableCross(userId, crossId int64) error {
	key, conn := s.ukey(userId)+":cross", s.r.Get()
	defer conn.Close()

	till := time.Now().Unix()
	if _, err := conn.Do("ZADD", key, till, crossId); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbCacheSaver) SaveCross(userId int64, l SimpleLocation) ([]int64, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	conn := s.r.Get()
	defer conn.Close()

	reply, err := redis.Values(s.saveScript.Do(conn, userId, b, time.Now().Unix()))
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

func (s *BreadcrumbCacheSaver) LoadCross(userId, crossId int64) (SimpleLocation, bool, error) {
	key, conn := s.ckey(crossId, userId), s.r.Get()
	defer conn.Close()

	var ret SimpleLocation
	reply, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
		return ret, false, nil
	}
	if err != nil {
		return ret, false, err
	}
	if err := json.Unmarshal(reply, &ret); err != nil {
		logger.ERROR("can't unmashal location value: %s with %s", err, string(reply))
		return ret, false, err
	}
	return ret, true, nil
}

func (s *BreadcrumbCacheSaver) Save(userId int64, l SimpleLocation) error {
	key, conn := s.ukey(userId), s.r.Get()
	defer conn.Close()

	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	if _, err := conn.Do("SET", key, b, "EX", 600); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbCacheSaver) Load(userId int64) (SimpleLocation, error) {
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

const (
	GEOMARKS_CREATE = "INSERT IGNORE INTO `geomarks` (`id`, `type`, `cross_id`, `mark`, `touched_at`, `deleted`) VALUES (?, ?, ?, ?, UNIX_TIMESTAMP(), FALSE)"
	GEOMARKS_UPDATE = "UPDATE `geomarks` SET `mark`=?, `touched_at`=UNIX_TIMESTAMP(), `deleted`=FALSE WHERE `id`=? AND `type`=? AND `cross_id`=?"
	GEOMARKS_GET    = "SELECT  `mark` FROM `geomarks` WHERE `cross_id`=? AND `deleted`=FALSE"
	GEOMARKS_DELETE = "UPDATE `geomarks` SET `deleted`=TRUE, `touched_at`=UNIX_TIMESTAMP() WHERE `id`=? AND `type`=? AND `cross_id`=? AND `deleted`=FALSE"
)

type GeomarksSaver struct {
	Db *sql.DB
}

func (s *GeomarksSaver) Set(crossId int64, mark Geomark) error {
	b, err := json.Marshal(mark)
	if err != nil {
		return err
	}
	n, err := s.Db.Exec(GEOMARKS_CREATE, mark.Id, mark.Type, crossId, string(b))
	if err != nil {
		return err
	}
	ret, err := n.RowsAffected()
	if err != nil {
		return err
	}
	if ret == 0 {
		mark.CreatedBy, mark.CreatedAt = mark.UpdatedBy, mark.UpdatedAt
		b, err = json.Marshal(mark)
		if err != nil {
			return err
		}
		_, err := s.Db.Exec(GEOMARKS_UPDATE, string(b), mark.Id, mark.Type, crossId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GeomarksSaver) Get(crossId int64) ([]Geomark, error) {
	rows, err := s.Db.Query(GEOMARKS_GET, crossId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Geomark
	for rows.Next() {
		var b string
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		var mark Geomark
		if err := json.Unmarshal([]byte(b), &mark); err != nil {
			return nil, err
		}
		ret = append(ret, mark)
	}
	return ret, nil
}

func (s *GeomarksSaver) Delete(crossId int64, markType, markId string) error {
	if _, err := s.Db.Exec(GEOMARKS_DELETE, markId, markType, crossId); err != nil {
		return err
	}
	return nil
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
