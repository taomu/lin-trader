package types

import (
	"fmt"
)

// 公共订单信息，可以用于各个交易所下单
type Order struct {
	ClientId  string  `json:"clientId"`
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	PosSide   string  `json:"posSide"`
	OrderType string  `json:"orderType"`
	Price     float64 `json:"price"`
	Quantity  float64 `json:"quantity"`
}

func ToOkxOrder(order *Order, toOkxSymbol func(string) (string, error), symbolInfo SymbolInfo) (map[string]interface{}, error) {
	okxSymbol, err := toOkxSymbol(order.Symbol)
	if err != nil {
		return nil, err
	}
	price := fmt.Sprintf("%.*f", symbolInfo.PricePrec, order.Price)
	sz := fmt.Sprintf("%.*f", symbolInfo.QtyPrec, order.Quantity/symbolInfo.CtVal)
	return map[string]interface{}{
		"clOrdId": order.ClientId,
		"InstId":  okxSymbol,
		"side":    order.Side,
		"posSide": order.PosSide,
		"ordType": order.OrderType,
		"px":      price,
		"sz":      sz,
	}, nil
}
