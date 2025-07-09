package margin

import (
	"github.com/taomu/lin-trader/margin/data"
	"github.com/taomu/lin-trader/margin/exchange/binance"
	bndata "github.com/taomu/lin-trader/margin/exchange/binance/data"
)

type Broker interface {
	Test()
	GetAllPairs() ([]bndata.Pair, error)
	ListSchedule() ([]bndata.Schedule, error)
	DelistSchedule() ([]bndata.Schedule, error)
}

func NewBroker(plat data.PLAT, apikey, apisecret, passphrase string) Broker {
	apiInfo := &data.ApiInfo{
		Key:        apikey,
		Secret:     apisecret,
		Passphrase: passphrase,
	}
	if plat == data.PLAT_BINANCE {
		return binance.NewBroker(apiInfo)
	}
	return nil
}
