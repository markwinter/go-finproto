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

func (o OrderReplace) Type() uint8 {
	return MESSAGE_ORDER_REPLACE
}

func (o OrderReplace) Bytes() []byte {
	data := make([]byte, orderReplaceSize)

	data[0] = MESSAGE_ORDER_REPLACE
	binary.BigEndian.PutUint16(data[1:3], o.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(o.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], o.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], o.OriginalReference)
	binary.BigEndian.PutUint64(data[19:27], o.NewReference)
	binary.BigEndian.PutUint32(data[27:31], o.Shares)
	binary.BigEndian.PutUint32(data[31:], o.Price)

	return data
}

func ParseOrderReplace(data []byte) (OrderReplace, error) {
	if len(data) != orderReplaceSize {
		return OrderReplace{}, NewInvalidPacketSize(orderReplaceSize, len(data))
	}

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
	}, nil
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
