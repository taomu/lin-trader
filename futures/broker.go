package futures

import (
	"fmt"

	"github.com/taomu/lin-trader/futures/exchange/binance"
	"github.com/taomu/lin-trader/futures/exchange/okx"
	"github.com/taomu/lin-trader/futures/types"
	"github.com/taomu/lin-trader/pkg/lintypes"
)

type Broker interface {
	//获取broker的公共参数SymbolInfos、BalanceAll、BalanceAvail、Positions
	GetDatas() *types.BrokerDatas
	//获取溢价指数
	GetPremium(symbol string) ([]types.Premium, error)
	//获取资金费率
	GetFundingInfo(symbol string) (*types.FundingRate, error)
	//获取所有交易对
	GetSymbolInfos() (map[string]types.SymbolInfo, error)
	//获取24小时内的交易对价格变化
	GetTickers24h() ([]types.Ticker24H, error)
	//订阅深度
	SubDepth(symbol string, onData func(updateData *types.Depth, snapData *types.Depth))
	//订阅深度轻量
	SubDepthLite(symbol string, onData func(updateData *types.Depth))
	//取消订阅深度
	UnSubDepth(symbol string)
	//取消订阅深度轻量
	UnSubDepthLite(symbol string)
	//设置ws主机地址
	SetWsHost(host string)
	//设置rest主机地址
	SetRestHost(host string)
	//初始化broker
	Init()
	//取消订单
	CancelOrder(clientOrderId string, symbol string) error // 取消订单
	//清除所有连接定时器等
	ClearAll()
	//获取持仓
	GetPositions() ([]*types.Position, error)
	//订阅账户
	SubAccount(onData func(updateData types.WsData))
	//下单
	PlaceOrder(order *types.Order) error
	//获取杠杆层级 bn需要api, symbol=""时 获取全部
	GetLeverageBracket(symbol string) (map[string][]types.LeverageBracket, error)
	//查询账户持仓方向
	GetDualSidePosition() (string, error)
	//转为标准symbol
	ToStdSymbol(symbol string) (string, error)
	//转为原始symbol
	ToOriSymbol(symbol string) (string, error)
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
