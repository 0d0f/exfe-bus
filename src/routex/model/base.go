package rmodel

import (
	"fmt"
	"model"
)

type TutorialData struct {
	Offset    int64   `json:"offset"`
	Accuracy  float64 `json:"acc"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Token struct {
	TokenType  string `json:"token_type"`
	UserId     int64  `json:"user_id"`
	CrossId    uint64 `json:"cross_id"`
	IdentityId int64  `json:"identity_id"`

	Cross    model.Cross    `json:"-"`
	Identity model.Identity `json:"-"`
}

type Identity struct {
	model.Identity
	Type   string `json:"type,omitempty"`
	Action string `json:"action,omitempty"`
}

type Invitation struct {
	Identity      model.Identity `json:"identity,omitempty"`
	Notifications []string       `json:"notification_identities,omitempty"`
	Type          string         `json:"type,omitempty"`
	Action        string         `json:"action,omitempty"`
}

type SimpleLocation struct {
	Timestamp int64      `json:"t,omitempty"`
	GPS       [3]float64 `json:"gps,omitempty"` // latitude, longitude, accuracy
}

func (l *SimpleLocation) ToMars(c GeoConversionRepo) {
	l.convert(c.EarthToMars)
}

func (l *SimpleLocation) ToEarth(c GeoConversionRepo) {
	l.convert(c.MarsToEarth)
}

func (l *SimpleLocation) convert(f func(lat, lng float64) (float64, float64)) {
	gps := l.GPS
	gps[0], gps[1] = f(gps[0], gps[1])
	l.GPS = gps
}

func (l SimpleLocation) MarshalJSON() ([]byte, error) {
	ret := fmt.Sprintf(`{"t":%d,"gps":[%.7f,%.7f,%.0f]}`, l.Timestamp, l.GPS[0], l.GPS[1], l.GPS[2])
	return []byte(ret), nil
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

func (g *Geomark) HasTag(tag string) bool {
	for _, t := range g.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (g *Geomark) RemoveTag(tag string) bool {
	ret := false
	for i := len(g.Tags) - 1; i >= 0; i-- {
		if g.Tags[i] == tag {
			ret = true
			g.Tags = append(g.Tags[:i], g.Tags[i+1:]...)
		}
	}
	return ret
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
			p.convert(f)
			pos[i] = p
		}
		g.Positions = pos
	}
}
