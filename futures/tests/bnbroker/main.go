package main

import (
	"fmt"

	"github.com/taomu/lin-trader/futures"
	"github.com/taomu/lin-trader/pkg/lintypes"
)

func main() {
	key := "1bB4Q7G6CIDWXpf5ED760F7xu0YrxAfOfbzcHoTJOdsy2jtREvYh65rbbtWEJjBc"
	secret := "yoiYKr0O9GNJgoisl6nSR23Xmo9o5qwJ9K0rmD8LjSPejOcLn90pcHXysi0I0DVh"
	broker, _ := futures.NewBroker(lintypes.PLAT_BINANCE, key, secret, "")
	//测试获取深度数据
	// subDepthTest(broker)
	//测试查询持仓
	// queryPositions(broker)
	//测试查询溢价指数
	// queryPremium(broker)
	//测试查询资金费率
	// queryFundingInfo(broker)
	//测试获取杠杆层级
	// getLeverageBracket(broker)
	//获取持仓方向
	queryDualSidePosition(broker)
}
func subDepthTest(broker futures.Broker) {
	// broker.SubDepth("BTCUSDT", func(updateData *data.Depth, snapData *data.Depth) {

	// })
	broker.GetSymbolInfos()
	select {}
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
func queryPremium(broker futures.Broker) {
	premiums, err := broker.GetPremium("BTCUSDT")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, premium := range premiums {
		fmt.Println(premium)
	}
}
func queryFundingRate(broker futures.Broker) {
	fundingRate, err := broker.GetFundingRate("BTCUSDT")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(fundingRate)
}

func getLeverageBracket(broker futures.Broker) {
	leverageBrackets, err := broker.GetLeverageBracket("BTCUSDT")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, leverageBracket := range leverageBrackets["BTCUSDT"] {
		fmt.Println(leverageBracket)
	}
}

func queryDualSidePosition(broker futures.Broker) {
	dualSidePosition, err := broker.GetDualSidePosition()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dualSidePosition)
}
