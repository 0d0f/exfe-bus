package phone

import (
	"github.com/googollee/go-assert"
	"testing"
)

func TestFilter(t *testing.T) {
	type Test struct {
		codec string
		i     string
		c     string
		o     string
		ok    bool
	}
	var tests = []Test{
		{"gb2312", "测试emoji👿123", "", "测试emoji123", true},
		{"gb2312", "测试emoji👿123", "?", "测试emoji?123", true},
	}
	for i, test := range tests {
		o, err := filter(test.codec, test.i, test.c)
		assert.MustEqual(t, err == nil, test.ok, "test %d", i)
		assert.Equal(t, o, test.o, "test %d", i)
	}
}
