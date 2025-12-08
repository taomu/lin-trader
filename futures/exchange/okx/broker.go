package okx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
	okdata "github.com/taomu/lin-trader/futures/exchange/okx/data"
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
		wsHost:  "wss://ws.okx.com",
	}
}
func (b *Broker) Init() error {
	err := b.updateSymbolInfoAll()
	if err != nil {
		return fmt.Errorf("OK_Broker Init updateSymbolInfoAll err: %w", err)
	}
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

// 获取溢价指数
func (b *Broker) GetPremium(symbol string) ([]types.Premium, error) {
	return nil, fmt.Errorf("not implemented")
}
func (b *Broker) GetFundingRate(symbol string) (*types.FundingRate, error) {
	symbolOri, err := b.ToOriSymbol(symbol)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"instId": symbolOri,
	}
	resp, err := b.Api.GetFundingRate(params)
	if err != nil {
		return nil, fmt.Errorf("OK_Broker GetFundingRate request err: %w", err)
	}
	var fundingRateRes okdata.FundingRateResp
	err = json.Unmarshal([]byte(resp), &fundingRateRes)
	if err != nil {
		return nil, fmt.Errorf("OK_Broker GetFundingRate Unmarshal err: %w", err)
	}
	return okdata.TransferOkxFundingRate(&fundingRateRes, b.ToStdSymbol)
}
func (b *Broker) GetSymbolInfos() (map[string]types.SymbolInfo, error) {
	err := b.updateSymbolInfoAll()
	if err != nil {
		return nil, fmt.Errorf("OK_Broker GetSymbolInfos updateSymbolInfoAll err: %w", err)
	}
	return b.Datas.SymbolInfos, nil
}
func (b *Broker) updateSymbolInfoAll() error {
	// if len(b.Datas.SymbolInfos) == 0 {
	resp, err := NewRestApi().Instruments(map[string]interface{}{
		"instType": "SWAP",
	})
	if err != nil {
		return err
	}
	symbolInfos, err := types.TransferOkxSymbolInfo(resp)
	if err != nil {
		return err
	}
	for _, it := range symbolInfos {
		b.Datas.SymbolInfos[it.Symbol] = it
	}
	// }
	return nil
}
func (b *Broker) GetTickers24h() ([]types.Ticker24H, error) {
	resp, err := NewRestApi().Tickers24h(map[string]interface{}{
		"instType": "SWAP",
	})
	if err != nil {
		return nil, err
	}
	return types.TransferOkxTicker(resp)
}

func (b *Broker) SubDepth(symbol string, onData func(updateData *types.Depth, snapData *types.Depth)) {

}

func (b *Broker) SubDepthLite(symbol string, onData func(updateData *types.Depth)) {

}

func (b *Broker) UnSubDepth(symbol string) {

}

func (b *Broker) UnSubDepthLite(symbol string) {

}

