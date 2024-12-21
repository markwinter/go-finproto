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

	data[0] = MESSAGE_TRADE_NON_CROSS
	binary.BigEndian.PutUint16(data[1:3], t.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(t.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], t.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], t.Reference)

	data[19] = byte(t.OrderIndicator)

	binary.BigEndian.PutUint32(data[20:24], t.Shares)

	copy(data[24:32], []byte(fmt.Sprintf("%-8s", t.Stock)))

	binary.BigEndian.PutUint32(data[32:36], t.Price)
	binary.BigEndian.PutUint64(data[36:], t.MatchNumber)

	return data
}

func ParseTradeNonCross(data []byte) (TradeNonCross, error) {
	if len(data) != tradeNonCrossSize {
		return TradeNonCross{}, NewInvalidPacketSize(tradeNonCrossSize, len(data))
	}

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
	}, nil
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
