package ringcache

import (
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestRingCache(t *testing.T) {
	c := New(5)
	c.Push("1", "abc")
	c.Push("2", "123")
	c.Push("3", "123")
	c.Push("4", "123")
	c.Push("5", "123")
	c.Push("6", "123")

	keys := []string{"6", "5", "4", "3", "2"}
	for i, r := 0, c.head; i < c.head.Len(); i, r = i+1, r.Next() {
		assert.Equal(t, r.Value.(*cacheData).key, keys[i])
		assert.Equal(t, r.Value.(*cacheData).data, "123")
	}

	d := c.Get("3")
	assert.Equal(t, d, "123")

	d = c.Get("1")
	assert.Equal(t, d, nil)

	c.Push("3", "abc")
	d = c.Get("3")
	assert.Equal(t, d, "abc")
	keys = []string{"3", "6", "5", "4", "2"}
	for i, r := 0, c.head; i < c.head.Len(); i, r = i+1, r.Next() {
		assert.Equal(t, r.Value.(*cacheData).key, keys[i])
		if i > 0 {
			assert.Equal(t, r.Value.(*cacheData).data, "123")
		} else {
			assert.Equal(t, r.Value.(*cacheData).data, "abc")
		}
	}

	c.Push("1", "abc")
	d = c.Get("1")
	assert.Equal(t, d, "abc")
	d = c.Get("2")
	assert.Equal(t, d, nil)
	keys = []string{"1", "3", "6", "5", "4"}
	for i, r := 0, c.head; i < c.head.Len(); i, r = i+1, r.Next() {
		assert.Equal(t, r.Value.(*cacheData).key, keys[i])
		if i > 1 {
			assert.Equal(t, r.Value.(*cacheData).data, "123")
		} else {
			assert.Equal(t, r.Value.(*cacheData).data, "abc")
		}
	}
}
