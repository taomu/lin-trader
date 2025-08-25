package okx

import (
	"fmt"

	"github.com/taomu/lin-trader/futures/data"
	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	"github.com/taomu/lin-trader/pkg/types"
)

type Broker struct {
	ApiInfo *types.ApiInfo
}

func (b *Broker) Test() {
	fmt.Println("okx test")
}
func (b *Broker) GetPremium() ([]data.Premium, error) {
	resp, err := NewRestApi().GetPremium(map[string]interface{}{
		"instId": "ANY",
	})
	if err != nil {
		return nil, err
	}
	return data.TransferOkxPremium(resp)
}
func (b *Broker) GetFundingInfo() ([]bndata.FundingInfo, error) {
	return nil, nil
}
func (b *Broker) GetSymbolInfos() ([]data.SymbolInfo, error) {
	resp, err := NewRestApi().Instruments(map[string]interface{}{
		"instType": "SWAP",
	})
	if err != nil {
		return nil, err
	}
	return data.TransferOkxSymbolInfo(resp)
}
func (b *Broker) GetTickers24h() ([]data.Ticker24H, error) {
	resp, err := NewRestApi().Tickers24h(map[string]interface{}{
		"instType": "SWAP",
	})
	if err != nil {
		return nil, err
	}
	return data.TransferOkxTicker(resp)
}

func (b *Broker) SubDepth(symbol string, onData func(depthUpdateData *data.DepthUpdate, depthSnapData *data.Depth)) {
	// if b.wsDepth == nil {
	// 	b.wsDepth = util.NewExcWebsocket(b.WsUrl)
	// }
	// b.wsDepth.OnConnect = func() {
	// 	b.wsDepth.Push("btcusdt@depth@100ms")
	// }
	// b.wsDepth.OnMessage = func(msg string) {
	// 	fmt.Println(msg)
	// }
	// b.wsDepth.Connect()
}
