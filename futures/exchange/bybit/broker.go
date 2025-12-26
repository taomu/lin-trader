package bybit

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	"github.com/taomu/lin-trader/futures/types"
	"github.com/taomu/lin-trader/pkg/lintypes"
	"github.com/taomu/lin-trader/pkg/util"
)

type Broker struct {
	Datas       *types.BrokerDatas
	Api         *RestApi
	ApiInfo     *lintypes.ApiInfo
	wsAccount   *util.ExcWebsocket
	wsHost      string
	wsDepth     *util.ExcWebsocket
	wsDepthLite *util.ExcWebsocket
}

func NewBroker(apiInfo *lintypes.ApiInfo) *Broker {
	datas := &types.BrokerDatas{
		SymbolInfos: make(map[string]types.SymbolInfo),
		Positions:   make([]*types.Position, 0),
	}
	api := NewRestApi()
	return &Broker{
		Api:     api,
		ApiInfo: apiInfo,
		Datas:   datas,
		wsHost:  "wss://stream.bybit.com",
	}
}

func (b *Broker) Init() error {
	return nil
}

func (b *Broker) GetDatas() *types.BrokerDatas {
	return b.Datas
}

func (b *Broker) SetWsHost(host string) {
	if host != "" {
		b.wsHost = host
	}
}

func (b *Broker) SetRestHost(host string) {
	if host != "" {
		b.Api.SetHost(host)
	}
}

func (b *Broker) GetPremium(symbol string) ([]types.Premium, error) {
	return nil, fmt.Errorf("not implemented")
}

func (b *Broker) GetFundingRate(symbol string) (*types.FundingRate, error) {
	return nil, fmt.Errorf("not implemented")
}

func (b *Broker) GetSymbolInfos() (map[string]types.SymbolInfo, error) {
	err := b.updateSymbolInfoAll()
	if err != nil {
		return nil, err
	}
	return b.Datas.SymbolInfos, nil
}

