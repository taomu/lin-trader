package binance

import (
	"fmt"

	"github.com/taomu/lin-trader/futures/data"
	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	"github.com/taomu/lin-trader/pkg/util"
)

type Broker struct {
	ApiInfo *data.ApiInfo
	WsUrl   string
	wsDepth *util.ExcWebsocket
}

func NewBroker(apiInfo *data.ApiInfo) *Broker {
	return &Broker{
		ApiInfo: apiInfo,
		WsUrl:   "wss://fstream.binance.com/ws",
	}
}

func (b *Broker) Test() {
	fmt.Println("binance test")
}
func (b *Broker) GetPremium() ([]data.Premium, error) {
	resp, err := NewRestApi().PremiumIndex(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return data.TransferBinancePremium(resp)
}
func (b *Broker) GetFundingInfo() ([]bndata.FundingInfo, error) {
	params := map[string]interface{}{}
	resp, err := NewRestApi().FundingInfo(params)
	if err != nil {
		return nil, err
	}
	return bndata.TransferBinanceFundingInfo(resp)
}
func (b *Broker) GetSymbolInfos() ([]data.SymbolInfo, error) {
	resp, err := NewRestApi().ExchangeInfo(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return data.TransferBinanceSymbolInfo(resp)
}
func (b *Broker) GetTickers24h() ([]data.Ticker24H, error) {
	return nil, nil
}
func (b *Broker) SubDepth(onData func(depthUpdateData *bndata.DepthUpdate, depthSnapData *bndata.DepthSnapshot)) {
	if b.wsDepth == nil {
		b.wsDepth = util.NewExcWebsocket(b.WsUrl)
	}
	b.wsDepth.OnConnect = func() {
		b.wsDepth.Push("btcusdt@depth@100ms")
	}
	b.wsDepth.OnMessage = func(msg string) {
		fmt.Println(msg)
	}
	b.wsDepth.Connect()
}
