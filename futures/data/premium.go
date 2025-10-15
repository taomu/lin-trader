package data

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Premium struct {
	Symbol       string
	SymbolOri    string
	Rate         float64
	NextSettleTs int64
	MaxRate      float64
	MinRate      float64
	SettlePeriod int64
}

func TransferBinancePremium(resp string) ([]Premium, error) {
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(resp), &raw); err != nil {
		return nil, err
	}

	// 检查是否以 '[' 开头来判断是数组还是单个对象
	if len(raw) > 0 && raw[0] == '[' {
		// 处理数组情况
		var binanceData []struct {
			Symbol          string `json:"symbol"`
			LastFundingRate string `json:"lastFundingRate"`
			NextFundingTime int64  `json:"nextFundingTime"`
			Time            int64  `json:"time"`
		}

		if err := json.Unmarshal(raw, &binanceData); err != nil {
			return nil, err
		}

		var premiums []Premium
		for _, item := range binanceData {
			rate, _ := strconv.ParseFloat(item.LastFundingRate, 64)
			premiums = append(premiums, Premium{
				Symbol:       item.Symbol,
				Rate:         rate,
				NextSettleTs: item.NextFundingTime,
			})
		}
		return premiums, nil
	} else {
		// 处理单个对象情况
		var binanceItem struct {
			Symbol          string `json:"symbol"`
			LastFundingRate string `json:"lastFundingRate"`
			NextFundingTime int64  `json:"nextFundingTime"`
			Time            int64  `json:"time"`
		}

		if err := json.Unmarshal(raw, &binanceItem); err != nil {
			return nil, err
		}

		rate, _ := strconv.ParseFloat(binanceItem.LastFundingRate, 64)
		return []Premium{{
			Symbol:       binanceItem.Symbol,
			Rate:         rate,
			NextSettleTs: binanceItem.NextFundingTime,
		}}, nil
	}
}

func TransferOkxPremium(resp string) ([]Premium, error) {
	type okxRespData struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			InstId          string `json:"instId"`
			InstType        string `json:"instType"`
			FundingTime     string `json:"fundingTime"`
			NextFundingTime string `json:"nextFundingTime"`
			FundingRate     string `json:"fundingRate"`
			MaxFundingRate  string `json:"maxFundingRate"`
			MinFundingRate  string `json:"minFundingRate"`
		} `json:"data"`
	}

	var okxResp okxRespData
	if err := json.Unmarshal([]byte(resp), &okxResp); err != nil {
		return nil, err
	}

	var premiums []Premium
	for _, item := range okxResp.Data {
		rate, _ := strconv.ParseFloat(item.FundingRate, 64)
		maxRate, _ := strconv.ParseFloat(item.MaxFundingRate, 64)
		minRate, _ := strconv.ParseFloat(item.MinFundingRate, 64)
		parts := strings.Split(strings.Replace(item.InstId, "-SWAP", "", -1), "-")
		symbol := strings.Join(parts, "")
		nextSettleTs, _ := strconv.ParseInt(item.NextFundingTime, 10, 64)
		fundingTime, _ := strconv.ParseInt(item.FundingTime, 10, 64)
		premiums = append(premiums, Premium{
			Symbol:       symbol,
			SymbolOri:    item.InstId,
			Rate:         rate,
			NextSettleTs: fundingTime,
			MaxRate:      maxRate,
			MinRate:      minRate,
			SettlePeriod: nextSettleTs - fundingTime,
		})
	}
	return premiums, nil
}