func (b *Broker) GetTickers24h() ([]types.Ticker24H, error) {
	resp, err := b.Api.Tickers24h(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	_ = resp
	return nil, nil
}

func (b *Broker) SubDepth(symbol string, onData func(updateData *types.Depth, snapData *types.Depth)) {
}

func (b *Broker) SubDepthLite(symbol string, onData func(updateData *types.Depth)) {}

func (b *Broker) UnSubDepth(symbol string) {}

func (b *Broker) UnSubDepthLite(symbol string) {}

func (b *Broker) GetPositions() ([]*types.Position, error) {
	resp, err := b.Api.GetPositions(map[string]interface{}{
		"category":   "linear",
		"settleCoin": "USDT",
	}, b.ApiInfo)
	if err != nil {
		return nil, err
	}
	// fmt.Println(resp)
	var raw struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string `json:"category"`
			List     []struct {
				Symbol        string `json:"symbol"`
				Side          string `json:"side"`
				Size          string `json:"size"`
				AvgPrice      string `json:"avgPrice"`
				UnrealisedPnl string `json:"unrealisedPnl"`
			} `json:"list"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(resp), &raw); err != nil {
		return nil, err
	}
	if raw.RetCode != 0 {
		return nil, fmt.Errorf("GetPositions err: %s", raw.RetMsg)
	}
	positions := make([]*types.Position, 0, len(raw.Result.List))
	for _, it := range raw.Result.List {
		pa := util.StrToFloat64(it.Size)
		ep := util.StrToFloat64(it.AvgPrice)
		up := util.StrToFloat64(it.UnrealisedPnl)
		side := strings.ToUpper(it.Side)
		posSide := side
		if side == "BUY" {
			posSide = types.PosSideLong
		}
		if side == "SELL" {
			posSide = types.PosSideShort
		}
		positions = append(positions, &types.Position{
			Symbol:           it.Symbol,
			PosSide:          posSide,
			PosAmt:           pa,
			EntryPrice:       ep,
			UnrealizedProfit: up,
		})
	}
	return positions, nil
}

func (b *Broker) SubAccount(onData func(updateData types.WsData)) {
	wsURL := b.wsHost
	b.wsAccount = util.NewExcWebsocket(wsURL)
	b.wsAccount.OnConnect = func() {
		fmt.Println("bybit account connect")
	}
	b.wsAccount.OnMessage = func(msg string) {
		var _ = msg
	}
	if err := b.wsAccount.Connect(); err != nil {
		fmt.Println("bybit account ws connect err:", err)
		return
	}
}

func (b *Broker) PlaceOrder(order *types.Order) error {
	return fmt.Errorf("not implemented")
}

func (b *Broker) CancelOrder(clientOrderId string, symbol string) error {
	params := map[string]interface{}{
		"orderId": clientOrderId,
		"symbol":  symbol,
	}
	_, err := b.Api.CancelOrder(params, b.ApiInfo)
	return err
}

func (b *Broker) ClearAll() {
	if b.wsAccount != nil {
		b.wsAccount.Close()
	}
	if b.wsDepth != nil {
		b.wsDepth.Close()
	}
	if b.wsDepthLite != nil {
		b.wsDepthLite.Close()
	}
}

func (b *Broker) GetLeverageBracket(symbol string) (map[string][]types.LeverageBracket, error) {
	params := map[string]interface{}{}
	cate := "linear"
	if symbol != "" {
		s := strings.ToUpper(symbol)
		if strings.HasSuffix(s, "USD") && !strings.HasSuffix(s, "USDT") && !strings.HasSuffix(s, "USDC") {
			cate = "inverse"
		}
		params["symbol"] = s
	}
	params["category"] = cate
	resp, err := b.Api.GetRiskLimit(params)
	if err != nil {
		return nil, err
	}
	var raw struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string `json:"category"`
			List     []struct {
				ID                int    `json:"id"`
				Symbol            string `json:"symbol"`
				RiskLimitValue    string `json:"riskLimitValue"`
				MaintenanceMargin string `json:"maintenanceMargin"`
				InitialMargin     string `json:"initialMargin"`
				MaxLeverage       string `json:"maxLeverage"`
			} `json:"list"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(resp), &raw); err != nil {
		return nil, err
	}
	if raw.RetCode != 0 {
		return nil, fmt.Errorf("GetLeverageBracket err: %s", raw.RetMsg)
	}
	type rec struct {
		ID     int
		Cap    float64
		MMR    float64
		MaxLev float64
	}
	group := make(map[string][]rec)
	for _, it := range raw.Result.List {
		cap := util.StrToFloat64(it.RiskLimitValue)
		mmr := util.StrToFloat64(it.MaintenanceMargin)
		maxLev := util.StrToFloat64(it.MaxLeverage)
		stdSymbol, _ := b.ToStdSymbol(it.Symbol)
		group[stdSymbol] = append(group[stdSymbol], rec{ID: it.ID, Cap: cap, MMR: mmr, MaxLev: maxLev})
	}
	ret := make(map[string][]types.LeverageBracket)
	for sym, arr := range group {
		sort.Slice(arr, func(i, j int) bool { return arr[i].Cap < arr[j].Cap })
		brs := make([]types.LeverageBracket, 0, len(arr))
		for i, r := range arr {
			floor := 0.0
			if i > 0 {
				floor = arr[i-1].Cap
			}
			brs = append(brs, types.LeverageBracket{
				Bracket:          r.ID,
				InitialLeverage:  r.MaxLev,
				NotionalCap:      r.Cap,
				NotionalFloor:    floor,
				QtyCap:           0,
				QtyFloor:         0,
				MaintMarginRatio: r.MMR,
				Cum:              0,
			})
		}
		ret[sym] = brs
	}
	return ret, nil
}

func (b *Broker) GetDualSidePosition() (string, error) {
	return "single", nil
}

func (b *Broker) ToStdSymbol(symbol string) (string, error) {
	return strings.ToUpper(symbol), nil
}

func (b *Broker) ToOriSymbol(symbol string) (string, error) {
	return symbol, nil
}

func (b *Broker) GetFundingInfo() ([]bndata.FundingInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

// 更新交易对信息（线性与反向合约）
func (b *Broker) updateSymbolInfoAll() error {
	categories := []string{"linear", "inverse"}
	for _, cate := range categories {
		resp, err := b.Api.Instruments(map[string]interface{}{
			"category": cate,
		})
		if err != nil {
			return err
		}
		fmt.Println(resp)
		symbolInfos, err := types.TransferBybitSymbolInfo(resp)
		if err != nil {
			return err
		}
		for _, it := range symbolInfos {
			std, _ := b.ToStdSymbol(it.Symbol)
			it.Symbol = std
			b.Datas.SymbolInfos[std] = it
		}
	}
	return nil
}
