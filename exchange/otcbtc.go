package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

type Ticker struct {
	Buy  string
	Sell string
	Low  string
	High string
	Last string
	Vol  string
}

func (t *Ticker) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["buy"], _ = strconv.ParseFloat(t.Buy, 64)
	data["sell"], _ = strconv.ParseFloat(t.Sell, 64)
	data["low"], _ = strconv.ParseFloat(t.Low, 64)
	data["high"], _ = strconv.ParseFloat(t.High, 64)
	data["last"], _ = strconv.ParseFloat(t.Last, 64)
	data["vol"], _ = strconv.ParseFloat(t.Vol, 64)
	return data
}

type TickerInfo struct {
	At     int64
	Ticker Ticker
}

type OTCBTC struct {
	Addr string
	data chan *client.Point
}

func NewOTCBTC(addr string) *OTCBTC {
	o := new(OTCBTC)
	o.Addr = addr
	o.data = make(chan *client.Point, 1024)
	return o
}

func (o *OTCBTC) Start() (err error) {
	go func() {
		o.Refresh()
		time.Sleep(time.Second * 30)
	}()
	return
}

func (o *OTCBTC) Message() (msg chan *client.Point) {
	return o.data
}

func (o *OTCBTC) Refresh() (err error) {
	apiURL := o.Addr + "/api/v2/tickers"
	resp, err := http.Get(apiURL)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("status %d", resp.StatusCode)
		return
	}
	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	cache := make(map[string]*TickerInfo)
	err = json.Unmarshal(buf, &cache)
	if err != nil {
		return
	}
	o.toPoint(cache)
	return
}

func (o *OTCBTC) toPoint(datas map[string]*TickerInfo) (err error) {
	var pt *client.Point
	tags := map[string]string{}
	for k, v := range datas {
		tags["chain"] = k
		pt, err = client.NewPoint("otcbtc", tags, v.Ticker.ToMap(), time.Now())
		if err != nil {
			log.Println("create point error:", err.Error())
			continue
		}
		o.data <- pt
	}
	return
}
