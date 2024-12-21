/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type TradeBroken struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	MatchNumber    uint64
}

func (t TradeBroken) Type() uint8 {
	return MESSAGE_TRADE_BROKEN
}

func (t TradeBroken) Bytes() []byte {
	data := make([]byte, tradeBrokenSize)

	data[0] = MESSAGE_TRADE_BROKEN
	binary.BigEndian.PutUint16(data[1:3], t.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(t.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], t.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:], t.MatchNumber)

	return data
}

func ParseTradeBroken(data []byte) (TradeBroken, error) {
	if len(data) != tradeBrokenSize {
		return TradeBroken{}, NewInvalidPacketSize(tradeBrokenSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return TradeBroken{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		MatchNumber:    binary.BigEndian.Uint64(data[11:]),
	}, nil
}

func (o TradeBroken) String() string {
	return fmt.Sprintf("[Trade Broken]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Match Number: %v\n",
		o.StockLocate, o.TrackingNumber, o.Timestamp,
		o.MatchNumber,
	)
}
