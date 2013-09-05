package rmodel

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	"model"
)

func NewRoutexModel(config *model.Config, db *sql.DB, pool *redis.Pool) (RoutexRepo, BreadcrumbCache, BreadcrumbsRepo, GeomarksRepo, GeoConversionRepo, error) {
	routexRepo, err := NewRoutexSaver(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	breadcrumbsCache := NewBreadcrumbCacheSaver(pool)
	breadcrumbsRepo, err := NewBreadcrumbsSaver(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	geomarksRepo, err := NewGeomarkSaver(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	conversion, err := NewGeoConversion(config, db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return routexRepo, breadcrumbsCache, breadcrumbsRepo, geomarksRepo, conversion, nil
}

type GeoConversionRepo interface {
	EarthToMars(lat, lng float64) (float64, float64)
	MarsToEarth(lat, lng float64) (float64, float64)
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
	GetWindowEnd(userId, crossId int64) (int64, error)
	Save(userId int64, l SimpleLocation) error
	Load(userId, crossId, afterTimestamp int64) ([]SimpleLocation, error)
	Update(userId int64, l SimpleLocation) error
}

type GeomarksRepo interface {
	Set(crossId int64, mark Geomark) error
	Get(crossId int64) ([]Geomark, error)
	Delete(crossId int64, type_, id string) error
}
