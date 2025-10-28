package types

type Position struct {
	Symbol           string  `json:"symbol"`
	PosSide          string  `json:"posSide"`
	PosAmt           float64 `json:"posAmt"`
	EntryPrice       float64 `json:"entryPrice"`
	UnrealizedProfit float64 `json:"unrealizedProfit"`
}
