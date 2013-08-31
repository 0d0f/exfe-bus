package routex

import (
	"github.com/googollee/go-assert"
	"testing"
)

func TestDistance(t *testing.T) {
	type Test struct {
		latA     float64
		lngA     float64
		latB     float64
		lngB     float64
		distance float64
	}
	var tests = []Test{
		{31.1773232, 121.5272407, 31.1774146, 121.5270696, 19.190085024966443},
	}
	for i, test := range tests {
		assert.Equal(t, Distance(test.latA, test.lngA, test.latB, test.lngB), test.distance, "test %d", i)
	}
}
