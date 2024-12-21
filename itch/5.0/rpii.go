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

type RpiInterestFlag uint8

const (
	RPI_INTEREST_BUY  RpiInterestFlag = 'B'
	RPI_INTEREST_SELL RpiInterestFlag = 'S'
	RPI_INTEREST_BOTH RpiInterestFlag = 'A'
	RPI_INTEREST_NONE RpiInterestFlag = 'N'
)

type Rpii struct {
	Stock          string
	Timestamp      time.Duration
	StockLocate    uint16
	TrackingNumber uint16
	InterestFlag   RpiInterestFlag
}

func (r Rpii) Type() uint8 {
	return MESSAGE_RPII
}

func (r Rpii) Bytes() []byte {
	data := make([]byte, rpiiSize)
	// TODO: implement
	return data
}

func ParseRpii(data []byte) (Rpii, error) {
	if len(data) != rpiiSize {
		return Rpii{}, NewInvalidPacketSize(rpiiSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return Rpii{
		StockLocate:    locate,
		TrackingNumber: tracking,
		Timestamp:      time.Duration(t),
		Stock:          strings.TrimSpace(string(data[11:19])),
		InterestFlag:   RpiInterestFlag(data[19]),
	}, nil
}

func (n Rpii) String() string {
	return fmt.Sprintf("[RPII]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Stock: %v\n"+
		"RPI Interest: %v\n",
		n.StockLocate, n.TrackingNumber, n.Timestamp,
		n.Stock, n.InterestFlag,
	)
}

func (i RpiInterestFlag) String() string {
	switch i {
	case RPI_INTEREST_BUY:
		return "Buy"
	case RPI_INTEREST_SELL:
		return "Sell"
	case RPI_INTEREST_BOTH:
		return "Both sides"
	case RPI_INTEREST_NONE:
		return "No orders available"
	}

	return "Unknown RpiInterestFlag"
}
