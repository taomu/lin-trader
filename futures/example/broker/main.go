package main

import (
	"fmt"

	traderfu "github.com/taomu/lin-trader/futures"
	"github.com/taomu/lin-trader/futures/data"
)

func main() {
	brokerOk, _ := traderfu.NewBroker(data.PLAT_OKX, "", "", "")
	brokerOk.Test()
	dataok, _ := brokerOk.GetPremium()
	fmt.Println(dataok)
	brokerBn, _ := traderfu.NewBroker(data.PLAT_BINANCE, "", "", "")
	brokerBn.Test()
	// databn, _ := brokerBn.GetFundingInfo()
	// fmt.Println(databn)
}
