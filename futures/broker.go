package futures

import (
	"fmt"

	"github.com/taomu/lin-trader/futures/data"
	"github.com/taomu/lin-trader/futures/exchange/binance"
	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	"github.com/taomu/lin-trader/futures/exchange/okx"
)

type BrokerPublic interface {
	GetPremium() ([]data.Premium, error)           //获取资金费率信息
	GetFundingInfo() ([]bndata.FundingInfo, error) //获取资金费率限制，仅用于binance
	GetSymbolInfos() ([]data.SymbolInfo, error)    //获取所有合约交易对信息
	GetTickers24h() ([]data.Ticker24H, error)           //获取所有合约的最新价格
}

type BrokerPrivate interface {
}

type Broker interface {
	BrokerPublic
	BrokerPrivate
	Test()
}

func NewBroker(plat data.PLAT, apiKey, apiSecret, apiPass string) (Broker, error) {
	apiInfo := &data.ApiInfo{
		Key:        apiKey,
		Secret:     apiSecret,
		Passphrase: apiPass,
	}
	switch plat {
	case data.PLAT_BINANCE:
		return &binance.Broker{
			ApiInfo: apiInfo,
		}, nil
	case data.PLAT_OKX:
		return &okx.Broker{
			ApiInfo: apiInfo,
		}, nil
	default:
		return nil, fmt.Errorf("unknown platform: %s", plat)
	}
}
