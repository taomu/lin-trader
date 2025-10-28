package types

type Depth struct {
	Time              int64       `json:"ts"`
	Symbol            string      `json:"s"`
	Asks              []*BookItem `json:"a"`
	Bids              []*BookItem `json:"b"`
	FinalUpdateId     int64       `json:"u"`  //当前序列id
	LastFinalUpdateId int64       `json:"pu"` //上一个序列id
}

type BookItem struct {
	Price float64 `json:"p"`
	Qty   float64 `json:"q"`
}
