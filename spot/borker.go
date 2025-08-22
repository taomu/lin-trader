package spot

import (
	"github.com/taomu/lin-trader/spot/data"
	"github.com/taomu/lin-trader/pkg/constant"
	"github.com/taomu/lin-trader/pkg/types"
	"github.com/taomu/lin-trader/spot/exchange/okx"
)

type Broker interface {
	Test()
	GetSymbolInfos() ([]data.SymbolInfo, error)
}

func NewBroker(plat constant.PLAT, apikey, apisecret, passphrase string) Broker {
	apiInfo := &types.ApiInfo{
		Key:        apikey,
		Secret:     apisecret,
		Passphrase: passphrase,
	}
	if plat == constant.PLAT_OKX {
		return okx.NewBroker(apiInfo)
	}
	return nil
}
