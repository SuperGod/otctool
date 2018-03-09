package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Influx struct {
		Addr string
		User string
		Pwd  string
		DB   string `json:"db"`
	}
	Sleep int64 // sleep seconds
}

func LoadConfig(configFile string) (cfg *Config, err error) {
	f, err := os.Open(configFile)
	if err != nil {
		return
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	cfg = new(Config)
	err = json.Unmarshal(buf, cfg)
	return
}
