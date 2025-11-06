package types

/*
	"bracket": 1,   // 层级
	"initialLeverage": 75,  // 该层允许的最高初始杠杆倍数
	"notionalCap": 10000,  // 该层对应的名义价值上限
	"notionalFloor": 0,  // 该层对应的名义价值下限
	"maintMarginRatio": 0.0065, // 该层对应的维持保证金率
	"cum": 0.0 // 速算数
*/

type LeverageBracket struct {
	Bracket          int     `json:"bracket"`          // 杠杆层级
	InitialLeverage  float64 `json:"initialLeverage"`  // 该层允许的最高初始杠杆倍数
	NotionalCap      float64 `json:"notionalCap"`      // 该层对应的名义价值上限
	NotionalFloor    float64 `json:"notionalFloor"`    // 该层对应的名义价值下限
	MaintMarginRatio float64 `json:"maintMarginRatio"` // 该层对应的维持保证金率
	Cum              float64 `json:"cum"`              // 速算数
}
