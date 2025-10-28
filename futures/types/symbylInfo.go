package types

import (
	"encoding/json"
	"strconv"
	"strings"
)

type SymbolInfo struct {
	SymbolOri string  // 原始交易对
	Symbol    string  // 标准交易对
	CtVal     float64 // 合约面值
	PricePrec int     // 价格精度
	QtyPrec   int     // 数量精度
	LotPrec   int     // 张数精度
	MinLot    float64 // 最小下单张数
	MinQty    float64 // 最小下单数量
	Status    string  // 状态 TRADING 交易中 、PREOPEN 预上线 、PRESTOP 预下线、STOP 下线
	OnlineTs  int64   // 下线时间 单位ms
	OfflineTs int64   // 上线时间 单位ms
	RuleType  string  // 合约类型 NORMAL:正常交易 PREMARKET 盘前交易
}

func TransferBinanceSymbolInfo(resp string) ([]SymbolInfo, error) {
	var apiResponse struct {
		Symbols []struct {
			Symbol             string `json:"symbol"`
			PricePrecision     int    `json:"pricePrecision"`
			QuantityPrecision  int    `json:"quantityPrecision"`
			BaseAssetPrecision int    `json:"baseAssetPrecision"`
			QuotePrecision     int    `json:"quotePrecision"`
			DeliveryDate       int64  `json:"deliveryDate"` //下架时间 毫秒时间戳
			OnboardDate        int64  `json:"onboardDate"`  //上线时间 毫秒时间戳
			Filters            []struct {
				FilterType  string `json:"filterType"`
				MinQty      string `json:"minQty,omitempty"`
				MinNotional string `json:"minNotional,omitempty"`
			} `json:"filters"`
			Status         string `json:"status"`
			UnderlyingType string `json:"underlyingType"`
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
		minQtyFloat, err := strconv.ParseFloat(minQty, 64)
		if err != nil {
			return nil, err
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
		ruleType := "NORMAL"
		if symbol.UnderlyingType == "PREMARKET" {
			ruleType = "PREMARKET"
		}
		symbolInfos = append(symbolInfos, SymbolInfo{
			SymbolOri: symbol.Symbol,
			Symbol:    symbol.Symbol,
			PricePrec: symbol.PricePrecision,
			QtyPrec:   symbol.QuantityPrecision,
			MinQty:    minQtyFloat,
			// 币安没有这些字段，设为空或默认值
			CtVal:     0,
			LotPrec:   0,
			MinLot:    0,
			Status:    status,
			OnlineTs:  symbol.OnboardDate,
			OfflineTs: symbol.DeliveryDate,
			RuleType:  ruleType,
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
			ExpTime  string `json:"expTime"`  //下线时间
			ListTime string `json:"listTime"` //上线时间 毫秒
			State    string `json:"state"`
			RuleType string `json:"ruleType"`
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
		ruleType := "NORMAL"
		if item.RuleType == "pre_market" {
			ruleType = "PREMARKET"
		}
		ctValFloat, err := strconv.ParseFloat(item.CtVal, 64)
		if err != nil {
			return nil, err
		}
		tickSzInt, err := strconv.Atoi(item.TickSz)
		if err != nil {
			return nil, err
		}
		lotSzInt, err := strconv.Atoi(item.LotSz)
		if err != nil {
			return nil, err
		}
		minSzFloat, err := strconv.ParseFloat(item.MinSz, 64)
		if err != nil {
			return nil, err
		}
		symbolInfos = append(symbolInfos, SymbolInfo{
			SymbolOri: item.InstId,
			Symbol:    symbol,
			CtVal:     ctValFloat,
			PricePrec: tickSzInt,
			QtyPrec:   lotSzInt,
			LotPrec:   lotSzInt,
			MinLot:    minSzFloat,
			MinQty:    minSzFloat,
			Status:    status,
			OnlineTs:  onlineTs,
			OfflineTs: offlineTs,
			RuleType:  ruleType,
		})
	}
	return symbolInfos, nil
}
