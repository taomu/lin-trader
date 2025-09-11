package main

import (
	"fmt"

	"github.com/taomu/lin-trader/margin/data"
	"github.com/taomu/lin-trader/margin/exchange/binance"
)

func main() {
	api := &data.ApiInfo{
		Key:    "",
		Secret: "",
	}
	broker := binance.NewBroker(api)
	symbolInfos, err := broker.GetSymbolInfos()
	if err != nil {
		panic(err)
	}
	for _, symbolInfo := range symbolInfos {
		fmt.Printf("%v", symbolInfo)
	}
}
