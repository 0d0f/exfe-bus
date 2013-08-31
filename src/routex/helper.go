package routex

import (
	"math"
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
