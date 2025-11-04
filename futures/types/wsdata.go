package types

const (
	WsDataTypeTrade    = "trade"
	WsDataTypePosition = "position"
	WsDataTypeBalance  = "balance"
	WsDataTypeOrder    = "order"
	WsDataTypeUnknow   = "unknown"
)

type WsTrade struct {
	ClientId   string  `json:"clientId"`
	OrderId    string  `json:"orderId"`
	Symbol     string  `json:"symbol"`
	Side       string  `json:"side"`
	PosSide    string  `json:"posSide"`
	OrderType  string  `json:"orderType"`
	Price      float64 `json:"price"`
	Quantity   float64 `json:"quantity"`
	OrderEvent string  `json:"orderEvent"`
	Status     string  `json:"status"`
	TradeId    int64   `json:"tradeId"`
	Profit     float64 `json:"profit"`
	FeeAsset   string  `json:"feeAsset"`
	Fee        float64 `json:"fee"`
}

type WsBalance struct {
	BalanceAll   float64 `json:"balanceAll"`
	BalanceAvail float64 `json:"balanceAvailable"`
	MEvent       string  `json:"MEvent"`
	EventSymbol  string  `json:"EventSymbol"` //当 MEvent==’FUNDING_FEE’时，为对应的交易对
}

type WsOrder struct {
	ClientId   string  `json:"clientId"`
	OrderId    string  `json:"orderId"`
	Symbol     string  `json:"symbol"`
	Side       string  `json:"side"`
	PosSide    string  `json:"posSide"`
	Price      float64 `json:"price"`
	Quantity   float64 `json:"quantity"`
	OrderEvent string  `json:"orderEvent"`
	Status     string  `json:"status"`
	AvgPrice   float64 `json:"avgPrice"`
}

type WsData struct {
	DataType string
	DataTs   int64
	Order    WsOrder
	Trade    WsTrade
	Position []*Position
	Balance  WsBalance
}
