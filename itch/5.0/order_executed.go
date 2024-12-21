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
	Timestamp      time.Duration
	Reference      uint64
	MatchNumber    uint64
	Shares         uint32
	StockLocate    uint16
	TrackingNumber uint16
}

func (o OrderExecuted) Type() uint8 {
	return MESSAGE_ORDER_EXECUTED
}

func (o OrderExecuted) Bytes() []byte {
	data := make([]byte, orderExecutedSize)

	data[0] = MESSAGE_ORDER_EXECUTED
	binary.BigEndian.PutUint16(data[1:3], o.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(o.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], o.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], o.Reference)
	binary.BigEndian.PutUint32(data[19:23], o.Shares)
	binary.BigEndian.PutUint64(data[23:], o.MatchNumber)

	return data
}

type OrderExecutedPrice struct {
	Timestamp      time.Duration
	Reference      uint64
	MatchNumber    uint64
	Shares         uint32
	ExecutionPrice uint32
	StockLocate    uint16
	TrackingNumber uint16
	Printable      bool
}

func (o OrderExecutedPrice) Type() uint8 {
	return MESSAGE_ORDER_EXECUTED_PRICE
}

func (o OrderExecutedPrice) Bytes() []byte {
	data := make([]byte, orderExecutedPriceSize)

	data[0] = MESSAGE_ORDER_EXECUTED_PRICE
	binary.BigEndian.PutUint16(data[1:3], o.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(o.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], o.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], o.Reference)
	binary.BigEndian.PutUint32(data[19:23], o.Shares)
	binary.BigEndian.PutUint64(data[23:31], o.MatchNumber)

	if o.Printable {
		data[31] = 'Y'
	} else {
		data[31] = 'N'
	}

	binary.BigEndian.PutUint32(data[32:], o.ExecutionPrice)

	return data
}

func ParseOrderExecuted(data []byte) (OrderExecuted, error) {
	if len(data) != orderExecutedSize {
		return OrderExecuted{}, NewInvalidPacketSize(orderExecutedSize, len(data))
	}

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
	}, nil
}

func ParseOrderExecutedPrice(data []byte) (OrderExecutedPrice, error) {
	if len(data) != orderExecutedPriceSize {
		return OrderExecutedPrice{}, NewInvalidPacketSize(orderExecutedPriceSize, len(data))
	}

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
	}, nil
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
