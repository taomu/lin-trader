package margin

import (
	"github.com/taomu/lin-trader/margin/data"
	"github.com/taomu/lin-trader/margin/exchange/binance"
)

type Broker interface {
	Test()
}

func NewBroker(plat data.PLAT) Broker {
	if plat == data.PLAT_BINANCE {
		return &binance.Broker{}
	}
	return nil
}
