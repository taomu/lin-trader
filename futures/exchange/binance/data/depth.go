package data

// WebSocket 深度事件（现货/合约基本一致）
type DepthUpdate struct {
	EventType string     `json:"e"` // "depthUpdate"
	EventTime int64      `json:"E"`
	Symbol    string     `json:"s"`
	FirstID   int64      `json:"U"`
	FinalID   int64      `json:"u"`
	Bids      [][]string `json:"b"` // [[price, qty], ...]
	Asks      [][]string `json:"a"`
}

// REST 快照
type DepthSnapshot struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// OrderBook：价->量；量为0表示删除
type OrderBook struct {
	Bids map[float64]float64
	Asks map[float64]float64
}
