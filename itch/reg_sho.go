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
	StockLocate    uint16
	TrackingNumber uint16
	Timestamp      time.Duration
	Stock          string
	Action         RegShoAction
}

func MakeRegSho(data []byte) Message {
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
	}
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
