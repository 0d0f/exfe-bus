package ringcache

import (
	"container/ring"
	"fmt"
)

type cacheData struct {
	key  string
	data interface{}
}

type RingCache struct {
	head *ring.Ring
}

func New(max int) *RingCache {
	return &RingCache{
		head: ring.New(max),
	}
}

func (c *RingCache) Push(key string, data interface{}) {
	r := c.findKey(key)
	switch r {
	case nil:
		c.head = c.head.Prev()
		c.head.Value.(*cacheData).key = key
		c.head.Value.(*cacheData).data = data
	case c.head:
		r.Value.(*cacheData).data = data
	default:
		r.Value.(*cacheData).data = data
		r.Prev().Unlink(1)
		c.head = c.head.Prev()
		c.head.Link(r)
		c.head = c.head.Next()
	}
}

func (c *RingCache) Get(key string) interface{} {
	data := c.findKey(key)
	if data == nil {
		return nil
	}
	return data.Value.(*cacheData).data
}

func (c *RingCache) findKey(key string) *ring.Ring {
	for i, n, r := 0, c.head.Len(), c.head; i < n; i, r = i+1, r.Next() {
		if r.Value == nil {
			r.Value = &cacheData{
				key:  key,
				data: nil,
			}
			return r
		}
		if r.Value.(*cacheData).key == key {
			return r
		}
	}
	return nil
}

func (c *RingCache) print() {
	for i, n, r := 0, c.head.Len(), c.head; i < n; i, r = i+1, r.Next() {
		if r.Value == nil {
			fmt.Printf("nil, ")
		} else {
			fmt.Printf("%s(%v), ", r.Value.(*cacheData).key, r.Value.(*cacheData).data)
		}
	}
	fmt.Println()
}
