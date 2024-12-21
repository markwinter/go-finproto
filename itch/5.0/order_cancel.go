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
	Timestamp      time.Duration
	Reference      uint64
	Shares         uint32
	StockLocate    uint16
	TrackingNumber uint16
}

func (o OrderCancel) Type() uint8 {
	return MESSAGE_ORDER_CANCEL
}

func (o OrderCancel) Bytes() []byte {
	data := make([]byte, orderCancelSize)

	data[0] = MESSAGE_ORDER_CANCEL
	binary.BigEndian.PutUint16(data[1:3], o.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(o.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], o.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], o.Reference)
	binary.BigEndian.PutUint32(data[19:], o.Shares)

	return data
}

func ParseOrderCancel(data []byte) (OrderCancel, error) {
	if len(data) != orderCancelSize {
		return OrderCancel{}, NewInvalidPacketSize(orderCancelSize, len(data))
	}

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
	}, nil
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
