package here

import (
	"math"
	"time"
)

type Group struct {
	Name            string
	CenterLatitude  float64
	CenterLongitude float64
	HasGps          bool
	Traits          map[string]int
	Data            map[string]*Data
}

func NewGroup() *Group {
	return &Group{
		Traits: make(map[string]int),
		Data:   make(map[string]*Data),
	}
}

func (g *Group) Add(data *Data) {
	g.Data[data.Token] = data
	g.calcuate()
}

func (g *Group) Remove(data *Data) {
	delete(g.Data, data.Token)
	g.calcuate()
}

func (g *Group) Clear(limit time.Duration) []*Data {
	data := make(map[string]*Data)
	var cleared []*Data
	for token, d := range g.Data {
		if d.UpdatedAt.Add(limit).Before(time.Now()) {
			cleared = append(cleared, d)
		} else {
			data[token] = d
		}
	}
	g.Data = data
	g.calcuate()
	return cleared
}

func (g *Group) Distance(u *Data) float64 {
	if !u.HasGPS() || !g.HasGps {
		return -1
	}
	a := math.Cos(g.CenterLatitude) * math.Cos(u.latitude) * math.Cos(g.CenterLongitude-u.longitude)
	b := math.Sin(g.CenterLatitude) * math.Sin(u.latitude)
	return math.Acos(a + b)
}

func (g *Group) HasTraits(traits []string) bool {
	for _, t := range traits {
		if _, ok := g.Traits[t]; ok {
			return true
		}
	}
	return false
}

func (g *Group) calcuate() {
	g.HasGps, g.CenterLatitude, g.CenterLongitude, g.Traits = false, 0, 0, make(map[string]int)
	n := 0
	for _, u := range g.Data {
		if u.HasGPS() {
			if n == 0 {
				g.CenterLatitude = u.latitude
				g.CenterLongitude = u.longitude
			} else {
				a := u.accuracy
				coeff := float64(n) * a
				g.CenterLatitude = (coeff*g.CenterLatitude + u.latitude) / (coeff + 1)
				g.CenterLongitude = (coeff*g.CenterLongitude + u.longitude) / (coeff + 1)
			}
			n += 1
			g.HasGps = true
		}
		for _, t := range u.Traits {
			g.Traits[t] += 1
		}
	}
}
