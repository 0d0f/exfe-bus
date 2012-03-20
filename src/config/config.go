package config

import (
	"os"
	"encoding/json"
)

func LoadFile(filename string, config interface{}) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		panic(err)
	}
}
