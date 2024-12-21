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

type TradingState uint8

const (
	STATE_HALTED    TradingState = 'H'
	STATE_PAUSED    TradingState = 'P'
	STATE_QUOTATION TradingState = 'Q'
	STATE_TRADING   TradingState = 'T'
)

type StockTradingAction struct {
	Stock          string
	Reason         string
	Timestamp      time.Duration
	StockLocate    uint16
	TrackingNumber uint16
	TradingState   TradingState
	Reserved       uint8
}

func (t StockTradingAction) Type() uint8 {
	return MESSAGE_STOCK_TRADING_ACTION
}

func (t StockTradingAction) Bytes() []byte {
	data := make([]byte, stockTradingActionSize)

	data[0] = MESSAGE_STOCK_TRADING_ACTION
	binary.BigEndian.PutUint16(data[1:3], t.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(t.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], t.TrackingNumber)

	copy(data[11:19], []byte(fmt.Sprintf("%-8s", t.Stock)))

	data[19] = byte(t.TradingState)
	data[20] = t.Reserved

	copy(data[21:25], []byte(fmt.Sprintf("%-4s", t.Reason)))

	return data
}

func ParseStockTradingAction(data []byte) (StockTradingAction, error) {
	if len(data) != stockTradingActionSize {
		return StockTradingAction{}, NewInvalidPacketSize(stockTradingActionSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return StockTradingAction{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Stock:          strings.TrimSpace(string(data[11:19])),
		TradingState:   TradingState(data[19]),
		Reserved:       data[20],
		Reason:         strings.TrimSpace(string(data[21:25])),
	}, nil
}

func (a StockTradingAction) String() string {
	return fmt.Sprintf("[Stock Trading Action]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Stock: %v\n"+
		"Trading State: %v\n"+
		"Reserved: %v\n"+
		"Reason: %v\n",
		a.StockLocate, a.TrackingNumber, a.Timestamp, a.Stock,
		a.TradingState, a.Reserved, a.Reason,
	)
}

func (t TradingState) String() string {
	switch t {
	case STATE_HALTED:
		return "Halted"
	case STATE_PAUSED:
		return "Paused"
	case STATE_QUOTATION:
		return "Quoatation only period for cross SRO halt or pause"
	case STATE_TRADING:
		return "Trading on Nasdaq"
	}

	return "Unknown Trading State"
}
