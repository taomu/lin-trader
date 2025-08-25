package data

type Depth struct {
	Time     int64
	Symbol   string
	Asks     []*BookItem
	Bids     []*BookItem
	Sequence int64 //当前序列id
	LastSeq  int64 //上一个序列id
}

type BookItem struct {
	Price float64
	Qty   float64
}
