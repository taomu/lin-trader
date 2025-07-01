package okx

import (
	"fmt"

	"github.com/taomu/lin-trader/futures/data"
	bndata "github.com/taomu/lin-trader/futures/exchange/binance/data"
)

type Broker struct {
	ApiInfo *data.ApiInfo
}

func (b *Broker) Test() {
	fmt.Println("okx test")
}
func (b *Broker) GetPremium() ([]data.Premium, error) {
	resp, err := NewRestApi().GetPremium(map[string]interface{}{
		"instId": "ANY",
	})
	if err != nil {
		return nil, err
	}
	return data.TransferOkxPremium(resp)
}
func (b *Broker) GetFundingInfo() ([]bndata.FundingInfo, error) {
	return nil, nil
}
func (b *Broker) GetSymbolInfos() ([]data.SymbolInfo, error) {
	resp, err := NewRestApi().Instruments(map[string]interface{}{
		"instType": "FUTURES",
	})
	if err != nil {
		return nil, err
	}
	return data.TransferOkxSymbolInfo(resp)
}
