package data

import (
	"strconv"

	"github.com/taomu/lin-trader/futures/types"
)

// WebSocket 深度事件（现货/合约基本一致）
type WsDepthUpdateRes struct {
	EventType         string     `json:"e"`  // "depthUpdate"
	EventTime         int64      `json:"E"`  //推送事件时间
	Symbol            string     `json:"s"`  //交易对
	Time              int64      `json:"T"`  //深度更新时间
	FirstUpdateID     int64      `json:"U"`  //当前推送第一个updateId
	FinalUpdateID     int64      `json:"u"`  //当前推送最后一个updateId
	LastFinalUpdateId int64      `json:"pu"` //上一次推送的最后一个updateId
	Bids              [][]string `json:"b"`  // [[price, qty], ...]
	Asks              [][]string `json:"a"`  // [[price, qty], ...]
}

// 深度快照
type DepthSnapshot struct {
	FinalID    int64      `json:"u"`    //当前推送最后一个updateId
	LastPushId int64      `json:"pu"`   //上一次推送的最后一个updateId
	Bids       [][]string `json:"bids"` // [[price, qty], ...]
	Asks       [][]string `json:"asks"` // [[price, qty], ...]
	Time       int64      `json:"time"` //深度更新时间
}

// OrderBook：价->量；量为0表示删除
type OrderBook struct {
	Bids map[float64]float64
	Asks map[float64]float64
}

// rest 请求depht数据
type DepthRes struct {
	FinalUpdateID int64      `json:"lastUpdateId"`
	Time          int64      `json:"time"`
	Bids          [][]string `json:"bids"`
	Asks          [][]string `json:"asks"`
}

func TransferBinanceDepthRes(depthRes DepthRes) *types.Depth {
	d := &types.Depth{
		Time:          depthRes.Time,
		FinalUpdateId: depthRes.FinalUpdateID,
		Asks:          make([]*types.BookItem, 0),
		Bids:          make([]*types.BookItem, 0),
	}
	for _, bid := range depthRes.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		qty, _ := strconv.ParseFloat(bid[1], 64)
		d.Bids = append(d.Bids, &types.BookItem{
			Price: price,
			Qty:   qty,
		})
	}
	for _, ask := range depthRes.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		qty, _ := strconv.ParseFloat(ask[1], 64)
		d.Asks = append(d.Asks, &types.BookItem{
			Price: price,
			Qty:   qty,
		})
	}
	return d
}

func TransferBinanceWsDepthUpdateRes(depthUpdate WsDepthUpdateRes) *types.Depth {
	d := &types.Depth{
		Time:              depthUpdate.Time,
		Symbol:            depthUpdate.Symbol,
		FinalUpdateId:     depthUpdate.FinalUpdateID,
		LastFinalUpdateId: depthUpdate.LastFinalUpdateId,
		Asks:              make([]*types.BookItem, 0),
		Bids:              make([]*types.BookItem, 0),
	}
	for _, bid := range depthUpdate.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		qty, _ := strconv.ParseFloat(bid[1], 64)
		d.Bids = append(d.Bids, &types.BookItem{
			Price: price,
			Qty:   qty,
		})
	}
	for _, ask := range depthUpdate.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		qty, _ := strconv.ParseFloat(ask[1], 64)
		d.Asks = append(d.Asks, &types.BookItem{
			Price: price,
			Qty:   qty,
		})
	}
	return d
}
