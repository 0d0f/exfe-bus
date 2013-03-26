package here

import (
	"math"
	"time"
)

type Identity struct {
	Id string `json:"id"`
}

type User struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Avatar     string     `json:"avatar"`
	Bio        string     `json:"bio"`
	Identities []Identity `json:"identities"`

	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Accuracy  float64  `json:"accuracy"`
	Traits    []string `json:"traits"`

	UpdatedAt time.Time `json:"-"`
}

type Group struct {
	Name            string
	CenterLatitude  float64
	CenterLongitude float64
	Traits          map[string]int
	Users           map[string]*User
}

func NewGroup() *Group {
	return &Group{
		Traits: make(map[string]int),
		Users:  make(map[string]*User),
	}
}

func (g *Group) Add(user *User) {
	user.UpdatedAt = time.Now()
	g.Users[user.Id] = user
	g.calcuate()
}

func (g *Group) Remove(user *User) {
	delete(g.Users, user.Id)
	g.calcuate()
}

func (g *Group) Clear(limit time.Duration) int {
	var remove []string
	for k, u := range g.Users {
		if u.UpdatedAt.Add(limit).Before(time.Now()) {
			remove = append(remove, k)
		}
	}
	for _, k := range remove {
		delete(g.Users, k)
	}
	g.calcuate()
	return len(remove)
}

func (g *Group) Distant(u *User) float64 {
	a := math.Cos(g.CenterLatitude) * math.Cos(u.Latitude) * math.Cos(g.CenterLongitude-u.Longitude)
	b := math.Sin(g.CenterLatitude) * math.Sin(u.Latitude)
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
	g.CenterLatitude, g.CenterLongitude, g.Traits = 0, 0, make(map[string]int)
	n := 0
	var userId string
	for k, u := range g.Users {
		if len(u.Traits) == 0 {
			a := u.Accuracy
			coeff := float64(n) * a
			g.CenterLatitude = (coeff*g.CenterLatitude + u.Latitude) / (coeff + 1)
			g.CenterLongitude = (coeff*g.CenterLongitude + u.Longitude) / (coeff + 1)
			n += 1
		}
		for _, t := range u.Traits {
			g.Traits[t] += 1
		}
		userId = k
	}
	if u, ok := g.Users[userId]; ok && (g.CenterLatitude == 0 || g.CenterLongitude == 0) {
		g.CenterLatitude, g.CenterLongitude = u.Latitude, u.Longitude
	}
}

type Cluster struct {
	Groups    map[string]*Group
	UserGroup map[string]string

	distantThreshold float64
	signThreshold    float64
	timeout          time.Duration
}

func NewCluster(threshold, signThreshold float64, timeout time.Duration) *Cluster {
	return &Cluster{
		Groups:           make(map[string]*Group),
		UserGroup:        make(map[string]string),
		distantThreshold: threshold,
		signThreshold:    signThreshold,
		timeout:          timeout,
	}
}

func (c *Cluster) AddUser(user *User) {
	groupKey := ""
	var distant float64 = -1
	for k, group := range c.Groups {
		d := group.Distant(user)
		if len(user.Traits) > 0 && d < c.signThreshold && group.HasTraits(user.Traits) {
			groupKey, distant = k, 0
		}
		if distant < 0 || d < distant {
			groupKey, distant = k, d
		}
	}
	var group *Group
	if groupKey != "" && distant < c.distantThreshold {
		group = c.Groups[groupKey]
	} else {
		group = NewGroup()
		groupKey = user.Id
	}
	group.Add(user)
	c.Groups[groupKey] = group
	c.UserGroup[user.Id] = groupKey
}

func (c *Cluster) Clear() []string {
	var remove []string
	var ret []string
	for k, group := range c.Groups {
		if group.Clear(c.timeout) > 0 {
			ret = append(ret, k)
		}
		if len(group.Users) == 0 {
			remove = append(remove, k)
		} else {
			c.Groups[k] = group
		}
	}
	for _, r := range remove {
		delete(c.Groups, r)
	}
	return ret
}
