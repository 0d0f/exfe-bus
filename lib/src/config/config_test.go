package config

import (
	"testing"
)

const data = "a_string: abcdef\na_uint: 1231"

func TestConfig(t *testing.T) {
	c := LoadString(data)
	if c.String("a_string") != "abcdef" {
		t.Error("Load 'a_string' error, got:", c.String("a_string"))
	}
	if c.Uint("a_uint") != 1231 {
		t.Error("Load 'a_uint' error, got:", c.Uint("a_uint"))
	}
}
