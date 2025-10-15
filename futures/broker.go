package futures

import (
	"fmt"

	"github.com/taomu/lin-trader/futures/data"
	"github.com/taomu/lin-trader/futures/exchange/binance"
	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	"github.com/taomu/lin-trader/futures/exchange/okx"
	"github.com/taomu/lin-trader/pkg/constant"
	"github.com/taomu/lin-trader/pkg/types"
)

// 交易所公告方法
type BrokerPublic interface {
	GetPremium(symbol string) ([]data.Premium, error)                                  //获取资金费率信息
	GetFundingInfo() ([]bndata.FundingInfo, error)                                     //获取资金费率限制，仅用于binance
	GetSymbolInfos() ([]data.SymbolInfo, error)                                        //获取所有合约交易对信息
	GetTickers24h() ([]data.Ticker24H, error)                                          //获取所有合约的最新价格
	SubDepth(symbol string, onData func(updateData *data.Depth, snapData *data.Depth)) //订阅深度数据
	UnSubDepth(symbol string)                                                          //取消订阅深度数据
}

// 交易所私有方法
type BrokerPrivate interface {
	GetPositions() ([]*data.Position, error) //获取持仓信息
}

// 获取变量
type BrokerVarsGetter interface {
	GetVars() *data.BrokerVars //获取所有变量
}

type Broker interface {
	BrokerPublic
	BrokerPrivate
	BrokerVarsGetter
}

func NewBroker(plat constant.PLAT, apiKey, apiSecret, apiPass string) (Broker, error) {
	apiInfo := &types.ApiInfo{
		Key:        apiKey,
		Secret:     apiSecret,
		Passphrase: apiPass,
	}
	vars := &data.BrokerVars{}
	switch plat {
	case constant.PLAT_BINANCE:
		return binance.NewBroker(apiInfo, vars), nil
	case constant.PLAT_OKX:
		return okx.NewBroker(apiInfo, vars), nil
	default:
		return nil, fmt.Errorf("unknown platform: %s", plat)
	}
}
