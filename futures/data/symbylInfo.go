package data

import (
	"encoding/json"
	"fmt"
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
	ExpTs     int64  // 下线时间 单位ms
}

func TransferBinanceSymbolInfo(resp string) ([]SymbolInfo, error) {
	var apiResponse struct {
		Symbols []struct {
			Symbol             string `json:"symbol"`
			PricePrecision     int    `json:"pricePrecision"`
			QuantityPrecision  int    `json:"quantityPrecision"`
			BaseAssetPrecision int    `json:"baseAssetPrecision"`
			QuotePrecision     int    `json:"quotePrecision"`
			Filters            []struct {
				FilterType  string `json:"filterType"`
				MinQty      string `json:"minQty,omitempty"`
				MinNotional string `json:"minNotional,omitempty"`
			} `json:"filters"`
			Status string `json:"status"`
		} `json:"symbols"`
	}
	if err := json.Unmarshal([]byte(resp), &apiResponse); err != nil {
		return nil, err
	}

	var symbolInfos []SymbolInfo
	for _, symbol := range apiResponse.Symbols {
		// 从filters中提取minQty
		var minQty string
		for _, filter := range symbol.Filters {
			if filter.FilterType == "LOT_SIZE" {
				minQty = filter.MinQty
				break
			}
		}
		status := "TRADING"
		if symbol.Status == "PENDING_TRADING" {
			status = "PREOPEN"
		}
		//下线  交割中 已交割 结算中 均作为stop
		if symbol.Status == "CLOSE" || symbol.Status == "DELIVERING" || symbol.Status == "DELIVERED" || symbol.Status == "SETTLING" {
			status = "STOP"
		}
		//预结算 预交割 均作为PRESTOP
		if symbol.Status == "PRE_SETTLE" || symbol.Status == "PRE_DELIVERING" {
			status = "PRESTOP"
		}
		symbolInfos = append(symbolInfos, SymbolInfo{
			SymbolOri: symbol.Symbol,
			Symbol:    symbol.Symbol,
			PricePrec: fmt.Sprintf("%d", symbol.PricePrecision),
			QtyPrec:   fmt.Sprintf("%d", symbol.QuantityPrecision),
			MinQty:    minQty,
			// 币安没有这些字段，设为空或默认值
			CtVal:   "",
			LotPrec: "",
			MinLot:  "",
			Status:  status,
		})
	}
	return symbolInfos, nil
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
			ExpTime  string `json:"expTime"`
			State    string `json:"state"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(resp), &apiResponse); err != nil || apiResponse.Code != "0" {
		return nil, err
	}

	var symbolInfos []SymbolInfo
	for _, item := range apiResponse.Data {
		parts := strings.Split(item.InstId, "-")
		symbol := parts[0] + parts[1]
		status := "TRADING"
		if item.State == "preopen" {
			status = "PREOPEN"
		}
		if item.State == "suspend" || item.State == "test" {
			status = "STOP"
		}
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
		})
	}
	return symbolInfos, nil
}
