package data

import (
	"fmt"
	"math"
	"strconv"

	"github.com/taomu/lin-trader/futures/types"
)

type WsPositionsMsg struct {
	Arg struct {
		Channel  string `json:"channel"`
		InstType string `json:"instType"`
		Uid      string `json:"uid"`
	} `json:"arg"`
	EventType string       `json:"eventType"`
	CurPage   int          `json:"curPage"`
	LastPage  bool         `json:"lastPage"`
	Data      []WsPosition `json:"data"`
}

type WsPosition struct {
	Adl                    string        `json:"adl"`
	AvailPos               string        `json:"availPos"`
	AvgPx                  string        `json:"avgPx"`
	BaseBal                string        `json:"baseBal"`
	BaseBorrowed           string        `json:"baseBorrowed"`
	BaseInterest           string        `json:"baseInterest"`
	BePx                   string        `json:"bePx"`
	BizRefId               string        `json:"bizRefId"`
	BizRefType             string        `json:"bizRefType"`
	CTime                  string        `json:"cTime"`
	Ccy                    string        `json:"ccy"`
	ClSpotInUseAmt         string        `json:"clSpotInUseAmt"`
	CloseOrderAlgo         []interface{} `json:"closeOrderAlgo"`
	DeltaBS                string        `json:"deltaBS"`
	DeltaPA                string        `json:"deltaPA"`
	Fee                    string        `json:"fee"`
	FundingFee             string        `json:"fundingFee"`
	GammaBS                string        `json:"gammaBS"`
	GammaPA                string        `json:"gammaPA"`
	HedgedPos              string        `json:"hedgedPos"`
	IdxPx                  string        `json:"idxPx"`
	Imr                    string        `json:"imr"`
	InstId                 string        `json:"instId"`
	InstType               string        `json:"instType"`
	Interest               string        `json:"interest"`
	Last                   string        `json:"last"`
	Lever                  string        `json:"lever"`
	Liab                   string        `json:"liab"`
	LiabCcy                string        `json:"liabCcy"`
	LiqPenalty             string        `json:"liqPenalty"`
	LiqPx                  string        `json:"liqPx"`
	Margin                 string        `json:"margin"`
	MarkPx                 string        `json:"markPx"`
	MaxSpotInUseAmt        string        `json:"maxSpotInUseAmt"`
	MgnMode                string        `json:"mgnMode"`
	MgnRatio               string        `json:"mgnRatio"`
	Mmr                    string        `json:"mmr"`
	NonSettleAvgPx         string        `json:"nonSettleAvgPx"`
	NotionalUsd            string        `json:"notionalUsd"`
	OptVal                 string        `json:"optVal"`
	PTime                  string        `json:"pTime"`
	PendingCloseOrdLiabVal string        `json:"pendingCloseOrdLiabVal"`
	Pnl                    string        `json:"pnl"`
	Pos                    string        `json:"pos"`
	PosCcy                 string        `json:"posCcy"`
	PosId                  string        `json:"posId"`
	PosSide                string        `json:"posSide"`
	QuoteBal               string        `json:"quoteBal"`
	QuoteBorrowed          string        `json:"quoteBorrowed"`
	QuoteInterest          string        `json:"quoteInterest"`
	RealizedPnl            string        `json:"realizedPnl"`
	SettledPnl             string        `json:"settledPnl"`
	SpotInUseAmt           string        `json:"spotInUseAmt"`
	SpotInUseCcy           string        `json:"spotInUseCcy"`
	ThetaBS                string        `json:"thetaBS"`
	ThetaPA                string        `json:"thetaPA"`
	TradeId                string        `json:"tradeId"`
	UTime                  string        `json:"uTime"`
	Upl                    string        `json:"upl"`
	UplLastPx              string        `json:"uplLastPx"`
	UplRatio               string        `json:"uplRatio"`
	UplRatioLastPx         string        `json:"uplRatioLastPx"`
	UsdPx                  string        `json:"usdPx"`
	VegaBS                 string        `json:"vegaBS"`
	VegaPA                 string        `json:"vegaPA"`
}

func TransformWsPositionsToWsData(msg WsPositionsMsg, symbolInfos map[string]types.SymbolInfo, toStdSymbol func(symbol string) (string, error)) (*types.WsData, error) {
	//打印symbolinfos
	wsdata := &types.WsData{
		DataType: types.WsDataTypePosition,
	}
	if len(msg.Data) == 0 {
		wsdata.DataTs = 0
		wsdata.Position = []*types.Position{}
		return wsdata, nil
	}
	wsdata.DataTs, _ = strconv.ParseInt(msg.Data[0].PTime, 10, 64)
	wsdata.Position = make([]*types.Position, 0, len(msg.Data))
	for _, pos := range msg.Data {
		symbol, err := toStdSymbol(pos.InstId)
		if err != nil {
			return nil, fmt.Errorf("TransformWsPositionsToWsData toStdSymbol err: %w", err)
		}
		symbolInfo, ok := symbolInfos[symbol]
		if !ok {
			return nil, fmt.Errorf("TransformWsPositionsToWsData symbol not found: %s", symbol)
		}
		posAmt, _ := strconv.ParseFloat(pos.Pos, 64)
		//格式化为symbolInfo.QtyPrec的精度
		posAmt = math.Round(posAmt*math.Pow10(symbolInfo.QtyPrec)) / math.Pow10(symbolInfo.QtyPrec)

		posSide := types.PosSideBoth
		if pos.PosSide == "short" {
			posSide = types.PosSideShort
		}
		if pos.PosSide == "long" {
			posSide = types.PosSideLong
		}
		upl, _ := strconv.ParseFloat(pos.Upl, 64)
		entryPrice, _ := strconv.ParseFloat(pos.AvgPx, 64)
		wsdata.Position = append(wsdata.Position, &types.Position{
			Symbol:           symbol,
			PosSide:          posSide,
			PosAmt:           posAmt,
			EntryPrice:       entryPrice,
			UnrealizedProfit: upl,
		})
	}
	return wsdata, nil
}
