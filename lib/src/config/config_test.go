package config

import (
	"testing"
)

const data = "a_string: abcdef\na_uint: 1231\na_int: 123"

func TestConfig(t *testing.T) {
	c := LoadString(data)
	if c.String("a_string") != "abcdef" {
		t.Error("Load 'a_string' error, got:", c.String("a_string"))
	}
	if c.Uint("a_uint") != 1231 {
		t.Error("Load 'a_uint' error, got:", c.Uint("a_uint"))
	}
	if c.Int("a_int") != 123 {
		t.Error("Load 'a_int' error, got:", c.Int("a_int"))
	}
}

func TestConfigString(t *testing.T) {
	data := "WithDoubleQuot: \"123\"\nWithSingleQuot: 'abc'"
	c := LoadString(data)
	if c.String("WithDoubleQuot") != "123" {
		t.Error("Load 'WithDoubleQuot' error, got:", c.String("WithDoubleQuot"))
	}
	if c.String("WithSingleQuot") != "abc" {
		t.Error("Load 'WithSingleQuot' error, got:", c.String("WithSingleQuot"))
	}
}
