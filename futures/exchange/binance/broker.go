package binance

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	"github.com/taomu/lin-trader/futures/types"
	"github.com/taomu/lin-trader/pkg/lintypes"
	"github.com/taomu/lin-trader/pkg/util"
)

type Broker struct {
	Datas     *types.BrokerDatas
	Api       *RestApi
	ApiInfo   *lintypes.ApiInfo
	wsAccount *util.ExcWebsocket
	wsUrl     string
	wsDepth   *util.ExcWebsocket
	Depth     *types.Depth
}

func NewBroker(apiInfo *lintypes.ApiInfo) *Broker {
	datas := &types.BrokerDatas{
		SymbolInfos: make(map[string]types.SymbolInfo),
		Positions:   make([]*types.Position, 0),
	}
	return &Broker{
		ApiInfo: apiInfo,
		wsUrl:   "wss://fstream.binance.com/ws",
		Api:     NewRestApi(),
		Datas:   datas,
	}
}

func (b *Broker) GetDatas() *types.BrokerDatas {
	return b.Datas
}

func (b *Broker) SetWsHost(host string) {
	if host != "" {
		b.wsUrl = host + "/ws"
	}
}

func (b *Broker) SetRestHost(host string) {
	if host != "" {
		b.Api.SetHost(host)
	}
}

