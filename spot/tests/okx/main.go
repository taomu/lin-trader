package main

import (
	"fmt"

	"github.com/taomu/lin-trader/common"
	"github.com/taomu/lin-trader/spot"
)

func main() {
	broker := spot.NewBroker(common.PLAT_OKX, "", "", "")
	broker.Test()
	res, err := broker.GetSymbolInfos()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v", res)
}
