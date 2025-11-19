package main

import (
	"fmt"

	"github.com/taomu/lin-trader/futures"
	"github.com/taomu/lin-trader/pkg/lintypes"
)

func main() {
	key := "f99fbeeb-a81d-4be3-8b4b-1eeb59370aee"
	secret := "F593239483463966769FBBE46743DBAA"
	passphrase := "@Carp008"
	broker, _ := futures.NewBroker(lintypes.PLAT_OKX, key, secret, passphrase)
	// showOkExample(broker)
	//测试查询持仓
	// queryPositions(broker)
	//测试查询杠杆层级
	// queryLeverageBrackets(broker)
	//测试查询资金费率
	queryFundingRate(broker)
}
func showOkExample(broker futures.Broker) {
	// premium, _ := broker.GetPremium()
	// fmt.Println(premium)
	tickers, _ := broker.GetTickers24h()
	fmt.Println(tickers)
}
func queryPositions(broker futures.Broker) {
	positions, err := broker.GetPositions()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, pos := range positions {
		fmt.Println(pos)
	}
}

func queryLeverageBrackets(broker futures.Broker) {
	brackets, err := broker.GetLeverageBracket("BTCUSDT")
	fmt.Println("查询到杠杆层级数量：", len(brackets))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("查询到杠杆层级：%+v\n", brackets)
}

// 查询资金费率
func queryFundingRate(broker futures.Broker) {
	fundingRate, err := broker.GetFundingRate("BTCUSDT")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fundingRate)
}
