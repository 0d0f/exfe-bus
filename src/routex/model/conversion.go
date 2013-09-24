package rmodel

import (
	"database/sql"
	"github.com/googollee/eviltransform/go"
	"model"
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
}

func NewGeoConversion(config *model.Config, db *sql.DB) (*GeoConversion, error) {
	return &GeoConversion{}, nil
}

func (c *GeoConversion) MarsToEarth(lat, lng float64) (float64, float64) {
	return transform.GCJtoWGS(lat, lng)
}

func (c *GeoConversion) EarthToMars(lat, lng float64) (float64, float64) {
	return transform.WGStoGCJ(lat, lng)
}
