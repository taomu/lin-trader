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
	Datas     *types.BrokerDatas
	Api       *RestApi
	ApiInfo   *lintypes.ApiInfo
	wsAccount *util.ExcWebsocket
	wsHost    string
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
	params := map[string]interface{}{
		"instId": symbol,
	}
	if symbol == "" {
		params["instId"] = "ANY"
	}
	resp, err := b.Api.GetPremium(params)
	if err != nil {
		return nil, err
	}
	return types.TransferOkxPremium(resp)
}
func (b *Broker) GetFundingInfo() ([]bndata.FundingInfo, error) {
	return nil, nil
}
func (b *Broker) GetSymbolInfos() ([]types.SymbolInfo, error) {
	resp, err := NewRestApi().Instruments(map[string]interface{}{
		"instType": "SWAP",
	})
	if err != nil {
		return nil, err
	}
	return types.TransferOkxSymbolInfo(resp)
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
func (b *Broker) SubAccount() {
	// 私有WS地址
	wsURL := "wss://ws.okx.com:8443/ws/v5/private"
	b.wsAccount = util.NewExcWebsocket(wsURL)

	// 登录签名：timestamp + method + requestPath + body（body为空）
	sign := func(secret, timestamp, method, requestPath, body string) string {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(timestamp + method + requestPath + body))
		return base64.StdEncoding.EncodeToString(mac.Sum(nil))
	}

	b.wsAccount.OnConnect = func() {
		// 登录
		ts := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		s := sign(b.ApiInfo.Secret, ts, "GET", "/users/self/verify", "")
		login := struct {
			Op   string `json:"op"`
			Args []struct {
				ApiKey     string `json:"apiKey"`
				Passphrase string `json:"passphrase"`
				Timestamp  string `json:"timestamp"`
				Sign       string `json:"sign"`
			} `json:"args"`
		}{
			Op: "login",
			Args: []struct {
				ApiKey     string `json:"apiKey"`
				Passphrase string `json:"passphrase"`
				Timestamp  string `json:"timestamp"`
				Sign       string `json:"sign"`
			}{
				{
					ApiKey:     b.ApiInfo.Key,
					Passphrase: b.ApiInfo.Passphrase,
					Timestamp:  ts,
					Sign:       s,
				},
			},
		}
		raw, _ := json.Marshal(login)
		_ = b.wsAccount.Push(string(raw))

		// 订阅账户与仓位
		sub := struct {
			Op   string        `json:"op"`
			Args []interface{} `json:"args"`
		}{
			Op: "subscribe",
			Args: []interface{}{
				map[string]string{"channel": "account"},
				map[string]string{"channel": "positions", "instType": "SWAP"},
			},
		}
		raw2, _ := json.Marshal(sub)
		_ = b.wsAccount.Push(string(raw2))
	}

	b.wsAccount.OnMessage = func(msg string) {
		// 包装解析
		var env struct {
			Event string `json:"event"`
			Arg   struct {
				Channel  string `json:"channel"`
				InstType string `json:"instType"`
			} `json:"arg"`
			Data []json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal([]byte(msg), &env); err != nil {
			return
		}
		// 忽略非数据事件
		if env.Event != "" {
			return
		}

		switch env.Arg.Channel {
		case "account":
			// 账户资金
			for _, d := range env.Data {
				var ad struct {
					TotalEq string `json:"totalEq"`
					AvailEq string `json:"availEq"`
					Details []struct {
						Ccy      string `json:"ccy"`
						CashBal  string `json:"cashBal"`
						AvailBal string `json:"availBal"`
					} `json:"details"`
				}
				if err := json.Unmarshal(d, &ad); err != nil {
					continue
				}
				// 优先使用 availEq/totalEq
				if v, err := strconv.ParseFloat(ad.TotalEq, 64); err == nil {
					b.Datas.BalanceAll = v
				}
				if v, err := strconv.ParseFloat(ad.AvailEq, 64); err == nil {
					b.Datas.BalanceAvail = v
				}
				// 如果有USDT详情，进一步精确可用余额
				for _, det := range ad.Details {
					if det.Ccy == "USDT" {
						if v, err := strconv.ParseFloat(det.AvailBal, 64); err == nil {
							b.Datas.BalanceAvail = v
						}
						break
					}
				}
			}
		case "positions":
			// 仓位
			positions := make([]*types.Position, 0, len(env.Data))
			for _, d := range env.Data {
				var p struct {
					InstId  string `json:"instId"`
					Pos     string `json:"pos"`
					PosSide string `json:"posSide"`
					AvgPx   string `json:"avgPx"`
				}
				if err := json.Unmarshal(d, &p); err != nil {
					continue
				}
				posAmt, err := strconv.ParseFloat(p.Pos, 64)
				if err != nil || posAmt == 0 {
					continue
				}
				parts := strings.Split(strings.Replace(p.InstId, "-SWAP", "", -1), "-")
				symbol := strings.Join(parts, "")
				side := ""
				if p.PosSide == "long" {
					side = "LONG"
				} else if p.PosSide == "short" {
					side = "SHORT"
				}
				entryPrice, _ := strconv.ParseFloat(p.AvgPx, 64)
				positions = append(positions, &types.Position{
					Symbol:     symbol,
					PosAmt:     posAmt,
					PosSide:    side,
					EntryPrice: entryPrice,
				})
			}
			b.Datas.Positions = positions
		default:
			// ignore other channels
		}
	}

	if err := b.wsAccount.Connect(); err != nil {
		fmt.Println("okx account ws connect err:", err)
		return
	}
}

func (b *Broker) PlaceOrder(order *types.Order) error {
	params, err := types.ToOkxOrder(order, b.toOkxSymbol, b.Datas.SymbolInfos[order.Symbol])
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

func (b *Broker) toOkxSymbol(symbol string) (string, error) {
	if len(symbol) < 4 {
		return symbol, nil
	}
	last4 := symbol[len(symbol)-4:]
	if last4 == "USDT" || last4 == "USDC" {
		return strings.ReplaceAll(symbol, last4, last4+"-SWAP"), nil
	}
	return "", fmt.Errorf("toOkxSymbol error: %s", symbol)
}

func (b *Broker) toCommonSymbol(symbol string) (string, error) {
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
