package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"strconv"
)

type TickerInfo struct {
	At     int64
	Ticker struct {
		Buy  string
		Sell string
		Low  string
		High string
		Last string
		Vol  string
	}
}

type Api struct {
	host  string
	cache map[string]*TickerInfo
}

func NewApi(host string) *Api {
	api := new(Api)
	api.host = host
	return api
}

func (api *Api) Refresh() (err error) {
	apiURL := api.host + "/api/v2/tickers"
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
	api.cache = make(map[string]*TickerInfo)
	err = json.Unmarshal(buf, &api.cache)
	return
}

func (api *Api) GetPrice(from, to string) (price float64, err error) {
	key := from + "_" + to
	v, ok := api.cache[key]
	if !ok {
		err = fmt.Errorf("no such market: %s", key)
		return
	}
	price, err = strconv.ParseFloat(v.Ticker.Last, 64)
	return
}
