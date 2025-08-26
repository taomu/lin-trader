package binance

import (
	"encoding/json"
	"fmt"
	"sort"
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
	Depth   *data.Depth
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
func (b *Broker) SubDepth(symbol string, onData func(updateData *data.Depth, snapData *data.Depth)) {
	if b.wsDepth == nil {
		b.wsDepth = util.NewExcWebsocket(b.WsUrl)
	}
	b.wsDepth.OnConnect = func() {
		fmt.Println("binance depth connect")
		msg := `{"method": "SUBSCRIBE","params": ["` + strings.ToLower(symbol) + `@depth@100ms"],"id": 1}`
		b.wsDepth.Push(msg)
	}
	b.wsDepth.OnMessage = func(msg string) {
		// fmt.Println(msg)
		var wsDepthUpdateRes bndata.WsDepthUpdateRes
		err := json.Unmarshal([]byte(msg), &wsDepthUpdateRes)
		if err != nil {
			fmt.Println("binance depth update unmarshal err:", err)
			return
		}
		depthUpdate := bndata.TransferBinanceWsDepthUpdateRes(wsDepthUpdateRes)
		b.depthMerge(*depthUpdate)
		onData(depthUpdate, b.Depth)
	}
	b.wsDepth.Connect()
	go b.initDepth(symbol)
}

func (b *Broker) UnSubDepth(symbol string) {
	if b.wsDepth == nil {
		return
	}
	msg := `{"method": "UNSUBSCRIBE","params": ["` + strings.ToLower(symbol) + `@depth@100ms"],"id": 1}`
	b.wsDepth.Push(msg)
}

func (b *Broker) depthMerge(depthUpdate data.Depth) {
	if b.Depth == nil {
		return
	}
	b.Depth.Time = depthUpdate.Time
	b.Depth.Symbol = depthUpdate.Symbol
	b.Depth.Sequence = depthUpdate.Sequence
	b.Depth.LastSeq = depthUpdate.LastSeq
	//如果depthUpdate.Asks 的 price也存在 b.Depth.Asks 中，就更新 qty，否则添加
	for _, ask := range depthUpdate.Asks {
		exist := false
		for _, item := range b.Depth.Asks {
			if item.Price == ask.Price {
				item.Qty = ask.Qty
				exist = true
				break
			}
		}
		if !exist {
			b.Depth.Asks = append(b.Depth.Asks, ask)
		}
	}
	//如果depthUpdate.Bids 的 price也存在 b.Depth.Bids 中，就更新 qty，否则添加
	for _, bid := range depthUpdate.Bids {
		exist := false
		for _, item := range b.Depth.Bids {
			if item.Price == bid.Price {
				item.Qty = bid.Qty
				exist = true
				break
			}
		}
		if !exist {
			b.Depth.Bids = append(b.Depth.Bids, bid)
		}
	}
	//去除b.Depth.Bids 和 b.Depth.Aids  中 qty 为 0 的项
	for i := 0; i < len(b.Depth.Bids); i++ {
		if b.Depth.Bids[i].Qty == 0 {
			b.Depth.Bids = append(b.Depth.Bids[:i], b.Depth.Bids[i+1:]...)
			i--
		}
	}
	for i := 0; i < len(b.Depth.Asks); i++ {
		if b.Depth.Asks[i].Qty == 0 {
			b.Depth.Asks = append(b.Depth.Asks[:i], b.Depth.Asks[i+1:]...)
			i--
		}
	}
	//对b.Depth.Bids 和 b.Depth.Aids 进行排序
	sort.Slice(b.Depth.Bids, func(i, j int) bool {
		return b.Depth.Bids[i].Price > b.Depth.Bids[j].Price
	})
	sort.Slice(b.Depth.Asks, func(i, j int) bool {
		return b.Depth.Asks[i].Price < b.Depth.Asks[j].Price
	})
	//取前500项
	b.Depth.Asks = b.Depth.Asks[:500]
	b.Depth.Bids = b.Depth.Bids[:500]

	// //打印b.Depth.Bids 和 b.Depth.Aids 长度
	// fmt.Println("合并后asks长度:", len(b.Depth.Asks))
	// fmt.Println("合并后bids长度:", len(b.Depth.Bids))
	// //打印b.Depth.Bids 和 b.Depth.Aids 前3项，循环打印出数据，不是地址
	// for i := 0; i < 3; i++ {
	// 	fmt.Println("合并后asks第", i, "项:", b.Depth.Asks[i])
	// 	fmt.Println("合并后bids第", i, "项:", b.Depth.Bids[i])
	// }
	// //打印b.Depth.Bids 和 b.Depth.Aids 后3项
	// for i := 0; i < 3; i++ {
	// 	fmt.Println("合并后asks后3项第", i, "项:", b.Depth.Asks[len(b.Depth.Asks)-1-i])
	// 	fmt.Println("合并后bids后3项第", i, "项:", b.Depth.Bids[len(b.Depth.Bids)-1-i])
	// }
	// //分隔符
	// fmt.Println("-----------------")
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
	b.Depth = bndata.TransferBinanceDepthRes(depth)
	fmt.Println("初始化获取到快照asks长度:", len(b.Depth.Asks), "bids长度:", len(b.Depth.Bids))
}
