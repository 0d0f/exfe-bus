package here

import (
	"math"
	"time"
)

type Identity struct {
	Id string
}

type User struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Avatar     string     `json:"avatar"`
	Bio        string     `json:"bio"`
	Identities []Identity `json:"identities"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  int     `json:"accuracy"`
	Sign      string  `json:"sign"`
	OldSign   string  `json:"old_sign"`

	UpdatedAt time.Time `json:"-"`
}

type Group struct {
	Name            string
	CenterLatitude  float64
	CenterLongitude float64
	Signs           map[string]int
	Users           map[string]*User
}

func NewGroup() *Group {
	return &Group{
		Signs: make(map[string]int),
		Users: make(map[string]*User),
	}
}

func (g *Group) Add(user *User) {
	user.UpdatedAt = time.Now()
	g.Users[user.Id] = user
	g.calcuateCenter()
	if user.Sign != "" {
		g.Signs[user.Sign] += 1
	}
	if user.OldSign != "" {
		i := g.Signs[user.OldSign]
		i -= 1
		if i <= 0 {
			delete(g.Signs, user.OldSign)
		}
	}
}

func (g *Group) Clear(limit time.Duration) int {
	var remove []string
	for k, u := range g.Users {
		if u.UpdatedAt.Add(limit).Before(time.Now()) {
			remove = append(remove, k)
		}
	}
	for _, k := range remove {
		user := g.Users[k]
		if user.Sign != "" {
			i := g.Signs[user.Sign]
			i -= 1
			if i <= 0 {
				delete(g.Signs, user.Sign)
			} else {
				g.Signs[user.Sign] = i
			}
		}
		delete(g.Users, k)
	}
	return len(remove)
}

func (g *Group) calcuateCenter() {
	g.CenterLatitude, g.CenterLongitude = 0, 0
	n := 0
	for _, u := range g.Users {
		a := u.Accuracy
		coeff := float64(a * n)
		g.CenterLatitude = (coeff*g.CenterLatitude + u.Latitude) / (coeff + 1)
		g.CenterLongitude = (coeff*g.CenterLongitude + u.Longitude) / (coeff + 1)
		n += 1
	}
}

func (g *Group) Distant(u *User) float64 {
	a := math.Cos(g.CenterLatitude) * math.Cos(u.Latitude) * math.Cos(g.CenterLongitude-u.Longitude)
	b := math.Sin(g.CenterLatitude) * math.Sin(u.Latitude)
	return math.Acos(a + b)
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
		if user.Sign != "" && d < c.signThreshold && group.Signs[user.Sign] > 0 {
			groupKey, distant = k, 0
		}
		if distant < 0 || d < distant {
			groupKey, distant = k, d
		}
	}
	group := NewGroup()
	if groupKey != "" && distant < c.distantThreshold {
		group = c.Groups[groupKey]
	} else {
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
