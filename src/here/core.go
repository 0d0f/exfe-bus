package here

import (
	"math"
	"strconv"
	"time"
)

type Identity struct {
	ExternalID       string `json:"external_id"`
	ExternalUsername string `json:"external_username"`
	Provider         string `json:"provider"`
}

type Card struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Avatar     string     `json:"avatar"`
	Bio        string     `json:"bio"`
	IsMe       bool       `json:"is_me"`
	Identities []Identity `json:"identities"`
}

type Data struct {
	Token     string   `json:"token"`
	Latitude  string   `json:"latitude"`
	Longitude string   `json:"longitude"`
	Accuracy  string   `json:"accuracy"`
	Traits    []string `json:"traits"`
	Card      Card     `json:"card"`

	latitude  float64 `json:"-"`
	longitude float64 `json:"-"`
	accuracy  float64 `json:"-"`

	UpdatedAt time.Time `json:"-"`
}

func (d Data) HasGPS() bool {
	return d.Latitude != "" && d.Longitude != "" && d.Accuracy != ""
}

func (d *Data) Init() error {
	if d.HasGPS() {
		var err error
		d.latitude, err = strconv.ParseFloat(d.Latitude, 64)
		if err != nil {
			return err
		}
		d.longitude, err = strconv.ParseFloat(d.Longitude, 64)
		if err != nil {
			return err
		}
		d.accuracy, err = strconv.ParseFloat(d.Accuracy, 64)
		if err != nil {
			return err
		}
	} else {
		d.Latitude = ""
		d.Longitude = ""
		d.Accuracy = ""
	}
	return nil
}

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

func (g *Group) Clear(limit time.Duration) []string {
	data := make(map[string]*Data)
	var clearedTokens []string
	for token, d := range g.Data {
		if d.UpdatedAt.Add(limit).Before(time.Now()) {
			clearedTokens = append(clearedTokens, token)
		} else {
			data[token] = d
		}
	}
	g.Data = data
	g.calcuate()
	return clearedTokens
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

type Cluster struct {
	Groups     map[string]*Group
	TokenGroup map[string]string

	distantThreshold float64
	signThreshold    float64
	timeout          time.Duration
}

func NewCluster(threshold, signThreshold float64, timeout time.Duration) *Cluster {
	return &Cluster{
		Groups:           make(map[string]*Group),
		TokenGroup:       make(map[string]string),
		distantThreshold: threshold,
		signThreshold:    signThreshold,
		timeout:          timeout,
	}
}

func (c *Cluster) Add(data *Data) error {
	err := data.Init()
	if err != nil {
		return err
	}

	data.UpdatedAt = time.Now()
	groupId, ok := c.TokenGroup[data.Token]
	if ok {
		group := c.Groups[groupId]
		oldData := group.Data[data.Token]
		data.Card.Id = oldData.Card.Id
		group.Remove(data)
	}

	groupId = ""
	var distant float64 = -1
	for k, group := range c.Groups {
		d := group.Distance(data)
		if d < 0 {
			d = c.distantThreshold
		}
		if group.HasTraits(data.Traits) && d < c.signThreshold {
			groupId, distant = k, 0
		}
		if distant < 0 || d < distant {
			groupId, distant = k, d
		}
	}
	var group *Group
	if groupId != "" && distant < c.distantThreshold {
		group = c.Groups[groupId]
	} else {
		group = NewGroup()
		group.Name = data.Token
		groupId = data.Token
	}
	group.Add(data)
	c.Groups[groupId] = group
	c.TokenGroup[data.Token] = groupId
	return nil
}

func (c *Cluster) Clear() []string {
	groups := make(map[string]*Group)
	var clearedTokens []string
	for k, group := range c.Groups {
		clearedTokens = append(clearedTokens, group.Clear(c.timeout)...)
		if len(group.Data) > 0 {
			groups[k] = group
		}
	}
	c.Groups = groups
	return clearedTokens
}
