package binance

import (
	"fmt"

	"github.com/taomu/lin-trader/futures/data"
	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
)

type Broker struct {
	ApiInfo *data.ApiInfo
}

func (b *Broker) Test() {
	fmt.Println("binance test")
}
func (b *Broker) GetPremium() ([]data.Premium, error) {
	resp, err := NewRestApi().PremiumIndex(map[string]interface{}{
		"instId": "ANY",
	})
	if err != nil {
		return nil, err
	}
	return data.TransferBinancePremium(resp)
}
func (b *Broker) GetFundingInfo() ([]bndata.FundingInfo, error) {
	params := map[string]interface{}{}
	resp, err := NewRestApi().FundingInfo(params)
	if err != nil {
		return nil, err
	}
	return bndata.TransferBinanceFundingInfo(resp)
}
func (b *Broker) GetSymbolInfos() ([]data.SymbolInfo, error) {
	resp, err := NewRestApi().ExchangeInfo(map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return data.TransferBinanceSymbolInfo(resp)
}
