/*
 * Copyright (c) 2022 Mark Edward Winter
 */
package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type OrderExecuted struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Reference      uint64
	Shares         uint32
	MatchNumber    uint64
}

type OrderExecutedPrice struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Reference      uint64
	Shares         uint32
	MatchNumber    uint64
	Printable      bool
	ExecutionPrice uint32
}

func MakeOrderExecuted(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return OrderExecuted{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Reference:      binary.BigEndian.Uint64(data[11:19]),
		Shares:         binary.BigEndian.Uint32(data[19:23]),
		MatchNumber:    binary.BigEndian.Uint64(data[23:31]),
	}
}

func MakeOrderExecutedPrice(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	printable := false
	if data[31] == 'Y' {
		printable = true
	}

	return OrderExecutedPrice{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Reference:      binary.BigEndian.Uint64(data[11:19]),
		Shares:         binary.BigEndian.Uint32(data[19:23]),
		MatchNumber:    binary.BigEndian.Uint64(data[23:31]),
		Printable:      printable,
		ExecutionPrice: binary.BigEndian.Uint32(data[32:]),
	}
}

func (o OrderExecuted) String() string {
	return fmt.Sprintf("[Order Executed]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Reference: %v\n"+
		"Shares: %v\n"+
		"Match Number: %v\n",
		o.StockLocate, o.TrackingNumber, o.Timestamp,
		o.Reference, o.Shares, o.MatchNumber,
	)
}

func (o OrderExecutedPrice) String() string {
	return fmt.Sprintf("[Order Executed with Price]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Reference: %v\n"+
		"Shares: %v\n"+
		"Match Number: %v\n"+
		"Printable: %v\n"+
		"Execution Price: %v\n",
		o.StockLocate, o.TrackingNumber, o.Timestamp,
		o.Reference, o.Shares, o.MatchNumber,
		o.Printable, float64(o.ExecutionPrice)/10000,
	)
}
