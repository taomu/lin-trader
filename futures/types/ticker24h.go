package types

import (
	"encoding/json"
	"strings"
)

type Ticker24H struct {
	Symbol    string
	SymbolOri string
	Last      float64
	Open24    float64
	High24    float64
	Low24     float64
	Vol24     float64
	Rise24    float64
}

func TransferOkxTicker(resp string) ([]Ticker24H, error) {
	var apiResponse struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			InstId  string  `json:"instId"`
			Last    float64 `json:"last,string"`
			Open24h float64 `json:"open24h,string"`
			High24h float64 `json:"high24h,string"`
			Low24h  float64 `json:"low24h,string"`
			Vol24h  float64 `json:"vol24h,string"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(resp), &apiResponse); err != nil || apiResponse.Code != "0" {
		return nil, err
	}
	var tickers []Ticker24H
	for _, item := range apiResponse.Data {
		rise24 := (item.Last - item.Open24h) / item.Open24h
		parts := strings.Split(strings.Replace(item.InstId, "-SWAP", "", -1), "-")
		symbol := strings.Join(parts, "")
		tickers = append(tickers, Ticker24H{
			Symbol:    symbol,
			SymbolOri: item.InstId,
			Last:      item.Last,
			Open24:    item.Open24h,
			High24:    item.High24h,
			Low24:     item.Low24h,
			Vol24:     item.Vol24h,
			Rise24:    float64(int(rise24*10000)) / 10000, // 保留4位小数
		})
	}
	return tickers, nil
}
