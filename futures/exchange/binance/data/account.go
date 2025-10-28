package data

import (
	"strconv"

	cmTypes "github.com/taomu/lin-trader/futures/types"
)

type AccountRes struct {
	TotalInitialMargin          string     `json:"totalInitialMargin"`          // 当前所需起始保证金总额(存在逐仓请忽略), 仅计算usdt资产
	TotalMaintMargin            string     `json:"totalMaintMargin"`            // 维持保证金总额, 仅计算usdt资产
	TotalWalletBalance          string     `json:"totalWalletBalance"`          // 账户总余额, 仅计算usdt资产
	TotalUnrealizedProfit       string     `json:"totalUnrealizedProfit"`       // 持仓未实现盈亏总额, 仅计算usdt资产
	TotalMarginBalance          string     `json:"totalMarginBalance"`          // 保证金总余额, 仅计算usdt资产
	TotalPositionInitialMargin  string     `json:"totalPositionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalOpenOrderInitialMargin string     `json:"totalOpenOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalCrossWalletBalance     string     `json:"totalCrossWalletBalance"`     // 全仓账户余额, 仅计算usdt资产
	TotalCrossUnPnl             string     `json:"totalCrossUnPnl"`             // 全仓持仓未实现盈亏总额, 仅计算usdt资产
	AvailableBalance            string     `json:"availableBalance"`            // 可用余额, 仅计算usdt资产
	MaxWithdrawAmount           string     `json:"maxWithdrawAmount"`           // 最大可转出余额, 仅计算usdt资产
	Assets                      []Asset    `json:"assets"`                      // 资产列表
	Positions                   []Position `json:"positions"`                   // 持仓列表，仅有仓位或挂单的交易对会被返回
}

type Asset struct {
	Asset                  string `json:"asset"`                  // 资产
	WalletBalance          string `json:"walletBalance"`          // 余额
	UnrealizedProfit       string `json:"unrealizedProfit"`       // 未实现盈亏
	MarginBalance          string `json:"marginBalance"`          // 保证金余额
	MaintMargin            string `json:"maintMargin"`            // 维持保证金
	InitialMargin          string `json:"initialMargin"`          // 当前所需起始保证金
	PositionInitialMargin  string `json:"positionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格)
	CrossWalletBalance     string `json:"crossWalletBalance"`     // 全仓账户余额
	CrossUnPnl             string `json:"crossUnPnl"`             // 全仓持仓未实现盈亏
	AvailableBalance       string `json:"availableBalance"`       // 可用余额
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`      // 最大可转出余额
	UpdateTime             int64  `json:"updateTime"`             // 更新时间
}

type Position struct {
	Symbol           string `json:"symbol"`           // 交易对
	PositionSide     string `json:"positionSide"`     // 持仓方向
	PositionAmt      string `json:"positionAmt"`      // 持仓数量
	UnrealizedProfit string `json:"unrealizedProfit"` // 持仓未实现盈亏
	IsolatedMargin   string `json:"isolatedMargin"`   // 逐仓保证金
	Notional         string `json:"notional"`         // 名义价值
	IsolatedWallet   string `json:"isolatedWallet"`   // 逐仓钱包余额
	InitialMargin    string `json:"initialMargin"`    // 持仓所需起始保证金(基于最新标记价格)
	MaintMargin      string `json:"maintMargin"`      // 当前杠杆下用户可用的最大名义价值
	UpdateTime       int64  `json:"updateTime"`       // 更新时间
}

// 转化结果为仓位
func TransferBinanceAccountResToPos(accountRes AccountRes) []*cmTypes.Position {
	positions := make([]*cmTypes.Position, 0, len(accountRes.Positions))
	for _, p := range accountRes.Positions {
		posAmt, _ := strconv.ParseFloat(p.PositionAmt, 64)
		unrealizedProfit, _ := strconv.ParseFloat(p.UnrealizedProfit, 64)
		if posAmt == 0 {
			continue
		}
		positions = append(positions, &cmTypes.Position{
			Symbol:           p.Symbol,
			PosSide:          p.PositionSide,
			PosAmt:           posAmt,
			UnrealizedProfit: unrealizedProfit,
		})
	}
	return positions
}
