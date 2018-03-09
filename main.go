package main

import (
	"log"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

func main() {
	cfg := Config{}
	clt, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     cfg.Influx.Addr,
		Username: cfg.Influx.User,
		Password: cfg.Influx.Pwd,
	})
	if err != nil {
		panic(err.Error())
	}
	defer clt.Close()
	api := NewApi("https://bb.otcbtc.com")
	for {
		time.Sleep(time.Duration(cfg.Sleep) * time.Second)
		err = api.Refresh()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		tags := make(map[string]string)
		bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  cfg.Influx.DB,
			Precision: "s",
		})
		var pt *client.Point
		for k, v := range api.Datas() {
			tags["chain"] = k
			pt, err = client.NewPoint("block", tags, v.Ticker.ToMap(), time.Now())
			if err != nil {
				log.Println("create point error:", err.Error())
				continue
			}
			bp.AddPoint(pt)
			err = clt.Write(bp)
			if err != nil {
				log.Println("write point error:", err.Error())
			}
		}
	}
}
