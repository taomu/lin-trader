package data

import (
	"encoding/json"
)

type OrderResp struct {
	Code    string          `json:"code"`
	Msg     string          `json:"msg"`
	Data    []OrderRespItem `json:"data"`
	InTime  string          `json:"inTime"`
	OutTime string          `json:"outTime"`
}

type OrderRespItem struct {
	ClOrdId string `json:"clOrdId"`
	OrdId   string `json:"ordId"`
	Tag     string `json:"tag"`
	Ts      string `json:"ts"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

func ParseOrderResp(resp string) (*OrderResp, error) {
	var r OrderResp
	if err := json.Unmarshal([]byte(resp), &r); err != nil {
		return nil, err
	}
	return &r, nil
}