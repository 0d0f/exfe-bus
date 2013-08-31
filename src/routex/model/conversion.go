package model

import (
	"database/sql"
	"fmt"
	"ringcache"
)

const (
	GEOCONVERSION_GET = "SELECT `offset_lat`, `offset_lng` FROM `gps_conversion` WHERE `lat`=? AND `lng`=? LIMIT 1"
)

type Offset struct {
	latOffset int
	lngOffset int
}

type GeoConversion struct {
	db    *sql.DB
	get   *sql.Stmt
	cache *ringcache.RingCache
}

func NewGeoConversion(db *sql.DB) (*GeoConversion, error) {
	p := NewErrPrepare(db)
	ret := &GeoConversion{
		db:    db,
		get:   p.Prepare(GEOCONVERSION_GET),
		cache: ringcache.New(200),
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GeoConversion) loadCache(key string) *Offset {
	data := c.cache.Get(key)
	if data == nil {
		return nil
	}
	offset, ok := data.(Offset)
	if !ok {
		return nil
	}
	return &offset
}

func (c *GeoConversion) Offset(lat, lng float64) (float64, float64) {
	latI := int(lat * 10)
	lngI := int(lng * 10)
	key := fmt.Sprintf("%d,%d", latI, lngI)

	var offsetLat, offsetLng int
	if offset := c.loadCache(key); offset != nil {
		offsetLat, offsetLng = offset.latOffset, offset.lngOffset
	} else {
		row := c.get.QueryRow(latI, lngI)
		if err := row.Scan(&offsetLat, &offsetLng); err != nil {
			return 0, 0
		}
		c.cache.Push(key, Offset{offsetLat, offsetLng})
	}
	return float64(offsetLat) * 0.0001, float64(offsetLng) * 0.0001
}

func (c *GeoConversion) MarsToEarth(lat, lng float64) (float64, float64) {
	offsetLat, offsetLong := c.Offset(lat, lng)
	lat = lat - offsetLat
	lng = lng - offsetLong
	return lat, lng
}

func (c *GeoConversion) EarthToMars(lat, lng float64) (float64, float64) {
	offsetLat, offsetLong := c.Offset(lat, lng)
	lat = lat + offsetLat
	lng = lng + offsetLong
	return lat, lng
}
