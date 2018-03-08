package main

import (
	"fmt"
	"log"
)

func main() {
	api := NewApi("https://bb.otcbtc.com")
	err := api.Refresh()
	if err != nil {
		log.Println(err.Error())
	}
	price, err := api.GetPrice("eos", "cny")
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println("price:", price)
}
