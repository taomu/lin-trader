package main

import (
	"fmt"

	"github.com/taomu/lin-trader/futures"
	"github.com/taomu/lin-trader/futures/data"
)

func main() {
	showOkExample()
	// showBnExample()

}
func showOkExample() {
	broker, _ := futures.NewBroker(data.PLAT_OKX, "", "", "")
	broker.Test()
	// premium, _ := broker.GetPremium()
	// fmt.Println(premium)
	tickers, _ := broker.GetTickers24h()
	fmt.Println(tickers)
}
func showBnExample() {
	broker, _ := futures.NewBroker(data.PLAT_BINANCE, "", "", "")
	broker.Test()
	databn, _ := broker.GetFundingInfo()
	fmt.Println(databn)
}
