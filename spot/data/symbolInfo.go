package data

import (
	"encoding/json"
	"strconv"
	"strings"
)

type SymbolInfo struct {
	SymbolOri string // 原始交易对
	Symbol    string // 标准交易对
	CtVal     string // 合约面值
	PricePrec string // 价格精度
	QtyPrec   string // 数量精度
	LotPrec   string // 张数精度
	MinLot    string // 最小下单张数
	MinQty    string // 最小下单数量
	Status    string // 状态 TRADING 交易中 、PREOPEN 预上线 、PRESTOP 预下线、STOP 下线
	OnlineTs  int64  // 下线时间 单位ms
	OfflineTs int64  // 上线时间 单位ms
}

func TransferOkxSymbolInfo(resp string) ([]SymbolInfo, error) {
	var apiResponse struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			InstId   string `json:"instId"`
			InstType string `json:"instType"`
			BaseCcy  string `json:"baseCcy"`
			QuoteCcy string `json:"quoteCcy"`
			CtVal    string `json:"ctVal"`
			TickSz   string `json:"tickSz"`
			LotSz    string `json:"lotSz"`
			MinSz    string `json:"minSz"`
			ExpTime  string `json:"expTime"`  //下线时间
			ListTime string `json:"listTime"` //上线时间 毫秒
			State    string `json:"state"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(resp), &apiResponse); err != nil || apiResponse.Code != "0" {
		return nil, err
	}

	var symbolInfos []SymbolInfo
	for _, item := range apiResponse.Data {
		parts := strings.Split(strings.Replace(item.InstId, "-SWAP", "", -1), "-")
		symbol := strings.Join(parts, "")
		status := "TRADING"
		if item.State == "preopen" {
			status = "PREOPEN"
		}
		if item.State == "suspend" || item.State == "test" {
			status = "STOP"
		}
		onlineTs, _ := strconv.ParseInt(item.ListTime, 10, 64)
		offlineTs, _ := strconv.ParseInt(item.ExpTime, 10, 64)
		symbolInfos = append(symbolInfos, SymbolInfo{
			SymbolOri: item.InstId,
			Symbol:    symbol,
			CtVal:     item.CtVal,
			PricePrec: item.TickSz,
			QtyPrec:   item.LotSz,
			LotPrec:   item.LotSz,
			MinLot:    item.MinSz,
			MinQty:    item.MinSz,
			Status:    status,
			OnlineTs:  onlineTs,
			OfflineTs: offlineTs,
		})
	}
	return symbolInfos, nil
}
