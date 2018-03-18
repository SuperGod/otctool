package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Source struct {
	DB     string `json:"db"`
	Params map[string]string
}

type Config struct {
	Influx struct {
		Addr string
		User string
		Pwd  string
		DB   string `json:"db"`
	}
	Currency []string
	Sources  map[string]*Source
	Sleep    int64 // sleep seconds
	Batch    int
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
