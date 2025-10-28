package types

type BrokerDatas struct {
	SymbolInfos  map[string]SymbolInfo
	BalanceAll   float64     //总余额
	BalanceAvail float64     //可用余额
	Positions    []*Position //持仓信息
}
