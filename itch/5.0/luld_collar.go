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

func (l LuldCollar) Type() uint8 {
	return MESSAGE_LULD_COLLAR
}

func (l LuldCollar) Bytes() []byte {
	data := make([]byte, luldSize)

	data[0] = MESSAGE_LULD_COLLAR
	binary.BigEndian.PutUint16(data[1:3], l.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(l.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], l.TrackingNumber)

	copy(data[11:19], []byte(fmt.Sprintf("%-8s", l.Stock)))

	binary.BigEndian.PutUint32(data[19:23], l.ReferencePrice)
	binary.BigEndian.PutUint32(data[23:27], l.UpperPrice)
	binary.BigEndian.PutUint32(data[27:31], l.LowerPrice)
	binary.BigEndian.PutUint32(data[31:], l.Extension)

	return data
}

func ParseLuldCollar(data []byte) (LuldCollar, error) {
	if len(data) != luldSize {
		return LuldCollar{}, NewInvalidPacketSize(luldSize, len(data))
	}

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
	}, nil
}

func (l LuldCollar) String() string {
	return fmt.Sprintf("[LULD Collar]\n"+
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
