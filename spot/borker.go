package spot

import (
	"github.com/taomu/lin-trader/spot/data"
	"github.com/taomu/lin-trader/pkg/lintypes"
	"github.com/taomu/lin-trader/spot/exchange/okx"
)

type Broker interface {
	Test()
	GetSymbolInfos() ([]data.SymbolInfo, error)
}

func NewBroker(plat lintypes.PLAT, apikey, apisecret, passphrase string) Broker {
	apiInfo := &lintypes.ApiInfo{
		Key:        apikey,
		Secret:     apisecret,
		Passphrase: passphrase,
	}
	if plat == lintypes.PLAT_OKX {
		return okx.NewBroker(apiInfo)
	}
	return nil
}
