package config

import (
	"fmt"
	"github.com/kylelemons/go-gypsy"
	"strconv"
	"strings"
)

type Configure struct {
	file *yaml.File
}

func LoadFile(filename string) *Configure {
	return &Configure{
		file: yaml.ConfigFile(filename),
	}
}

func LoadString(data string) *Configure {
	return &Configure{
		file: yaml.Config(data),
	}
}

func (c *Configure) String(key string) string {
	value, err := c.file.Get(key)
	if err != nil {
		panic(fmt.Sprintf("Load configure key(%s) error: %s", key, err.Error()))
	}
	return strings.Trim(value, "\"'")
}

func (c *Configure) Uint(key string) uint {
	value := c.String(key)
	i, err := strconv.ParseUint(value, 10, 0)
	if err != nil {
		panic(fmt.Sprintf("Configure key(%s)'s value(%s) can't convert to int: %s", key, value, err.Error()))
	}
	return uint(i)
}

func (c *Configure) Int(key string) int {
	value := c.String(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Configure key(%s)'s value(%s) can't convert to int: %s", key, value, err.Error()))
	}
	return i
}
