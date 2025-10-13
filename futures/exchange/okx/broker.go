package okx

import (
	"encoding/json"
	"fmt"

	"github.com/taomu/lin-trader/futures/data"
	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	okdata "github.com/taomu/lin-trader/futures/exchange/okx/data"
	"github.com/taomu/lin-trader/pkg/types"
)

type Broker struct {
	ApiInfo *types.ApiInfo
}

func (b *Broker) Test() {
	fmt.Println("okx test")
}
func (b *Broker) GetPremium(symbol string) ([]data.Premium, error) {
	params := map[string]interface{}{
		"instId": symbol,
	}
	if symbol == "" {
		params["instId"] = "ANY"
	}
	resp, err := NewRestApi().GetPremium(params)
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

func (b *Broker) SubDepth(symbol string, onData func(updateData *data.Depth, snapData *data.Depth)) {
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

func (b *Broker) UnSubDepth(symbol string) {
	// if b.wsDepth == nil {
	// 	return
	// }
	// msg := `{"method": "UNSUBSCRIBE","params": ["` + strings.ToLower(symbol) + `@depth@100ms"],"id": 1}`
	// b.wsDepth.Push(msg)
}

func (b *Broker) GetPositions() ([]*data.Position, error) {
	resp, err := NewRestApi().GetPositions(map[string]interface{}{
		"instType": "SWAP",
	}, b.ApiInfo)
	if err != nil {
		return nil, err
	}
	var positionsRes okdata.PositionsRes
	if err := json.Unmarshal([]byte(resp), &positionsRes); err != nil {
		return nil, err
	}
	positions := okdata.TransformPositionToPos(positionsRes)
	return positions, nil
}
