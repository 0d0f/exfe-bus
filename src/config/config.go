package config

import (
	"flag"
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

func LoadHelper(defaultname string, config interface{}) {
	var configFile string
	flag.StringVar(&configFile, "config", "twitter.json", "Specify the configuration file")
	flag.Parse()
	LoadFile(configFile, config)
}
