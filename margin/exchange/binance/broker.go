package binance

import (
	"fmt"

	"github.com/taomu/lin-trader/margin/data"
	bndata "github.com/taomu/lin-trader/margin/exchange/binance/data"
)

type Broker struct {
	ApiInfo *data.ApiInfo
}

func NewBroker(apiInfo *data.ApiInfo) *Broker {
	return &Broker{
		ApiInfo: apiInfo,
	}
}
func (b *Broker) Test() {
	fmt.Println("binance test")
}
func (b *Broker) GetAllPairs() ([]bndata.Pair, error) {
	resp, err := NewRestApi().GetAllPairs(map[string]interface{}{}, b.ApiInfo)
	if err != nil {
		return nil, err
	}
	return bndata.TransferBinancePair(resp)
}
func (b *Broker) ListSchedule() ([]bndata.Schedule, error) {
	resp, err := NewRestApi().ListSchedule(map[string]interface{}{}, b.ApiInfo)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", resp)
	return bndata.TransferBinanceSchedule(resp)
}
func (b *Broker) DelistSchedule() ([]bndata.Schedule, error) {
	resp, err := NewRestApi().DelistSchedule(map[string]interface{}{}, b.ApiInfo)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", resp)
	return bndata.TransferBinanceSchedule(resp)
}
