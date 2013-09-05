package rmodel

import (
	"broker"
	"database/sql"
	"encoding/json"
	"fmt"
	"logger"
	"model"
	"net/url"
	"ringcache"
)

const (
	GEOCONVERSION_GET      = "SELECT `offset_lat`, `offset_lng` FROM `gps_conversion` WHERE `lat`=FORMAT(?, ?) AND `lng`=FORMAT(?, ?) ORDER BY `accuracy` DESC"
	GEOCONVERSION_INSERT_2 = "INSERT IGNORE INTO `gps_conversion` (`lat`, `lng`, `offset_lat`, `offset_lng`, accuracy) VALUES(FORMAT(?, 2), FORMAT(?, 2), ?, ?, 2)"
)

type Offset struct {
	latOffset int
	lngOffset int
}

type GeoConversion struct {
	config  *model.Config
	db      *sql.DB
	get     *sql.Stmt
	insert2 *sql.Stmt
	cache   *ringcache.RingCache
}

func NewGeoConversion(config *model.Config, db *sql.DB) (*GeoConversion, error) {
	p := NewErrPrepare(db)
	ret := &GeoConversion{
		config:  config,
		db:      db,
		get:     p.Prepare(GEOCONVERSION_GET),
		insert2: p.Prepare(GEOCONVERSION_INSERT_2),
		cache:   ringcache.New(200),
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GeoConversion) cacheKey(lat, lng float64) string {
	return fmt.Sprintf("%.2f,%.2f", lat, lng)
}

func (c *GeoConversion) loadCache(lat, lng float64) *Offset {
	data := c.cache.Get(c.cacheKey(lat, lng))
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
	var offsetLat, offsetLng int
	if offset := c.loadCache(lat, lng); offset != nil {
		offsetLat, offsetLng = offset.latOffset, offset.lngOffset
	} else {
		row := c.get.QueryRow(lat, 2, lng, 2)
		if err := row.Scan(&offsetLat, &offsetLng); err != nil {
			row := c.get.QueryRow(lat, 1, lng, 1)
			if err := row.Scan(&offsetLat, &offsetLng); err != nil {
				if err != sql.ErrNoRows {
					logger.ERROR("geo_conversion offset for lat=%s, lng=%s is not int", lat, lng)
				}
				logger.DEBUG("no offset for %.7f, %.7f", lat, lng)
				return 0, 0
			}
			logger.DEBUG("low accuracy offset for %.7f, %.7f", lat, lng)
			go c.queryNavi(lat, lng)
		} else {
			logger.DEBUG("high accuracy offset for %.7f, %.7f", lat, lng)
			c.cache.Push(c.cacheKey(lat, lng), Offset{offsetLat, offsetLng})
		}
	}
	return float64(offsetLat) * 0.000001, float64(offsetLng) * 0.000001
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

func (c *GeoConversion) queryNavi(lat, lng float64) {
	logger.DEBUG("query high accuracy offset for %.7f, %.7f", lat, lng)
	query := make(url.Values)
	query.Set("locations", fmt.Sprintf("%.6f,%.6f", lng, lat))
	query.Set("coordsys", "gps")
	query.Set("output", "json")
	query.Set("key", c.config.AutoNavi.Key)
	u := fmt.Sprintf("http://restapi.amap.com/v3/assistant/coordinate/convert?%s", query.Encode())
	resp, err := broker.HttpResponse(broker.Http("GET", u, "", nil))
	if err != nil {
		logger.ERROR("get %s failed: %s", u, err)
		return
	}
	defer resp.Close()
	var ret struct {
		Info      string
		Locations string
	}
	decoder := json.NewDecoder(resp)
	if err := decoder.Decode(&ret); err != nil {
		logger.ERROR("decode %s failed: %s", u, err)
		return
	}
	if ret.Info != "ok" {
		logger.ERROR("response %s failed: %+v", u, ret)
		return
	}
	var naviLat, naviLng float64
	if _, err := fmt.Sscanf(ret.Locations, "%f,%f", &naviLng, &naviLat); err != nil {
		logger.ERROR("sscan %s failed: %s", ret.Locations, err)
		return
	}
	offsetLat, offsetLng := int((naviLat-lat)*1e6), int((naviLng-lng)*1e6)
	if offsetLat == 0 && offsetLng == 0 {
		return
	}
	if _, err := c.insert2.Exec(lat, lng, offsetLat, offsetLng); err != nil {
		logger.ERROR("insert geo conv failed: %s", err)
		return
	}
	logger.DEBUG("save high accuracy offset for %.7f, %.7f", lat, lng)
}
