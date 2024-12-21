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

type OrderIndicator uint8

const (
	ORDER_INDICATOR_BUY  OrderIndicator = 'B'
	ORDER_INDICATOR_SELL OrderIndicator = 'S'
)

type OrderAdd struct {
	Stock          string
	Timestamp      time.Duration
	Reference      uint64
	Shares         uint32
	Price          uint32
	StockLocate    uint16
	TrackingNumber uint16
	OrderIndicator OrderIndicator
}

func (o OrderAdd) Type() uint8 {
	return MESSAGE_ORDER_ADD
}

func (o OrderAdd) Bytes() []byte {
	data := make([]byte, orderAddSize)
	// TODO: implement

	data[0] = MESSAGE_ORDER_ADD
	binary.BigEndian.PutUint16(data[1:3], o.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(o.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], o.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], o.Reference)
	data[19] = byte(o.OrderIndicator)

	binary.BigEndian.PutUint32(data[20:24], o.Shares)

	copy(data[24:32], []byte(fmt.Sprintf("%-8s", o.Stock)))

	binary.BigEndian.PutUint32(data[32:], o.Price)

	return data
}

type OrderAddAttributed struct {
	Stock          string
	Attribution    string
	Timestamp      time.Duration
	Reference      uint64
	Shares         uint32
	Price          uint32
	StockLocate    uint16
	TrackingNumber uint16
	OrderIndicator OrderIndicator
}

func (o OrderAddAttributed) Type() uint8 {
	return MESSAGE_ORDER_ADD_ATTRIBUTED
}

func (o OrderAddAttributed) Bytes() []byte {
	data := make([]byte, orderAddAttrSize)

	data[0] = MESSAGE_ORDER_ADD_ATTRIBUTED
	binary.BigEndian.PutUint16(data[1:3], o.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(o.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], o.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], o.Reference)
	data[19] = byte(o.OrderIndicator)

	binary.BigEndian.PutUint32(data[20:24], o.Shares)

	copy(data[24:32], []byte(fmt.Sprintf("%-8s", o.Stock)))

	binary.BigEndian.PutUint32(data[32:36], o.Price)

	copy(data[36:], []byte(fmt.Sprintf("%-4s", o.Attribution)))

	return data
}

func ParseOrderAdd(data []byte) (OrderAdd, error) {
	if len(data) != orderAddSize {
		return OrderAdd{}, NewInvalidPacketSize(orderAddSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return OrderAdd{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Reference:      binary.BigEndian.Uint64(data[11:19]),
		OrderIndicator: OrderIndicator(data[19]),
		Shares:         binary.BigEndian.Uint32(data[20:24]),
		Stock:          strings.TrimSpace(string(data[24:32])),
		Price:          binary.BigEndian.Uint32(data[32:]),
	}, nil
}

func ParseOrderAddAttributed(data []byte) (OrderAddAttributed, error) {
	if len(data) != orderAddAttrSize {
		return OrderAddAttributed{}, NewInvalidPacketSize(orderAddAttrSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return OrderAddAttributed{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Reference:      binary.BigEndian.Uint64(data[11:19]),
		OrderIndicator: OrderIndicator(data[19]),
		Shares:         binary.BigEndian.Uint32(data[20:24]),
		Stock:          strings.TrimSpace(string(data[24:32])),
		Price:          binary.BigEndian.Uint32(data[32:36]),
		Attribution:    strings.TrimSpace(string(data[36:])),
	}, nil
}

func (a OrderAdd) String() string {
	return fmt.Sprintf("[Order Add]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Reference: %v\n"+
		"Order Indicator: %v\n"+
		"Shares: %v\n"+
		"Stock: %v\n"+
		"Price: %v\n",
		a.StockLocate, a.TrackingNumber, a.Timestamp,
		a.Reference, a.OrderIndicator, a.Shares, a.Stock,
		float64(a.Price)/10000,
	)
}

func (a OrderAddAttributed) String() string {
	return fmt.Sprintf("[Order Add Attributed]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Reference: %v\n"+
		"Order Indicator: %v\n"+
		"Shares: %v\n"+
		"Stock: %v\n"+
		"Price: %v\n"+
		"Attribution: %v\n",
		a.StockLocate, a.TrackingNumber, a.Timestamp,
		a.Reference, a.OrderIndicator, a.Shares, a.Stock,
		float64(a.Price)/10000, a.Attribution,
	)
}

func (o OrderIndicator) String() string {
	switch o {
	case ORDER_INDICATOR_BUY:
		return "Buy"
	case ORDER_INDICATOR_SELL:
		return "Sell"
	}

	return "Unknown OrderIndicator"
}
