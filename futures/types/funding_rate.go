package types

type FundingRate struct {
	Symbol       string
	SymbolOri    string
	Rate         float64
	NextSettleTs int64
	MaxRate      float64
	MinRate      float64
	SettlePeriod int64
}
