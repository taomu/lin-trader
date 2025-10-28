package data

import (
	"strconv"
	"strings"

	cmTypes "github.com/taomu/lin-trader/futures/types"
)

type PositionsRes struct {
	Code string     `json:"code"`
	Data []Position `json:"data"`
	Msg  string     `json:"msg"`
}

type Position struct {
	Adl                    string        `json:"adl"`                    // ADL等级
	AvailPos               string        `json:"availPos"`               // 可平仓数量
	AvgPx                  string        `json:"avgPx"`                  // 开仓均价
	BaseBal                string        `json:"baseBal"`                // 币余额
	BaseBorrowed           string        `json:"baseBorrowed"`           // 已借币数量
	BaseInterest           string        `json:"baseInterest"`           // 币利息
	BePx                   string        `json:"bePx"`                   // 盈亏平衡价
	BizRefId               string        `json:"bizRefId"`               // 业务关联ID
	BizRefType             string        `json:"bizRefType"`             // 业务关联类型
	CTime                  string        `json:"cTime"`                  // 创建时间
	Ccy                    string        `json:"ccy"`                    // 币种
	ClSpotInUseAmt         string        `json:"clSpotInUseAmt"`         // 现货对冲占用数量
	CloseOrderAlgo         []interface{} `json:"closeOrderAlgo"`         // 平仓策略委托
	DeltaBS                string        `json:"deltaBS"`                // 美金本位持仓仓位delta
	DeltaPA                string        `json:"deltaPA"`                // 币本位持仓仓位delta
	Fee                    string        `json:"fee"`                    // 手续费
	FundingFee             string        `json:"fundingFee"`             // 资金费用
	GammaBS                string        `json:"gammaBS"`                // 美金本位持仓仓位gamma
	GammaPA                string        `json:"gammaPA"`                // 币本位持仓仓位gamma
	IdxPx                  string        `json:"idxPx"`                  // 指数价格
	Imr                    string        `json:"imr"`                    // 初始保证金
	InstId                 string        `json:"instId"`                 // 产品ID
	InstType               string        `json:"instType"`               // 产品类型
	Interest               string        `json:"interest"`               // 利息
	Last                   string        `json:"last"`                   // 最新成交价
	Lever                  string        `json:"lever"`                  // 杠杆倍数
	Liab                   string        `json:"liab"`                   // 负债额
	LiabCcy                string        `json:"liabCcy"`                // 负债币种
	LiqPenalty             string        `json:"liqPenalty"`             // 强平罚金
	LiqPx                  string        `json:"liqPx"`                  // 预估强平价
	Margin                 string        `json:"margin"`                 // 保证金余额
	MarkPx                 string        `json:"markPx"`                 // 标记价格
	MaxSpotInUseAmt        string        `json:"maxSpotInUseAmt"`        // 最大现货对冲占用数量
	MgnMode                string        `json:"mgnMode"`                // 保证金模式
	MgnRatio               string        `json:"mgnRatio"`               // 保证金率
	Mmr                    string        `json:"mmr"`                    // 维持保证金
	NotionalUsd            string        `json:"notionalUsd"`            // 以美金价值
	OptVal                 string        `json:"optVal"`                 // 期权市值
	PendingCloseOrdLiabVal string        `json:"pendingCloseOrdLiabVal"` // 挂平仓单负债价值
	Pnl                    string        `json:"pnl"`                    // 收益
	Pos                    string        `json:"pos"`                    // 持仓数量
	PosCcy                 string        `json:"posCcy"`                 // 持仓币种
	PosId                  string        `json:"posId"`                  // 持仓ID
	PosSide                string        `json:"posSide"`                // 持仓方向
	QuoteBal               string        `json:"quoteBal"`               // 计价货币余额
	QuoteBorrowed          string        `json:"quoteBorrowed"`          // 已借计价货币数量
	QuoteInterest          string        `json:"quoteInterest"`          // 计价货币利息
	RealizedPnl            string        `json:"realizedPnl"`            // 已实现收益
	SpotInUseAmt           string        `json:"spotInUseAmt"`           // 现货对冲占用数量
	SpotInUseCcy           string        `json:"spotInUseCcy"`           // 现货对冲占用币种
	ThetaBS                string        `json:"thetaBS"`                // 美金本位持仓仓位theta
	ThetaPA                string        `json:"thetaPA"`                // 币本位持仓仓位theta
	TradeId                string        `json:"tradeId"`                // 最新成交ID
	UTime                  string        `json:"uTime"`                  // 更新时间
	Upl                    string        `json:"upl"`                    // 未实现收益
	UplLastPx              string        `json:"uplLastPx"`              // 以最新成交价计算的未实现收益
	UplRatio               string        `json:"uplRatio"`               // 未实现收益率
	UplRatioLastPx         string        `json:"uplRatioLastPx"`         // 以最新成交价计算的未实现收益率
	UsdPx                  string        `json:"usdPx"`                  // 美金价格
	VegaBS                 string        `json:"vegaBS"`                 // 美金本位持仓仓位vega
	VegaPA                 string        `json:"vegaPA"`                 // 币本位持仓仓位vega
	NonSettleAvgPx         string        `json:"nonSettleAvgPx"`         // 非结算币种开仓均价
	SettledPnl             string        `json:"settledPnl"`             // 已结算收益
}

func TransformPositionToPos(positionsRes PositionsRes) []*cmTypes.Position {
	if len(positionsRes.Data) == 0 {
		return nil
	}
	result := make([]*cmTypes.Position, 0, len(positionsRes.Data))
	for _, pos := range positionsRes.Data {
		posAmt, err := strconv.ParseFloat(pos.Pos, 64)
		if err != nil {
			continue
		}
		parts := strings.Split(strings.Replace(pos.InstId, "-SWAP", "", -1), "-")
		symbol := strings.Join(parts, "")
		posSide := ""
		if pos.PosSide == "long" {
			posSide = "LONG"
		}
		if pos.PosSide == "short" {
			posSide = "SHORT"
		}
		entryPrice, err := strconv.ParseFloat(pos.AvgPx, 64)
		if err != nil {
			continue
		}
		result = append(result, &cmTypes.Position{
			Symbol:     symbol,
			PosAmt:     posAmt,
			PosSide:    posSide,
			EntryPrice: entryPrice,
		})
	}
	return result
}