func (b *Broker) GetPremium(symbol string) ([]types.Premium, error) {
	params := map[string]interface{}{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	resp, err := b.Api.PremiumIndex(params)
	if err != nil {
		return nil, err
	}
	return types.TransferBinancePremium(resp)
}
func (b *Broker) GetFundingInfo() ([]bndata.FundingInfo, error) {
	params := map[string]interface{}{}
	resp, err := b.Api.FundingInfo(params)
	if err != nil {
		return nil, err
	}
	return bndata.TransferBinanceFundingInfo(resp)
}
func (b *Broker) GetSymbolInfos() ([]types.SymbolInfo, error) {
	resp, err := b.Api.ExchangeInfo(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return types.TransferBinanceSymbolInfo(resp)
}
func (b *Broker) GetTickers24h() ([]types.Ticker24H, error) {
	return nil, nil
}
func (b *Broker) SubDepth(symbol string, onData func(updateData *types.Depth, snapData *types.Depth)) {
	if b.wsDepth == nil {
		b.wsDepth = util.NewExcWebsocket(b.wsUrl)
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
		dMap := make(map[string]float64)
		for _, ask := range depthUpdate.Asks {
			priceStr := fmt.Sprintf("%f", ask.Price)

			//判断 priceStr是否已存在dMap
			_, ok := dMap[priceStr]
			if ok {
				fmt.Println("priceStr已存在", priceStr)
			}
			dMap[priceStr] = ask.Qty

		}
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

func (b *Broker) depthMerge(depthUpdate types.Depth) {
	if b.Depth == nil {
		return
	}
	// 添加对 Bids 和 Asks 的空指针检查
	if b.Depth.Bids == nil || b.Depth.Asks == nil {
		return
	}
	b.Depth.Time = depthUpdate.Time
	b.Depth.Symbol = depthUpdate.Symbol
	b.Depth.FinalUpdateId = depthUpdate.FinalUpdateId
	b.Depth.LastFinalUpdateId = depthUpdate.LastFinalUpdateId
	//先把depthUpdate.Asks 转map[string]float64
	askMap := make(map[string]float64)
	for _, ask := range b.Depth.Asks {
		if ask == nil {
			continue
		}
		askMap[fmt.Sprintf("%f", ask.Price)] = ask.Qty
	}
	//更新数据
	for _, ask := range depthUpdate.Asks {
		askMap[fmt.Sprintf("%f", ask.Price)] = ask.Qty
	}
	//把askMap 转 b.Depth.Asks
	b.Depth.Asks = make([]*types.BookItem, 0)
	for price, qty := range askMap {
		if qty == 0 {
			continue
		}
		b.Depth.Asks = append(b.Depth.Asks, &types.BookItem{
			Price: util.StrToFloat64(price),
			Qty:   qty,
		})
	}
	//先把depthUpdate.Bids 转map[string]float64
	bidMap := make(map[string]float64)
	for _, bid := range b.Depth.Bids {
		if bid == nil {
			continue
		}
		bidMap[fmt.Sprintf("%f", bid.Price)] = bid.Qty
	}
	//更新数据
	for _, bid := range depthUpdate.Bids {
		bidMap[fmt.Sprintf("%f", bid.Price)] = bid.Qty
	}
	//把bidMap 转 b.Depth.Bids
	b.Depth.Bids = make([]*types.BookItem, 0)
	for price, qty := range bidMap {
		if qty == 0 {
			continue
		}
		b.Depth.Bids = append(b.Depth.Bids, &types.BookItem{
			Price: util.StrToFloat64(price),
			Qty:   qty,
		})
	}

	//对b.Depth.Bids 和 b.Depth.Aids 进行排序
	sort.Slice(b.Depth.Bids, func(i, j int) bool {
		return b.Depth.Bids[i].Price > b.Depth.Bids[j].Price
	})
	sort.Slice(b.Depth.Asks, func(i, j int) bool {
		return b.Depth.Asks[i].Price < b.Depth.Asks[j].Price
	})
	//取前500项
	b.Depth.Asks = b.Depth.Asks[:int(math.Min(float64(len(b.Depth.Asks)), 500))]
	b.Depth.Bids = b.Depth.Bids[:int(math.Min(float64(len(b.Depth.Bids)), 500))]

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

func (b *Broker) GetPositions() ([]*types.Position, error) {
	params := map[string]interface{}{}
	resp, err := b.Api.Account(params, b.ApiInfo)
	if err != nil {
		fmt.Println("binance account err:", err)
		return nil, err
	}
	var account bndata.AccountRes
	err = json.Unmarshal([]byte(resp), &account)
	if err != nil {
		fmt.Println("binance account unmarshal err:", err)
		return nil, err
	}
	return bndata.TransferBinanceAccountResToPos(account), err
}

// 订阅账户信息推送，维护 Vars 中的仓位与资金
func (b *Broker) SubAccount() {
	// 1) 获取 listenKey
	lkResp, err := b.Api.StartUserDataStream(b.ApiInfo.Key)
	if err != nil {
		fmt.Println("start user data stream err:", err)
		return
	}
	var lk struct {
		ListenKey string `json:"listenKey"`
	}
	if err := json.Unmarshal([]byte(lkResp), &lk); err != nil || lk.ListenKey == "" {
		fmt.Println("parse listenKey err:", err, "resp:", lkResp)
		return
	}

	// 2) 建立账户 WebSocket
	wsURL := "wss://fstream.binance.com/ws/" + lk.ListenKey
	b.wsAccount = util.NewExcWebsocket(wsURL)
	b.wsAccount.OnConnect = func() {
		fmt.Println("binance account connect")
	}
	b.wsAccount.OnMessage = func(msg string) {
		// 检测事件类型
		var header struct {
			EventType string `json:"e"`
		}
		if err := json.Unmarshal([]byte(msg), &header); err != nil {
			return
		}
		if header.EventType != "ACCOUNT_UPDATE" {
			return
		}

		// 解析 ACCOUNT_UPDATE
		var accUpdate struct {
			EventType string `json:"e"`
			EventTime int64  `json:"E"`
			Acc       struct {
				Balances []struct {
					Asset string `json:"a"`
					Wb    string `json:"wb"` // 钱包余额
					Cw    string `json:"cw"` // Cross Wallet（用作可用余额近似）
				} `json:"B"`
				Positions []struct {
					Symbol string `json:"s"`
					Pa     string `json:"pa"` // 持仓数量
					Ep     string `json:"ep"` // 开仓均价
					Up     string `json:"up"` // 未实现盈亏
					Ps     string `json:"ps"` // 持仓方向：BOTH/LONG/SHORT
				} `json:"P"`
			} `json:"a"`
		}
		if err := json.Unmarshal([]byte(msg), &accUpdate); err != nil {
			fmt.Println("account update unmarshal err:", err)
			return
		}

		// 更新资金（USDT）
		for _, binfo := range accUpdate.Acc.Balances {
			if binfo.Asset == "USDT" {
				if v, err := strconv.ParseFloat(binfo.Wb, 64); err == nil {
					b.Datas.BalanceAll = v
				}
				if v, err := strconv.ParseFloat(binfo.Cw, 64); err == nil {
					b.Datas.BalanceAvail = v
				}
				break
			}
		}

		// 更新仓位
		positions := make([]*types.Position, 0, len(accUpdate.Acc.Positions))
		for _, p := range accUpdate.Acc.Positions {
			pa, _ := strconv.ParseFloat(p.Pa, 64)
			if pa == 0 {
				continue
			}
			ep, _ := strconv.ParseFloat(p.Ep, 64)
			up, _ := strconv.ParseFloat(p.Up, 64)
			positions = append(positions, &types.Position{
				Symbol:           p.Symbol,
				PosSide:          p.Ps,
				PosAmt:           pa,
				EntryPrice:       ep,
				UnrealizedProfit: up,
			})
		}
		b.Datas.Positions = positions
	}
	if err := b.wsAccount.Connect(); err != nil {
		fmt.Println("binance account ws connect err:", err)
		return
	}

	// 3) 保活 listenKey（每 30 分钟）
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := b.Api.KeepaliveUserDataStream(b.ApiInfo.Key, lk.ListenKey); err != nil {
				fmt.Println("keepalive listenKey err:", err)
			}
		}
	}()
}

func (b *Broker) PlaceOrder(order *types.Order) error {
	return nil
}
