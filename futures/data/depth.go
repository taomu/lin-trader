package data

type Depth struct {
	Time     int64       `json:"ts"`
	Symbol   string      `json:"s"`
	Asks     []*BookItem `json:"a"`
	Bids     []*BookItem `json:"b"`
	Sequence int64       `json:"id"`  //当前序列id
	LastSeq  int64       `json:"lid"` //上一个序列id
}

type BookItem struct {
	Price float64 `json:"p"`
	Qty   float64 `json:"q"`
}
