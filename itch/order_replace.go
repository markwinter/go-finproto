/*
 * Copyright (c) 2022 Mark Edward Winter
 */
package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type OrderReplace struct {
	StockLocate       uint16
	TrackingNumber    uint16
	Timestamp         time.Duration
	OriginalReference uint64
	NewReference      uint64
	Shares            uint32
	Price             uint32
}

func MakeOrderReplace(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return OrderReplace{
		StockLocate:       locate,
		TrackingNumber:    tracking,
		Timestamp:         time.Duration(t),
		OriginalReference: binary.BigEndian.Uint64(data[11:19]),
		NewReference:      binary.BigEndian.Uint64(data[19:27]),
		Shares:            binary.BigEndian.Uint32(data[27:31]),
		Price:             binary.BigEndian.Uint32(data[31:]),
	}
}

func (o OrderReplace) String() string {
	return fmt.Sprintf("[Order Replaced]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Original Reference: %v\n"+
		"New Reference: %v\n"+
		"Shares: %v\n"+
		"Price: %v\n",
		o.StockLocate, o.TrackingNumber, o.Timestamp,
		o.OriginalReference, o.NewReference,
		o.Shares, float64(o.Price)/10000,
	)
}
