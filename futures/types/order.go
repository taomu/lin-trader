package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/taomu/lin-trader/pkg/lintypes"
)

// 公共订单信息，可以用于各个交易所下单
type Order struct {
	ClientId    string  `json:"clientId"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	PosSide     string  `json:"posSide"`
	OrderType   string  `json:"orderType"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"quantity"`
	TimeInForce string  `json:"timeInForce"`
}

func ToOkxOrder(order *Order, toOkxSymbol func(string) (string, error), symbolInfo SymbolInfo) (map[string]interface{}, error) {
	okxSymbol, err := toOkxSymbol(order.Symbol)
	if err != nil {
		return nil, err
	}
	price := fmt.Sprintf("%.*f", symbolInfo.PricePrec, order.Price)
	sz := fmt.Sprintf("%.*f", symbolInfo.QtyPrec, order.Quantity/symbolInfo.CtVal)
	ordType := order.OrderType
	if order.TimeInForce == lintypes.ORDER_TIME_IN_FORCE_GTX {
		ordType = "post_only"
	}
	if order.TimeInForce == lintypes.ORDER_TIME_IN_FORCE_FOK {
		ordType = "fok"
	}
	if order.TimeInForce == lintypes.ORDER_TIME_IN_FORCE_IOC {
		ordType = "ioc"
	}
	if order.TimeInForce == lintypes.ORDER_TIME_IN_FORCE_GTC || order.TimeInForce == "" {
		ordType = strings.ToLower(ordType)
	}

	return map[string]interface{}{
		"clOrdId": order.ClientId,
		"InstId":  okxSymbol,
		"side":    strings.ToLower(order.Side),
		"posSide": strings.ToLower(order.PosSide),
		"ordType": ordType,
		"px":      price,
		"sz":      sz,
	}, nil
}

func ToBinanceOrderParams(order *Order, toBinanceSymbol func(string) (string, error), symbolInfo SymbolInfo) (map[string]interface{}, error) {
	binanceSymbol, err := toBinanceSymbol(order.Symbol)
	if err != nil {
		return nil, err
	}
	price := fmt.Sprintf("%.*f", symbolInfo.PricePrec, order.Price)
	sz := fmt.Sprintf("%.*f", symbolInfo.QtyPrec, order.Quantity)
	//判断sz转float64是否为0
	szFloat, err := strconv.ParseFloat(sz, 64)
	if err != nil {
		return nil, err
	}
	if szFloat == 0 {
		return nil, fmt.Errorf("quantity is 0")
	}
	timeInForce := order.TimeInForce
	if timeInForce == "" {
		timeInForce = lintypes.ORDER_TIME_IN_FORCE_GTC
	}
	params := map[string]interface{}{
		"newClientOrderId": order.ClientId,
		"symbol":           binanceSymbol,
		"side":             order.Side,
		"positionSide":     order.PosSide,
		"type":             order.OrderType,
		"price":            price,
		"quantity":         sz,
		"timeInForce":      timeInForce,
	}
	if order.OrderType == lintypes.ORDER_TYPE_MARKET {
		//删除timeInForce和price参数
		delete(params, "timeInForce")
		delete(params, "price")
	}
	return params, nil
}
