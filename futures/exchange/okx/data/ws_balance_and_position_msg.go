package data

import (
	"encoding/json"
	"strconv"

	"github.com/taomu/lin-trader/futures/types"
)

type WsBalanceAndPositionMsg struct {
	Event string `json:"event"`
	Arg   struct {
		Channel string `json:"channel"`
		Uid     string `json:"uid"`
	} `json:"arg"`
	Data []WsBalanceAndPositionData `json:"data"`
	Code string                     `json:"code"`
}

type WsBalanceAndPositionData struct {
	BalData   []BalanceDetail   `json:"balData"`
	PosData   []PositionDetail  `json:"posData"`
	Trades    []json.RawMessage `json:"trades"`
	EventType string            `json:"eventType"`
	PTime     string            `json:"pTime"`
}

type BalanceDetail struct {
	CashBal string `json:"cashBal"`
	Ccy     string `json:"ccy"`
	UTime   string `json:"uTime"`
}

type PositionDetail struct {
	InstId  string `json:"instId"`
	Pos     string `json:"pos"`
	PosSide string `json:"posSide"`
	AvgPx   string `json:"avgPx"`
}

func TransformBalanceAndPositionToWsData(msg WsBalanceAndPositionMsg) (*types.WsData, error) {
	balanceUSDT := 0.0
	var err error
	if len(msg.Data) != 0 {
		for _, bal := range msg.Data[0].BalData {
			if bal.Ccy == "USDT" {
				balanceUSDT, err = strconv.ParseFloat(bal.CashBal, 64)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	wsBalance := types.WsBalance{
		BalanceAll:   balanceUSDT,
		BalanceAvail: balanceUSDT,
	}
	var wsdata types.WsData
	wsdata.DataType = types.WsDataTypeBalance
	if len(msg.Data) > 0 {
		wsdata.DataTs, err = strconv.ParseInt(msg.Data[0].PTime, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	wsdata.Balance = wsBalance
	return &wsdata, nil
}
