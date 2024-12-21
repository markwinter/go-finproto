/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type MwcbStatus struct {
	Timestamp      time.Duration
	StockLocate    uint16
	TrackingNumber uint16
	BreachedLevel  uint8
}

func (m MwcbStatus) Type() uint8 {
	return MESSAGE_MWCB_STATUS
}

func (m MwcbStatus) Bytes() []byte {
	data := make([]byte, mwcbStatusSize)

	data[0] = MESSAGE_MWCB_STATUS
	binary.BigEndian.PutUint16(data[1:3], m.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(m.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], m.TrackingNumber)

	data[11] = m.BreachedLevel

	return data
}

func ParseMwcbStatus(data []byte) (MwcbStatus, error) {
	if len(data) != mwcbStatusSize {
		return MwcbStatus{}, NewInvalidPacketSize(mwcbStatusSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return MwcbStatus{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		BreachedLevel:  data[11],
	}, nil
}

func (l MwcbStatus) String() string {
	return fmt.Sprintf("[MWCB Status]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Breached Level: %v\n",
		l.StockLocate, l.TrackingNumber, l.Timestamp,
		l.BreachedLevel,
	)
}
