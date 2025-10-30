package types

const (
	WsDataTypeTrade    = "trade"
	WsDataTypePosition = "position"
	WsDataTypeBalance  = "balance"
	WsDataTypeOrder    = "order"
	WsDataTypeUnknow   = "unknown"
)

type WsTrade struct {
	ClientId    string  `json:"clientId"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	PosSide     string  `json:"posSide"`
	OrderType   string  `json:"orderType"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"quantity"`
	TimeInForce string  `json:"timeInForce"`
}

type WsBalance struct {
	BalanceAll   float64 `json:"balanceAll"`
	BalanceAvail float64 `json:"balanceAvailable"`
}

type WsOrder struct {
	ClientId    string  `json:"clientId"`
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	PosSide     string  `json:"posSide"`
	OrderType   string  `json:"orderType"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"quantity"`
	TimeInForce string  `json:"timeInForce"`
}

type WsData struct {
	DataType string
	Order    WsOrder
	Trade    WsTrade
	Position []*Position
	Balance  WsBalance
}
