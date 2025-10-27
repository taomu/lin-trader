package unifiedbn

import (
	"encoding/json"
	"fmt"

	"github.com/taomu/lin-trader/pkg/lintypes"
	"github.com/taomu/lin-trader/unifiedbn/data"
)

type Broker struct {
	Key    string
	Secret string
}

func NewBroker(key, secret string) *Broker {
	return &Broker{
		Key:    key,
		Secret: secret,
	}
}

func (b *Broker) Test() {
	fmt.Println("unifiedbn test")
}

// accountResponse is used to unmarshal the response from /papi/v1/account
type accountResponse struct {
	Positions []positionDetail `json:"positions"`
}

// positionDetail is used to unmarshal the position details from the account response
type positionDetail struct {
	Symbol           string      `json:"symbol"`
	PositionAmt      json.Number `json:"positionAmt"`
	PositionSide     string      `json:"positionSide"`
	EntryPrice       json.Number `json:"entryPrice"`
	UnrealizedProfit json.Number `json:"unrealizedProfit"`
}

func (b *Broker) GetFuturesPositions() ([]data.Position, error) {
	ra := NewRestApi()
	apiInfo := &lintypes.ApiInfo{
		Key:    b.Key,
		Secret: b.Secret,
	}
	params := make(map[string]interface{})
	resp, err := ra.Account(params, apiInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %v", err)
	}

	var accResp accountResponse
	err = json.Unmarshal([]byte(resp), &accResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal account response: %v", err)
	}

	var positions []data.Position
	for _, p := range accResp.Positions {
		posAmt, err := p.PositionAmt.Float64()
		if err != nil {
			// You might want to log this error
			continue
		}

		// Only include open positions
		if posAmt == 0 {
			continue
		}

		entryPrice, err := p.EntryPrice.Float64()
		if err != nil {
			// You might want to log this error
			continue
		}
		unrealizedProfit, err := p.UnrealizedProfit.Float64()
		if err != nil {
			// You might want to log this error
			continue
		}

		positions = append(positions, data.Position{
			Symbol:           p.Symbol,
			PosSide:          p.PositionSide,
			PosAmt:           posAmt,
			EntryPrice:       entryPrice,
			UnrealizedProfit: unrealizedProfit,
		})
	}

	return positions, nil
}
