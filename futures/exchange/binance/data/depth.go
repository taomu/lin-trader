package data

// WebSocket 深度事件（现货/合约基本一致）
type DepthUpdate struct {
	EventType  string     `json:"e"`  // "depthUpdate"
	EventTime  int64      `json:"E"`  //推送事件时间
	Symbol     string     `json:"s"`  //交易对
	Time       int64      `json:"T"`  //深度更新时间
	FirstID    int64      `json:"U"`  //当前推送第一个updateId
	FinalID    int64      `json:"u"`  //当前推送最后一个updateId
	LastPushId int64      `json:"pu"` //上一次推送的最后一个updateId
	Bids       [][]string `json:"b"`  // [[price, qty], ...]
	Asks       [][]string `json:"a"`  // [[price, qty], ...]
}

// 深度快照
type DepthSnapshot struct {
	FinalID    int64      `json:"u"`    //当前推送最后一个updateId
	LastPushId int64      `json:"pu"`   //上一次推送的最后一个updateId
	Bids       [][]string `json:"bids"` // [[price, qty], ...]
	Asks       [][]string `json:"asks"` // [[price, qty], ...]
	Time       int64      `json:"time"` //深度更新时间
}

// OrderBook：价->量；量为0表示删除
type OrderBook struct {
	Bids map[float64]float64
	Asks map[float64]float64
}

// rest 请求depht数据
type DepthRes struct {
	FinalID int64      `json:"lastUpdateId"`
	Time    int64      `json:"time"`
	Bids    [][]string `json:"bids"`
	Asks    [][]string `json:"asks"`
}
