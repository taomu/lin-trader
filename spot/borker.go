package spot

import (
	"github.com/taomu/lin-trader/common"
	"github.com/taomu/lin-trader/spot/data"
	"github.com/taomu/lin-trader/spot/exchange/okx"
)

type Broker interface {
	Test()
	GetSymbolInfos() ([]data.SymbolInfo, error)
}

func NewBroker(plat common.PLAT, apikey, apisecret, passphrase string) Broker {
	apiInfo := &common.ApiInfo{
		Key:        apikey,
		Secret:     apisecret,
		Passphrase: passphrase,
	}
	if plat == common.PLAT_OKX {
		return okx.NewBroker(apiInfo)
	}
	return nil
}
