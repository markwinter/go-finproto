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

type AddOrder struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Reference      uint64
	OrderIndicator OrderIndicator
	Shares         uint32
	Stock          string
	Price          uint32
}

type AddOrderAttributed struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Reference      uint64
	OrderIndicator OrderIndicator
	Shares         uint32
	Stock          string
	Price          uint32
	Attribution    string
}

func MakeAddOrder(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return AddOrder{
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

func MakeAddOrderAttributed(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return AddOrderAttributed{
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

func (a AddOrder) String() string {
	return fmt.Sprintf("[Add Order]\n"+
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

func (a AddOrderAttributed) String() string {
	return fmt.Sprintf("[Add Order]\n"+
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
