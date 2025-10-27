package main

import (
	"fmt"

	"github.com/taomu/lin-trader/pkg/lintypes"
	"github.com/taomu/lin-trader/spot"
)

func main() {
	broker := spot.NewBroker(lintypes.PLAT_OKX, "", "", "")
	broker.Test()
	res, err := broker.GetSymbolInfos()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v", res)
}
