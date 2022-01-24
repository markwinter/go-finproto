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

type LuldCollar struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Stock          string
	ReferencePrice uint32
	UpperPrice     uint32
	LowerPrice     uint32
	Extension      uint32
}

func MakeLuldCollar(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return LuldCollar{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Stock:          strings.TrimSpace(string(data[11:19])),
		ReferencePrice: binary.BigEndian.Uint32(data[19:23]),
		UpperPrice:     binary.BigEndian.Uint32(data[23:27]),
		LowerPrice:     binary.BigEndian.Uint32(data[27:31]),
		Extension:      binary.BigEndian.Uint32(data[31:]),
	}
}

func (l LuldCollar) String() string {
	return fmt.Sprintf("[IPO Quotation]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Stock: %v\n"+
		"Reference Price: %v\n"+
		"Upper Price: %v\n"+
		"Lower Price: %v\n"+
		"Extension: %v\n",
		l.StockLocate, l.TrackingNumber, l.Timestamp,
		l.Stock,
		float64(l.ReferencePrice)/10000,
		float64(l.UpperPrice)/10000,
		float64(l.LowerPrice)/10000,
		l.Extension,
	)
}
