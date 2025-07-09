package binance

import (
	"fmt"

	"github.com/taomu/lin-trader/margin/data"
	bndata "github.com/taomu/lin-trader/margin/exchange/binance/data"
)

type Broker struct {
	ApiInfo *data.ApiInfo
}

func (b *Broker) Test() {
	fmt.Println("binance test")
}
func (b *Broker) GetAllPairs() ([]bndata.Pair, error) {
	resp, err := NewRestApi().GetAllPairs(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return bndata.TransferBinancePair(resp)
}
