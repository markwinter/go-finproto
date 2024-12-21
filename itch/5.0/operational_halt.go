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
	Stock          string
	Timestamp      time.Duration
	StockLocate    uint16
	TrackingNumber uint16
	MarketCode     MarketCode
	HaltAction     HaltAction
}

func (o OperationalHalt) Type() uint8 {
	return MESSAGE_OPERATIONAL_HALT
}

func (o OperationalHalt) Bytes() []byte {
	data := make([]byte, operationalHaltSize)

	data[0] = MESSAGE_OPERATIONAL_HALT
	binary.BigEndian.PutUint16(data[1:3], o.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(o.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], o.TrackingNumber)

	copy(data[11:19], []byte(fmt.Sprintf("%-8s", o.Stock)))

	data[19] = byte(o.MarketCode)
	data[20] = byte(o.HaltAction)

	return data
}

func ParseOperationalHalt(data []byte) (OperationalHalt, error) {
	if len(data) != operationalHaltSize {
		return OperationalHalt{}, NewInvalidPacketSize(operationalHaltSize, len(data))
	}

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
	}, nil
}

func (h OperationalHalt) String() string {
	return fmt.Sprintf("[Operational Halt]\n"+
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