func (b *Broker) GetPositions() ([]*types.Position, error) {
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

// 订阅账户信息推送，维护 Vars 中的仓位与资金
func (b *Broker) SubAccount(onData func(updateData types.WsData)) {
	// 私有WS地址
	wsURL := b.wsHost + "/ws/v5/private"
	b.wsAccount = util.NewExcWebsocket(wsURL)
	fmt.Println("订阅子账户资金变化，wsURL:", wsURL)

	b.wsAccount.OnConnect = b.onWsAccountConnect

	b.wsAccount.OnMessage = func(msg string) {
		b.onWsAccountMessage(msg, onData)
	}

	if err := b.wsAccount.Connect(); err != nil {
		fmt.Println("okx account ws connect err:", err)
		return
	}
}

func (b *Broker) onWsAccountConnect() {
	fmt.Println("连接成功,登陆账号")
	// 登录签名：timestamp + method + requestPath + body（body为空）
	sign := func(secret, timestamp, method, requestPath, body string) string {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(timestamp + method + requestPath + body))
		return base64.StdEncoding.EncodeToString(mac.Sum(nil))
	}
	ts := time.Now().Unix()
	tsStr := strconv.FormatInt(ts, 10)
	s := sign(b.ApiInfo.Secret, tsStr, "GET", "/users/self/verify", "")
	type AccountMsg struct {
		ApiKey     string `json:"apiKey"`
		Passphrase string `json:"passphrase"`
		Timestamp  string `json:"timestamp"`
		Sign       string `json:"sign"`
	}
	type LoginMsg struct {
		Op   string       `json:"op"`
		Args []AccountMsg `json:"args"`
	}

	login := LoginMsg{
		Op: "login",
		Args: []AccountMsg{
			{
				ApiKey:     b.ApiInfo.Key,
				Passphrase: b.ApiInfo.Passphrase,
				Timestamp:  tsStr,
				Sign:       s,
			},
		},
	}
	raw, _ := json.Marshal(login)
	fmt.Println("发送登录消息:", string(raw))
	err := b.wsAccount.Push(string(raw))
	if err != nil {
		fmt.Println("OK_Broker SubAccount login err:", err)
		return
	}
}
func (b *Broker) onWsAccountMessage(msg string, onData func(updateData types.WsData)) {
	// 包装解析
	var env struct {
		Event string `json:"event"`
		Arg   struct {
			Channel  string `json:"channel"`
			InstType string `json:"instType"`
		} `json:"arg"`
		Data []json.RawMessage `json:"data"`
		Code string            `json:"code"`
	}
	if err := json.Unmarshal([]byte(msg), &env); err != nil {
		fmt.Println("OK_Broker onWsAccountMessage msg Unmarshal err:", err)
		return
	}
	if env.Event == "login" && env.Code == "0" {
		// 登录成功 订阅账户与仓位
		fmt.Println("登录成功，订阅账户与仓位")
		b.sendSubAccountMsg()
		return
	}
	switch env.Arg.Channel {
	case "balance_and_position":
		// fmt.Println("收到账户消息：", string(msg))
		wsdata, err := b.onWsBalanceAndPositionMessage(msg)
		if err != nil {
			fmt.Println("OK_Broker onWsAccountMessage onWsBalanceAndPositionMessage err:", err)
			return
		}
		if wsdata != nil {
			onData(*wsdata)
			b.Datas.BalanceAll = wsdata.Balance.BalanceAll
			b.Datas.BalanceAvail = wsdata.Balance.BalanceAvail
		}
		// fmt.Printf("收到账户消息：%+v\n", wsdata)
	case "positions":
		wsdata, err := b.onWsPositionsMessage(msg)
		if err != nil {
			fmt.Println("OK_Broker onWsAccountMessage onWsPositionsMessage err:", err)
			return
		}
		if wsdata != nil {
			onData(*wsdata)
			b.Datas.Positions = wsdata.Position
		}
	default:
		fmt.Println("ignore other channels:", env.Arg.Channel, "msg:", string(msg))
	}
}

func (b *Broker) onWsBalanceAndPositionMessage(msg string) (*types.WsData, error) {
	var env okdata.WsBalanceAndPositionMsg
	if err := json.Unmarshal([]byte(msg), &env); err != nil {
		return nil, fmt.Errorf("onWsBalanceAndPositionMessage msg Unmarshal err: %w", err)
	}
	// fmt.Printf("收到账户与仓位消息：%+v\n", env)
	wsdata, err := okdata.TransformBalanceAndPositionToWsData(env)
	if err != nil {
		return nil, fmt.Errorf("onWsBalanceAndPositionMessage TransformBalanceAndPositionToWsData err: %w", err)
	}
	return wsdata, nil
}

func (b *Broker) onWsPositionsMessage(msg string) (*types.WsData, error) {
	var env okdata.WsPositionsMsg
	if err := json.Unmarshal([]byte(msg), &env); err != nil {
		return nil, fmt.Errorf("onWsPositionsMessage msg Unmarshal err: %w", err)
	}
	// fmt.Printf("收到仓位消息：%+v\n", env)
	wsdata, err := okdata.TransformWsPositionsToWsData(env, b.Datas.SymbolInfos, b.ToStdSymbol)
	if err != nil {
		return nil, fmt.Errorf("onWsPositionsMessage TransformWsPositionsToWsData err: %w", err)
	}
	return wsdata, nil
}

func (b *Broker) sendSubAccountMsg() {
	sub := struct {
		Op   string        `json:"op"`
		Args []interface{} `json:"args"`
	}{
		Op: "subscribe",
		Args: []interface{}{
			map[string]string{"channel": "balance_and_position"},
			map[string]string{"channel": "positions", "instType": "SWAP"},
		},
	}
	raw2, _ := json.Marshal(sub)
	b.wsAccount.Push(string(raw2))
}

func (b *Broker) PlaceOrder(order *types.Order) error {
	params, err := types.ToOkxOrder(order, b.ToOriSymbol, b.Datas.SymbolInfos[order.Symbol])
	if err != nil {
		return err
	}
	fmt.Println("okx place order params:", params)
	resp, err := b.Api.PlaceOrder(params, b.ApiInfo)
	if err != nil {
		return err
	}
	orderResp, err := okdata.ParseOrderResp(resp)
	if err != nil {
		return err
	}
	if orderResp.Code != "0" {
		return fmt.Errorf("okx place order error: %s", orderResp.Msg)
	}
	return nil
}

