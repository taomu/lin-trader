package data

import (
	"strconv"

	"github.com/taomu/lin-trader/futures/types"
)

type PremiumIndexResp struct {
	Symbol               string `json:"symbol"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	InterestRate         string `json:"interestRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	Time                 int64  `json:"time"`
}

func TransferBinanceFundingRate(resp *PremiumIndexResp) (*types.FundingRate, error) {
	rate, err := strconv.ParseFloat(resp.LastFundingRate, 64)
	if err != nil {
		return nil, err
	}
	return &types.FundingRate{
		Symbol:       resp.Symbol,
		SymbolOri:    resp.Symbol,
		Rate:         rate,
		NextSettleTs: resp.NextFundingTime,
	}, nil
}
