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

type MarketCode uint8
type HaltAction uint8

const (
	MARKET_CODE_NASDAQ MarketCode = 'Q'
	MARKET_CODE_BX     MarketCode = 'B'
	MARKET_CODE_PSX    MarketCode = 'X'

	HALT_ACTION_HALT   HaltAction = 'H'
	HALT_ACTION_LIFTED HaltAction = 'T'
)

type OperationalHalt struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Stock          string
	MarketCode     MarketCode
	HaltAction     HaltAction
}

func MakeOperationalHalt(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return OperationalHalt{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Stock:          strings.TrimSpace(string(data[11:19])),
		MarketCode:     MarketCode(data[19]),
		HaltAction:     HaltAction(data[20]),
	}
}

func (h OperationalHalt) String() string {
	return fmt.Sprintf("[IPO Quotation]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Stock: %v\n"+
		"Market Code: %v\n"+
		"Halt Action: %v\n",
		h.StockLocate, h.TrackingNumber, h.Timestamp,
		h.Stock, h.MarketCode, h.HaltAction,
	)
}

func (m MarketCode) String() string {
	switch m {
	case MARKET_CODE_NASDAQ:
		return "Nasdaq"
	case MARKET_CODE_BX:
		return "BX"
	case MARKET_CODE_PSX:
		return "PSX"
	}

	return "Unknown MarketCode"
}

func (h HaltAction) String() string {
	switch h {
	case HALT_ACTION_HALT:
		return "Halted"
	case HALT_ACTION_LIFTED:
		return "Halt lifted, trading resumed"
	}

	return "Unknown HaltAction"
}
