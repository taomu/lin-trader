package main

import (
	"fmt"

	"github.com/taomu/lin-trader/margin/exchange/binance/data"
)

func main() {
	// apiInfo := &data.ApiInfo{
	// 	Key:    "7RDKmv164I0MCIAlKxJaaBh02trRlAexzFJoANEcUldkq0vWetbJkJBslV7olKgh",
	// 	Secret: "YkZZa1Sa2TfC8sUgCKIbgTVTvkKK8dKNCfZ4oJQg5fHpcHDyTkWDJCiDqTUyR9zG",
	// }
	// broker := margin.NewBroker(data.PLAT_BINANCE, apiInfo)
	// // _, err := broker.GetAllPairs()
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	// // // fmt.Printf("%v", pairs)
	// schedules, err := broker.ListSchedule()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("schedules: %v", schedules)
	// delistSchedules, err := broker.DelistSchedule()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("delistSchedules: %v", delistSchedules)
	str := `[{"listTime":1752076800000,"crossMarginAssets":["LA"],"isolatedMarginSymbols":["LAUSDC","LAUSDT"]}]`
	schedules, err := data.TransferBinanceSchedule(str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("schedules: %v", schedules)
}
