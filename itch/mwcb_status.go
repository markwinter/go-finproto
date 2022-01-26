/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type McwbStatus struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	BreachedLevel  uint8
}

func MakeMcwbStatus(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return McwbStatus{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		BreachedLevel:  data[11],
	}
}

func (l McwbStatus) String() string {
	return fmt.Sprintf("[MWCB Status]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Breached Level: %v\n",
		l.StockLocate, l.TrackingNumber, l.Timestamp,
		l.BreachedLevel,
	)
}
