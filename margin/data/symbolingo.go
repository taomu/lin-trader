package data

import (
	"encoding/json"
	"fmt"
)

type SymbolInfo struct {
	SymbolOri              string // 原始交易对
	Symbol                 string // 标准交易对
	CtVal                  string // 合约面值
	QuoteAssetPrec         string // QuoteAsset精度
	QuotePrec              string // Quote精度
	LotPrec                string // 张数精度
	MinQty                 string // 最小下单数量
	IsSpotTradingAllowed   bool   //允许现货交易
	IsMarginTradingAllowed bool   //运行杠杆交易
	Status                 string // 状态
}

func TransferBinanceSymbolInfo(resp string) ([]SymbolInfo, error) {
	var apiResponse struct {
		Symbols []struct {
			Symbol                 string `json:"symbol"`
			QuoteAssetPrecision    int    `json:"quoteAssetPrecision"`
			BaseAssetPrecision     int    `json:"baseAssetPrecision"`
			QuotePrecision         int    `json:"quotePrecision"`
			IsSpotTradingAllowed   bool   `json:"isSpotTradingAllowed"`   //允许现货交易
			IsMarginTradingAllowed bool   `json:"isMarginTradingAllowed"` //允许杠杆交易
			Filters                []struct {
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
		// status := "TRADING"
		// if symbol.Status == "PENDING_TRADING" {
		// 	status = "PREOPEN"
		// }
		// //下线  交割中 已交割 结算中 均作为stop
		// if symbol.Status == "CLOSE" || symbol.Status == "DELIVERING" || symbol.Status == "DELIVERED" || symbol.Status == "SETTLING" {
		// 	status = "STOP"
		// }
		// //预结算 预交割 均作为PRESTOP
		// if symbol.Status == "PRE_SETTLE" || symbol.Status == "PRE_DELIVERING" {
		// 	status = "PRESTOP"
		// }
		symbolInfos = append(symbolInfos, SymbolInfo{
			SymbolOri:      symbol.Symbol,
			Symbol:         symbol.Symbol,
			QuoteAssetPrec: fmt.Sprintf("%d", symbol.QuoteAssetPrecision),
			QuotePrec:      fmt.Sprintf("%d", symbol.QuotePrecision),
			MinQty:         minQty,
			// 币安没有这些字段，设为空或默认值
			CtVal:                  "",
			LotPrec:                "",
			IsSpotTradingAllowed:   symbol.IsSpotTradingAllowed,
			IsMarginTradingAllowed: symbol.IsMarginTradingAllowed,
			Status:                 symbol.Status,
		})
	}
	return symbolInfos, nil
}
