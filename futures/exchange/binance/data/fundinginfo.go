package data

import (
	"encoding/json"
	"strconv"
)

type FundingInfo struct {
	Symbol   string `json:"symbol"`
	MaxRate  string `json:"adjustedFundingRateCap"`   // 最大资金费率
	MinRate  string `json:"adjustedFundingRateFloor"` // 最小资金费率
	Interval int    `json:"fundingIntervalHours"`     // 资金费率间隔 小时
}

func TransferBinanceFundingInfo(resp string) ([]FundingInfo, error) {
	var fundingInfo []FundingInfo
	err := json.Unmarshal([]byte(resp), &fundingInfo)
	if err != nil {
		return nil, err
	}

	for i := range fundingInfo {
		fundingInfo[i].MaxRate, err = formatFloatStr(fundingInfo[i].MaxRate)
		if err != nil {
			return nil, err
		}
		fundingInfo[i].MinRate, err = formatFloatStr(fundingInfo[i].MinRate)
		if err != nil {
			return nil, err
		}
	}

	return fundingInfo, nil
}

// floatStr 字符串去掉末尾的0
func formatFloatStr(floatStr string) (string, error) {
	// 转成 float64
	f, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		return "", err
	}
	// 再格式化成字符串
	// 'f' 表示小数格式，-1 表示由 Go 自动决定精度，去掉多余零
	s := strconv.FormatFloat(f, 'f', -1, 64)
	return s, nil
}
