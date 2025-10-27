package okx

import (
	"fmt"

	"github.com/taomu/lin-trader/pkg/lintypes"
	"github.com/taomu/lin-trader/spot/data"
)

type Broker struct {
}

func NewBroker(apiInfo *lintypes.ApiInfo) *Broker {
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
