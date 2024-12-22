/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/quagmt/udecimal"
)

type MwcbLevel struct {
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	LevelOne       udecimal.Decimal // Price (8)
	LevelTwo       udecimal.Decimal // Price (8)
	LevelThree     udecimal.Decimal // Price (8)
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

	levelOne, _ := priceToBytes(m.LevelOne, 8)
	levelTwo, _ := priceToBytes(m.LevelTwo, 8)
	levelThree, _ := priceToBytes(m.LevelThree, 8)

	copy(data[11:19], levelOne)
	copy(data[19:27], levelTwo)
	copy(data[27:], levelThree)

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

	levelOne, _ := bytesToPrice(data[11:19], 8)
	levelTwo, _ := bytesToPrice(data[19:27], 8)
	levelThree, _ := bytesToPrice(data[27:], 8)

	return MwcbLevel{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		LevelOne:       levelOne,
		LevelTwo:       levelTwo,
		LevelThree:     levelThree,
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
		l.LevelOne,
		l.LevelTwo,
		l.LevelThree,
	)
}
