package okx

import (
	"fmt"

	"github.com/taomu/lin-trader/pkg/types"
	"github.com/taomu/lin-trader/spot/data"
)

type Broker struct {
}

func NewBroker(apiInfo *types.ApiInfo) *Broker {
	return &Broker{}
}

func (b *Broker) Test() {
	fmt.Println("okx test")
}

func (b *Broker) GetSymbolInfos() ([]data.SymbolInfo, error) {
	resp, err := NewRestApi().Instruments(map[string]interface{}{
		"instType": "SPOT",
	})
	if err != nil {
		return nil, err
	}
	return data.TransferOkxSymbolInfo(resp)
}
