/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"time"
)

type MwcbLevel struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	LevelOne       uint64
	LevelTwo       uint64
	LevelThree     uint64
}

func (m MwcbLevel) Type() uint8 {
	return MESSAGE_MWCB_LEVEL
}

func (m MwcbLevel) Bytes() []byte {
	data := make([]byte, mwcbLevelSize)

	data[0] = MESSAGE_MWCB_LEVEL
	binary.BigEndian.PutUint16(data[1:3], m.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(m.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], m.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], m.LevelOne)
	binary.BigEndian.PutUint64(data[19:27], m.LevelTwo)
	binary.BigEndian.PutUint64(data[27:], m.LevelThree)

	return data
}

func ParseMwcbLevel(data []byte) (MwcbLevel, error) {
	if len(data) != mwcbLevelSize {
		return MwcbLevel{}, NewInvalidPacketSize(mwcbLevelSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return MwcbLevel{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		LevelOne:       binary.BigEndian.Uint64(data[11:19]),
		LevelTwo:       binary.BigEndian.Uint64(data[19:27]),
		LevelThree:     binary.BigEndian.Uint64(data[27:]),
	}, nil
}

func (l MwcbLevel) String() string {
	return fmt.Sprintf("[MWCB Levels]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Level One: %v\n"+
		"Level Two: %v\n"+
		"Level Three: %v\n",
		l.StockLocate, l.TrackingNumber, l.Timestamp,
		float64(l.LevelOne)/100000000,
		float64(l.LevelTwo)/100000000,
		float64(l.LevelThree)/100000000,
	)
}
