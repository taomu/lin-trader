package main

import (
	"fmt"

	"github.com/taomu/lin-trader/futures"
	"github.com/taomu/lin-trader/pkg/lintypes"
)

func main() {
	broker, err := futures.NewBroker(lintypes.PLAT_BYBIT, "", "", "")
	if err != nil {
		fmt.Println("NewBroker err:", err)
		return
	}
	_, err = broker.GetSymbolInfos()
	if err != nil {
		// fmt.Println("GetSymbolInfos err:", err)
		return
	}
	fmt.Println("GetSymbolInfos done")
}
