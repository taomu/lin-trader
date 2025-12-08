package main

import (
	"fmt"

	"github.com/taomu/lin-trader/futures"
	"github.com/taomu/lin-trader/pkg/lintypes"
)

func main() {
	key := "8cl9PSIlqVImMTwfSV"
	secret := "oNZbGn3OBsrzDhdumeRyTESHSokxzZ37L7VZ"
	broker, _ := futures.NewBroker(lintypes.PLAT_BYBIT, key, secret, "")
	positions, err := broker.GetPositions()
	if err != nil {
		fmt.Println("GetPositions err:", err)
		return
	}
	fmt.Println("持仓数量:", len(positions))
	for _, p := range positions {
		fmt.Printf("%+v\n", *p)
	}
}
