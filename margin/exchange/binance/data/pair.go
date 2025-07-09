package data

import "encoding/json"

type Pair struct {
	Symbol        string `json:"symbol"`
	Base          string `json:"base"`
	Quote         string `json:"quote"`
	IsBuyAllowed  bool   `json:"isBuyAllowed"`
	IsSellAllowed bool   `json:"isSellAllowed"`
	IsMarginTrade bool   `json:"isMarginTrade"`
	DeListTime    int64  `json:"delistTime"`
}

func TransferBinancePair(resp string) ([]Pair, error) {
	var pairs []Pair
	if err := json.Unmarshal([]byte(resp), &pairs); err != nil {
		return nil, err
	}
	return pairs, nil
}
