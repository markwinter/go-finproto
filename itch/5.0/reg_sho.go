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

type RegShoAction uint8

const (
	REGSHO_NO_PRICE_TEST RegShoAction = '0'
	REGSHO_INTRADAY_DROP RegShoAction = '1'
	REGSHO_REMAINS       RegShoAction = '2'
)

type RegSho struct {
	Stock          string
	Timestamp      time.Duration
	StockLocate    uint16
	TrackingNumber uint16
	Action         RegShoAction
}

func (r RegSho) Type() uint8 {
	return MESSAGE_REG_SHO
}

func (r RegSho) Bytes() []byte {
	data := make([]byte, regShoSize)

	data[0] = MESSAGE_REG_SHO
	binary.BigEndian.PutUint16(data[1:3], r.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(r.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], r.TrackingNumber)

	copy(data[11:19], []byte(fmt.Sprintf("%-8s", r.Stock)))

	data[19] = byte(r.Action)

	return data
}

func ParseRegSho(data []byte) (RegSho, error) {
	if len(data) != regShoSize {
		return RegSho{}, NewInvalidPacketSize(regShoSize, len(data))
	}

	fmt.Printf("%v\n", data)

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return RegSho{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Stock:          strings.TrimSpace(string(data[11:19])),
		Action:         RegShoAction(data[19]),
	}, nil
}

func (r RegSho) String() string {
	return fmt.Sprintf("[Reg SHO]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Stock: %v\n"+
		"Action: %v\n",
		r.StockLocate, r.TrackingNumber, r.Timestamp,
		r.Stock, r.Action,
	)
}

func (a RegShoAction) String() string {
	switch a {
	case REGSHO_NO_PRICE_TEST:
		return "No price test in place"
	case REGSHO_INTRADAY_DROP:
		return "Restriction in place due to intraday price drop"
	case REGSHO_REMAINS:
		return "Restriction remains in effect"
	}

	return "Unknown RegShoAction"
}