func (b *Broker) ToOriSymbol(symbol string) (string, error) {
	if len(symbol) < 4 {
		return symbol, nil
	}
	last4 := symbol[len(symbol)-4:]
	if last4 == "USDT" || last4 == "USDC" {
		return strings.ReplaceAll(symbol, last4, "-"+last4+"-SWAP"), nil
	}
	return "", fmt.Errorf("toOkxSymbol error: %s", symbol)
}

func (b *Broker) ToStdSymbol(symbol string) (string, error) {
	//根据-分割字符串
	parts := strings.Split(symbol, "-")
	if len(parts) != 3 {
		return "", fmt.Errorf("toCommonSymbol error: %s", symbol)
	}
	if parts[2] != "SWAP" {
		return "", fmt.Errorf("toCommonSymbol error: %s", symbol)
	}
	return parts[0] + parts[1], nil
}

func (b *Broker) CancelOrder(clientOrderId string, symbol string) error {
	okxSymbol, err := b.ToOriSymbol(symbol)
	if err != nil {
		return err
	}
	params := map[string]interface{}{
		"instId":  okxSymbol,
		"clOrdId": clientOrderId,
	}
	resp, err := b.Api.CancelOrder(params, b.ApiInfo)
	if err != nil {
		return err
	}
	orderResp, err := okdata.ParseOrderResp(resp)
	if err != nil {
		return err
	}
	if orderResp.Code != "0" {
		return fmt.Errorf("okx cancel order error: %s", orderResp.Msg)
	}
	return nil
}

// 清除所有连接定时器等
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
	params := map[string]interface{}{
		"instType": "SWAP",
		"tdMode":   "cross",
	}
	if symbol == "" {
		return nil, fmt.Errorf("GetLeverageBracket error symbol is empty")
	}
	instFamily := strings.ReplaceAll(symbol, "USDT", "-USDT")
	instFamily = strings.ReplaceAll(instFamily, "USDC", "-USDC")
	params["instFamily"] = instFamily
	resp, err := b.Api.GetPositionTiers(params)
	if err != nil {
		return nil, err
	}
	var apiResp struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			MaxLever   string `json:"maxLever"`
			InstFamily string `json:"instFamily"`
			Tier       string `json:"tier"`
			MinSz      string `json:"minSz"`
			MaxSz      string `json:"maxSz"`
			Imr        string `json:"imr"`
			Mmr        string `json:"mmr"`
		} `json:"data"`
	}
	if err = json.Unmarshal([]byte(resp), &apiResp); err != nil {
		return nil, err
	}
	if apiResp.Code != "0" {
		return nil, fmt.Errorf("okx position tiers error: %s", apiResp.Msg)
	}
	symbolInfoAll, err := b.GetSymbolInfos()
	if err != nil {
		return nil, err
	}
	symbolInfo, ok := symbolInfoAll[symbol]
	if !ok {
		return nil, fmt.Errorf("symbol %s not found in okx symbol infos", symbol)
	}
	brackets := make(map[string][]types.LeverageBracket)
	for _, it := range apiResp.Data {
		tier, _ := strconv.Atoi(it.Tier)
		mmr, _ := strconv.ParseFloat(it.Mmr, 64)
		lotMax, _ := strconv.ParseFloat(it.MaxSz, 64)
		lotMin, _ := strconv.ParseFloat(it.MinSz, 64)
		maxLever, _ := strconv.ParseFloat(it.MaxLever, 64)
		stdSymbol, _ := b.ToStdSymbol(it.InstFamily + "-SWAP")
		_, ok := brackets[stdSymbol]
		if !ok {
			brackets[stdSymbol] = make([]types.LeverageBracket, 0)
		}
		brackets[stdSymbol] = append(brackets[stdSymbol], types.LeverageBracket{
			Bracket:          tier,
			InitialLeverage:  maxLever,
			NotionalCap:      0,
			NotionalFloor:    0,
			QtyCap:           lotMax * symbolInfo.CtVal,
			QtyFloor:         lotMin * symbolInfo.CtVal,
			MaintMarginRatio: mmr,
			Cum:              0,
		})
	}
	return brackets, nil
}

func (b *Broker) GetDualSidePosition() (string, error) {
	return "", nil
}

func (b *Broker) GetFundingInfo() ([]bndata.FundingInfo, error) {
	return nil, fmt.Errorf("not implemented")
}
