package data

import (
	"strconv"

	"github.com/taomu/lin-trader/futures/types"
)

type FundingRateResp struct {
	Code string            `json:"code"`
	Data []FundingRateItem `json:"data"`
	Msg  string            `json:"msg"`
}

type FundingRateItem struct {
	FormulaType     string `json:"formulaType"`
	FundingRate     string `json:"fundingRate"`
	FundingTime     string `json:"fundingTime"`
	ImpactValue     string `json:"impactValue"`
	InstId          string `json:"instId"`
	InstType        string `json:"instType"`
	InterestRate    string `json:"interestRate"`
	MaxFundingRate  string `json:"maxFundingRate"`
	Method          string `json:"method"`
	MinFundingRate  string `json:"minFundingRate"`
	NextFundingRate string `json:"nextFundingRate"`
	NextFundingTime string `json:"nextFundingTime"`
	Premium         string `json:"premium"`
	SettFundingRate string `json:"settFundingRate"`
	SettState       string `json:"settState"`
	Ts              string `json:"ts"`
}

func TransferOkxFundingRate(resp *FundingRateResp, toStdSymbol func(string) (string, error)) (*types.FundingRate, error) {
	var rates []*types.FundingRate
	for _, item := range resp.Data {
		rate, _ := strconv.ParseFloat(item.FundingRate, 64)
		nextSettleTs, _ := strconv.ParseInt(item.FundingTime, 10, 64)
		maxRate, _ := strconv.ParseFloat(item.MaxFundingRate, 64)
		minRate, _ := strconv.ParseFloat(item.MinFundingRate, 64)
		stdSymbol, err := toStdSymbol(item.InstId)
		if err != nil {
			return nil, err
		}
		rates = append(rates, &types.FundingRate{
			Symbol:       stdSymbol,
			SymbolOri:    item.InstId,
			Rate:         rate,
			NextSettleTs: nextSettleTs,
			MaxRate:      maxRate,
			MinRate:      minRate,
		})
	}
	return rates[0], nil
}
