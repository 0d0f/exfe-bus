package here

import (
	"time"
)

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

func (c *Cluster) Add(data *Data) *Group {
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

	return group
}

func (c *Cluster) Clear() []Group {
	groups := make(map[string]*Group)
	var clearedGroups []Group
	for k, group := range c.Groups {
		datas := group.Clear(c.timeout)
		if len(datas) > 0 {
			for _, d := range datas {
				group := NewGroup()
				group.Data[d.Token] = d
				clearedGroups = append(clearedGroups, *group)
				delete(c.TokenGroup, d.Token)
			}
		}
		if len(group.Data) > 0 {
			clearedGroups = append(clearedGroups, *group)
			groups[k] = group
		}
	}
	c.Groups = groups
	return clearedGroups
}
