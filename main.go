package main

import (
	"flag"
	"log"

	"github.com/SuperGod/otctool/cfg"
	"github.com/SuperGod/otctool/exchange"
	client "github.com/influxdata/influxdb/client/v2"
)

var (
	configFile = flag.String("c", "config.json", "config file")
)

type Chainer interface {
	Start() error
	Message() chan *client.Point
	DB() string
}

func main() {
	flag.Parse()

	cfg, err := cfg.LoadConfig(*configFile)
	if err != nil {
		panic(err)
	}
	clt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     cfg.Influx.Addr,
		Username: cfg.Influx.User,
		Password: cfg.Influx.Pwd,
	})
	if err != nil {
		panic(err.Error())
	}
	defer clt.Close()
	huobiCfg, ok := cfg.Sources["huobi"]
	if !ok {
		panic("no huobi config found")
	}
	api, err := exchange.NewHuoBi(huobiCfg, cfg.Currency)
	if err != nil {
		panic(err.Error())
	}
	batch := cfg.Batch
	if batch == 0 {
		batch = 10
	}
	var bp client.BatchPoints
	bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  api.DB(),
		Precision: "s",
	})
	n := 0
	api.Start()
	for pt := range api.Message() {
		if n%batch == 0 {
			err = clt.Write(bp)
			if err != nil {
				log.Println("write point error:", err.Error())
				n++
				continue
			}
			bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
				Database:  api.DB(),
				Precision: "s",
			})
		}
		bp.AddPoint(pt)
		n++
	}
}
