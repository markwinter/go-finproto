/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type TradeNonCross struct {
	Stock          string
	Timestamp      time.Duration
	Reference      uint64
	MatchNumber    uint64
	Shares         uint32
	Price          uint32
	StockLocate    uint16
	TrackingNumber uint16
	OrderIndicator OrderIndicator
}

func (t TradeNonCross) Type() uint8 {
	return MESSAGE_TRADE_NON_CROSS
}

func (t TradeNonCross) Bytes() []byte {
	data := make([]byte, tradeNonCrossSize)
	// TODO: implement
	return data
}

func ParseTradeNonCross(data []byte) TradeNonCross {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return TradeNonCross{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Reference:      binary.BigEndian.Uint64(data[11:19]),
		OrderIndicator: OrderIndicator(data[19]),
		Shares:         binary.BigEndian.Uint32(data[20:24]),
		Stock:          strings.TrimSpace(string(data[24:32])),
		Price:          binary.BigEndian.Uint32(data[32:36]),
		MatchNumber:    binary.BigEndian.Uint64(data[36:]),
	}
}

func (o TradeNonCross) String() string {
	return fmt.Sprintf("[Trade Non-Cross]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Reference: %v\n"+
		"Order Indicator: %v\n"+
		"Shares: %v\n"+
		"Stock: %v\n"+
		"Price: %v\n"+
		"Match Number: %v\n",
		o.StockLocate, o.TrackingNumber, o.Timestamp,
		o.Reference, o.OrderIndicator, o.Shares, o.Stock,
		float64(o.Price)/10000, o.MatchNumber,
	)
}
