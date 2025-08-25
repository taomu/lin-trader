package data

type Depth struct {
	Time     int64
	Symbol   string
	Asks     []OrderBookItem
	Bids     []OrderBookItem
	DataId   int64
	Sequence int64
}

type DepthUpdate struct {
	Time     int64
	Symbol   string
	Asks     []OrderBookItem
	Bids     []OrderBookItem
	DataId   int64
	Sequence int64
}

type OrderBookItem struct {
	Price string
	Qty   string
}
