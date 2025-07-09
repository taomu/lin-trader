package data

import "encoding/json"

type Schedule struct {
	ListTime              string   `json:"listTime"`
	DelistTime            string   `json:"delistTime"`
	CrossMarginAssets     []string `json:"crossMarginAssets"`
	IsolatedMarginSymbols []string `json:"isolatedMarginSymbols"`
}

func TransferBinanceSchedule(resp string) ([]Schedule, error) {
	var schedules []Schedule
	err := json.Unmarshal([]byte(resp), &schedules)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}
