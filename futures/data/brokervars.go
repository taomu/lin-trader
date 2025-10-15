package data

type BrokerVars struct {
	Positions    []*Position //持仓信息
	BalanceAvail float64     //可用余额
	BalanceAll   float64     //总余额
}
