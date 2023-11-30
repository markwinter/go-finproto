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
	// TODO: implement
	return data
}

func ParseMwcbStatus(data []byte) MwcbStatus {
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
	}
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
