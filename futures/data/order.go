package data

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
