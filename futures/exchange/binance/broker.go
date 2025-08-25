package binance

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/taomu/lin-trader/futures/data"
	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	"github.com/taomu/lin-trader/pkg/types"
	"github.com/taomu/lin-trader/pkg/util"
)

type Broker struct {
	ApiInfo *types.ApiInfo
	WsUrl   string
	wsDepth *util.ExcWebsocket
	Depth   *bndata.DepthSnapshot
	Api     *RestApi
}

func NewBroker(apiInfo *types.ApiInfo) *Broker {
	return &Broker{
		ApiInfo: apiInfo,
		WsUrl:   "wss://fstream.binance.com/ws",
		Api:     NewRestApi(),
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
func (b *Broker) SubDepth(symbol string, onData func(depthUpdateData *data.DepthUpdate, depthSnapData *data.Depth)) {
	fmt.Println("binance sub depth" + b.WsUrl)
	if b.wsDepth == nil {
		b.wsDepth = util.NewExcWebsocket(b.WsUrl)
	}
	fmt.Println("binance sub depth2")
	b.wsDepth.OnConnect = func() {
		fmt.Println("binance depth connect")
		msg := `{"method": "SUBSCRIBE","params": ["` + strings.ToLower(symbol) + `@depth@100ms"],"id": 1}`
		b.wsDepth.Push(msg)
	}
	fmt.Println("binance sub depth3")
	b.wsDepth.OnMessage = func(msg string) {
		// fmt.Println(msg)
		var depthUpdate bndata.DepthUpdate
		err := json.Unmarshal([]byte(msg), &depthUpdate)
		if err != nil {
			fmt.Println("binance depth update unmarshal err:", err)
			return
		}
		fmt.Println("binance depth update:", depthUpdate)
	}
	// b.wsDepth.Connect()
	go b.initDepth(symbol)
}

func (b *Broker) depthMerge(bndata.DepthUpdate) {
	if b.Depth == nil {
		b.Depth = &bndata.DepthSnapshot{}
	}
}

func (b *Broker) initDepth(symbol string) {
	params := map[string]interface{}{
		"symbol": symbol,
	}
	resp, err := b.Api.Depth(params)
	if err != nil {
		fmt.Println("binance depth err:", err)
		return
	}
	var depth bndata.DepthRes
	err = json.Unmarshal([]byte(resp), &depth)
	if err != nil {
		fmt.Println("binance depth unmarshal err:", err)
		return
	}
	fmt.Printf("获取到快照：%+v", depth)
}
