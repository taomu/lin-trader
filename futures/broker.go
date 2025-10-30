package futures

import (
	"fmt"

	"github.com/taomu/lin-trader/futures/exchange/binance"
	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	"github.com/taomu/lin-trader/futures/exchange/okx"
	"github.com/taomu/lin-trader/futures/types"
	"github.com/taomu/lin-trader/pkg/lintypes"
)

// 交易所公告方法
type BrokerPublic interface {
	GetDatas() *types.BrokerDatas
	GetPremium(symbol string) ([]types.Premium, error)
	GetFundingInfo() ([]bndata.FundingInfo, error)
	GetSymbolInfos() ([]types.SymbolInfo, error)
	GetTickers24h() ([]types.Ticker24H, error)
	SubDepth(symbol string, onData func(updateData *types.Depth, snapData *types.Depth))
	UnSubDepth(symbol string)
	SetWsHost(host string)
	SetRestHost(host string)
}

// 交易所私有方法
type BrokerPrivate interface {
	GetPositions() ([]*types.Position, error)
	SubAccount(onData func(updateData types.WsData))
	PlaceOrder(order *types.Order) error // 下单
}

type Broker interface {
	BrokerPublic
	BrokerPrivate
}

func NewBroker(plat lintypes.PLAT, apiKey, apiSecret, apiPass string) (Broker, error) {
	apiInfo := &lintypes.ApiInfo{
		Key:        apiKey,
		Secret:     apiSecret,
		Passphrase: apiPass,
	}
	switch plat {
	case lintypes.PLAT_BINANCE:
		return binance.NewBroker(apiInfo), nil
	case lintypes.PLAT_OKX:
		return okx.NewBroker(apiInfo), nil
	default:
		return nil, fmt.Errorf("unknown platform: %s", plat)
	}
}
