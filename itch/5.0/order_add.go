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
	// TODO: implement
	return data
}

func ParseOrderAdd(data []byte) OrderAdd {
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
	}
}

func ParseOrderAddAttributed(data []byte) OrderAddAttributed {
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
	}
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
