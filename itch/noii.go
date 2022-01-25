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

type ImbalanceDirection uint8

const (
	IMBALANCE_BUY          ImbalanceDirection = 'B'
	IMBALANCE_SELL         ImbalanceDirection = 'S'
	IMBALANCE_NONE         ImbalanceDirection = 'N'
	IMBALANCE_INSUFFICIENT ImbalanceDirection = 'O'
)

type Noii struct {
	StockLocate        uint16
	TrackingNumber     uint16
	Timestamp          time.Duration
	PairedShares       uint64
	ImbalanceShares    uint64
	ImbalanceDirection ImbalanceDirection
	Stock              string
	FarPrice           uint32
	NearPrice          uint32
	CurrentPrice       uint32
	CrossType          CrossType
	VariationIndicator uint8
}

func MakeNoii(data []byte) Message {
	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	return Noii{
		StockLocate:        locate,
		TrackingNumber:     tracking,
		Timestamp:          time.Duration(t),
		PairedShares:       binary.BigEndian.Uint64(data[11:19]),
		ImbalanceShares:    binary.BigEndian.Uint64(data[19:27]),
		ImbalanceDirection: ImbalanceDirection(data[27]),
		Stock:              strings.TrimSpace(string(data[28:36])),
		FarPrice:           binary.BigEndian.Uint32(data[36:40]),
		NearPrice:          binary.BigEndian.Uint32(data[40:44]),
		CurrentPrice:       binary.BigEndian.Uint32(data[44:48]),
		CrossType:          CrossType(data[48]),
		VariationIndicator: data[49],
	}
}

func (n Noii) String() string {
	return fmt.Sprintf("[NOII]\n"+
		"Stock Locate: %v\n"+
		"Tracking Number: %v\n"+
		"Timestamp: %v\n"+
		"Paired Shares: %v\n"+
		"Imbalance Shares: %v\n"+
		"Imbalance Direction: %v\n"+
		"Stock: %v\n"+
		"Far Price: %v\n"+
		"Near Price: %v\n"+
		"Current Price: %v\n"+
		"Cross Type: %v\n"+
		"Variation Indicator: %v\n",
		n.StockLocate, n.TrackingNumber, n.Timestamp,
		n.PairedShares, n.ImbalanceShares, n.ImbalanceDirection,
		n.Stock, float64(n.FarPrice)/10000, float64(n.NearPrice)/10000,
		float64(n.CurrentPrice)/10000, n.CrossType, n.VariationIndicator,
	)
}

func (i ImbalanceDirection) String() string {
	switch i {
	case IMBALANCE_BUY:
		return "Buy"
	case IMBALANCE_SELL:
		return "Sell"
	case IMBALANCE_NONE:
		return "None"
	case IMBALANCE_INSUFFICIENT:
		return "Insufficient orders to calculate"
	}

	return "Unknown ImbalanceDirection"
}
