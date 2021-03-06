package exchange

import (
	"fmt"
	"log"
	"net/http"
	"time"

	goex "github.com/SuperGod/GoEx"
	"github.com/SuperGod/GoEx/huobi"
	"github.com/SuperGod/otctool/cfg"
	client "github.com/influxdata/influxdb/client/v2"
)

var (
	CurrencyMap = map[string]goex.CurrencyPair{"eos_usdt": goex.EOS_USDT,
		"eos_btc":  goex.EOS_BTC,
		"btc_usdt": goex.BTC_USDT}
)

type CommonApi struct {
	db            string
	api           goex.API
	data          chan *client.Point
	currencyPairs []goex.CurrencyPair
	errs          chan error
	batch         int
}

func NewHuoBi(source *cfg.Source, currency []string) (commonApi *CommonApi, err error) {
	key, ok := source.Params["key"]
	if !ok {
		err = fmt.Errorf("no key of huobi ")
		return
	}

	secret, ok := source.Params["secret"]
	if !ok {
		err = fmt.Errorf("no secret of huobi ")
		return
	}
	accountID, ok := source.Params["account"]
	if !ok {
		err = fmt.Errorf("no account id of huobi ")
		return
	}
	commonApi = NewHuoBiApiByKey(key, secret, accountID)
	commonApi.errs = make(chan error, 1024)
	for _, v := range currency {
		pair, ok := CurrencyMap[v]
		if !ok {
			log.Println("unsupport currency pair:", v)
			continue
		}
		commonApi.currencyPairs = append(commonApi.currencyPairs, pair)
	}
	commonApi.batch = source.Batch
	commonApi.db = source.DB
	return
}

func NewHuoBiApiByKey(accessKey, secretKey, accountID string) (commonApi *CommonApi) {
	clt := &http.Client{}
	commonApi = new(CommonApi)
	commonApi.api = huobi.NewHuobiProWithAddr(clt, accessKey, secretKey, accountID, "https://api.huobipro.com")
	commonApi.data = make(chan *client.Point, 1024)
	return
}

func (cApi *CommonApi) Start() (err error) {
	go func() {
		for {
			cApi.getTicker()
			// go cApi.getDepth()
			time.Sleep(time.Duration(cApi.batch) * time.Second)
		}
	}()
	go cApi.logError()
	return
}

func (cApi *CommonApi) logError() {
	for err := range cApi.errs {
		log.Println(err.Error())
	}
}

func (cApi *CommonApi) getTicker() {
	var t *goex.Ticker
	var err error
	var pt *client.Point
	tags := map[string]string{}
	for _, v := range cApi.currencyPairs {
		t, err = cApi.api.GetTicker(v)
		if err != nil {
			log.Println("get ticker failed:", err.Error(), v)
			cApi.errs <- err
			continue
		}
		tags["chain"] = v.String()
		sec, fields := Ticker2Map(t)
		pt, err = client.NewPoint("ticker", tags, fields, sec)
		if err != nil {
			cApi.errs <- fmt.Errorf("create point error:%s", err.Error())
			continue
		}
		cApi.data <- pt
	}
}

func (cApi *CommonApi) getDepth() {
	// var d *goex.Depth
	// var err error
	// for _, v := range cApi.currencyPairs {
	// d, err = cApi.api.GetDepth(10, v)
	// if err != nil {
	// 	cApi.errs <- err
	// 	continue
	// }
	// }
}

func (cApi *CommonApi) getKline() {
	// var kline []goex.Kline
	// var err error
	// for _, v := range cApi.currencyPairs {
	// d, err = cApi.api.GetKlineRecords(v, period int, size int, since int)(10, v)
	// if err != nil {
	// cApi.errs <- err
	// continue
	// }
	// }
}

func (cApi *CommonApi) Message() chan *client.Point {
	return cApi.data
}

func (cApi *CommonApi) DB() string {
	return cApi.db
}

func Ticker2Map(t *goex.Ticker) (date time.Time, data map[string]interface{}) {
	date = time.Unix(int64(t.Date/1000), 0)
	data = make(map[string]interface{})
	data["last"] = t.Last
	data["buy"] = t.Buy
	data["sell"] = t.Sell
	data["high"] = t.High
	data["low"] = t.Low
	data["vol"] = t.Vol
	return
}
