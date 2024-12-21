/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type OrderDelete struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Reference      uint64
}

func (o OrderDelete) Type() uint8 {
	return MESSAGE_ORDER_DELETE
}

func (o OrderDelete) Bytes() []byte {
	data := make([]byte, orderDeleteSize)

	data[0] = MESSAGE_ORDER_DELETE
	binary.BigEndian.PutUint16(data[1:3], o.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(o.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], o.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:], o.Reference)

	return data
}

func ParseOrderDelete(data []byte) (OrderDelete, error) {
	if len(data) != orderDeleteSize {
		return OrderDelete{}, NewInvalidPacketSize(orderDeleteSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return OrderDelete{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Reference:      binary.BigEndian.Uint64(data[11:19]),
	}, nil
}

func (o OrderDelete) String() string {
	return fmt.Sprintf("[Order Delete]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Reference: %v\n",
		o.StockLocate, o.TrackingNumber, o.Timestamp,
		o.Reference,
	)
}
