package exchange

import (
	"fmt"
	"testing"

	"github.com/SuperGod/otctool/cfg"
)

var ()

func TestHuoBi(t *testing.T) {
	cfg, err := cfg.LoadConfig("../config.json")
	if err != nil {
		t.Fatal(err.Error())
	}
	huobiCfg, ok := cfg.Sources["huobi"]
	if !ok {
		t.Fatal("no huobi found")
	}
	huobiCfg.Params["key"] = ""
	huobiCfg.Params["secret"] = ""
	huobiCfg.Params["account"] = "15000135390"
	api, err := NewHuoBi(huobiCfg, cfg.Currency)
	if err != nil {
		t.Fatal(err.Error())
	}

	api.Start()
	for msg := range api.Message() {
		fmt.Println("msg:", msg)
		t.Log("msg:", msg.String())
	}
}
