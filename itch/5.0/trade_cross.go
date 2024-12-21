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

type CrossType uint8

const (
	CROSS_TYPE_NASDAQ_OPEN            CrossType = 'O'
	CROSS_TYPE_NASDAQ_CLOSE           CrossType = 'C'
	CROSS_TYPE_IPO_HALTED             CrossType = 'H'
	CROSS_TYPE_EXTENDED_TRADING_CLOSE CrossType = 'A'
)

type TradeCross struct {
	Stock          string
	Timestamp      time.Duration
	MatchNumber    uint64
	Shares         uint64
	CrossPrice     uint32
	StockLocate    uint16
	TrackingNumber uint16
	CrossType      CrossType
}

func (t TradeCross) Type() uint8 {
	return MESSAGE_TRADE_CROSS
}

func (t TradeCross) Bytes() []byte {
	data := make([]byte, tradeCrossSize)

	data[0] = MESSAGE_TRADE_CROSS
	binary.BigEndian.PutUint16(data[1:3], t.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(t.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], t.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], t.Shares)

	copy(data[19:27], []byte(fmt.Sprintf("%-8s", t.Stock)))

	binary.BigEndian.PutUint32(data[27:31], t.CrossPrice)
	binary.BigEndian.PutUint64(data[31:39], t.MatchNumber)

	data[39] = byte(t.CrossType)

	return data
}

func ParseTradeCross(data []byte) (TradeCross, error) {
	if len(data) != tradeCrossSize {
		return TradeCross{}, NewInvalidPacketSize(tradeCrossSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return TradeCross{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Shares:         binary.BigEndian.Uint64(data[11:19]),
		Stock:          strings.TrimSpace(string(data[19:27])),
		CrossPrice:     binary.BigEndian.Uint32(data[27:31]),
		MatchNumber:    binary.BigEndian.Uint64(data[31:39]),
		CrossType:      CrossType(data[39]),
	}, nil
}

func (o TradeCross) String() string {
	return fmt.Sprintf("[Trade Cross]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Shares: %v\n"+
		"Stock: %v\n"+
		"Cross Price: %v\n"+
		"Match Number: %v\n"+
		"Cross Type: %v\n",
		o.StockLocate, o.TrackingNumber, o.Timestamp,
		o.Shares, o.Stock, float64(o.CrossPrice)/10000,
		o.MatchNumber, o.CrossType,
	)
}

func (c CrossType) String() string {
	switch c {
	case CROSS_TYPE_NASDAQ_OPEN:
		return "Nasdaq Opening Cross"
	case CROSS_TYPE_NASDAQ_CLOSE:
		return "Nasdaq Closing Cross"
	case CROSS_TYPE_IPO_HALTED:
		return "Cross for IPO and halted / paused securities"
	}

	return "Unknown CrossType"
}
