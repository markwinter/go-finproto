/*
 * Copyright (c) 2022 Mark Edward Winter
 */
package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type OrderCancel struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Reference      uint64
	Shares         uint32
}

func MakeOrderCancel(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return OrderCancel{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Reference:      binary.BigEndian.Uint64(data[11:19]),
		Shares:         binary.BigEndian.Uint32(data[19:23]),
	}
}

func (o OrderCancel) String() string {
	return fmt.Sprintf("[Order Cancelled]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Reference: %v\n"+
		"Shares: %v\n",
		o.StockLocate, o.TrackingNumber, o.Timestamp,
		o.Reference, o.Shares,
	)
}
