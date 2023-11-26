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
	Shares         uint32
	CrossPrice     uint32
	StockLocate    uint16
	TrackingNumber uint16
	CrossType      CrossType
}

func MakeTradeCross(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return TradeCross{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Shares:         binary.BigEndian.Uint32(data[11:19]),
		Stock:          strings.TrimSpace(string(data[19:27])),
		CrossPrice:     binary.BigEndian.Uint32(data[27:31]),
		MatchNumber:    binary.BigEndian.Uint64(data[31:39]),
		CrossType:      CrossType(data[39]),
	}
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
