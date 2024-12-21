/*
 * Copyright (c) 2022 Mark Edward Winter
 */

package itch

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/quagmt/udecimal"
)

type ImbalanceDirection uint8

const (
	IMBALANCE_BUY          ImbalanceDirection = 'B'
	IMBALANCE_SELL         ImbalanceDirection = 'S'
	IMBALANCE_NONE         ImbalanceDirection = 'N'
	IMBALANCE_INSUFFICIENT ImbalanceDirection = 'O'
)

type Noii struct {
	Stock              string
	Timestamp          time.Duration
	PairedShares       uint64
	ImbalanceShares    uint64
	FarPrice           udecimal.Decimal // Price (4)
	NearPrice          udecimal.Decimal // Price (4)
	CurrentPrice       udecimal.Decimal // Price (4)
	StockLocate        uint16
	TrackingNumber     uint16
	ImbalanceDirection ImbalanceDirection
	CrossType          CrossType
	VariationIndicator uint8
}

func (n Noii) Type() uint8 {
	return MESSAGE_NOII
}

func (n Noii) Bytes() []byte {
	data := make([]byte, noiiSize)

	data[0] = MESSAGE_NOII

	binary.BigEndian.PutUint16(data[1:3], n.StockLocate)

	// Order of these fields are important. We write timestamp to 3:11 first to let us write a uint64, then overwrite 3:5 with tracking number
	binary.BigEndian.PutUint64(data[3:11], uint64(n.Timestamp.Nanoseconds()))
	binary.BigEndian.PutUint16(data[3:5], n.TrackingNumber)

	binary.BigEndian.PutUint64(data[11:19], n.PairedShares)
	binary.BigEndian.PutUint64(data[19:27], n.ImbalanceShares)
	data[27] = byte(n.ImbalanceDirection)

	copy(data[28:36], []byte(fmt.Sprintf("%-8s", n.Stock)))

	farP, _ := priceToBytes(n.FarPrice, 4)
	nearP, _ := priceToBytes(n.NearPrice, 4)
	curP, _ := priceToBytes(n.CurrentPrice, 4)

	copy(data[36:40], farP)
	copy(data[40:44], nearP)
	copy(data[44:48], curP)

	data[48] = byte(n.CrossType)
	data[49] = n.VariationIndicator

	return data
}

func ParseNoii(data []byte) (Noii, error) {
	if len(data) != noiiSize {
		return Noii{}, NewInvalidPacketSize(noiiSize, len(data))
	}

	locate := binary.BigEndian.Uint16(data[1:3])
	tracking := binary.BigEndian.Uint16(data[3:5])
	data[3] = 0
	data[4] = 0
	t := binary.BigEndian.Uint64(data[3:11])

	farP, _ := bytesToPrice(data[36:40], 4)
	nearP, _ := bytesToPrice(data[40:44], 4)
	curP, _ := bytesToPrice(data[44:48], 4)

	return Noii{
		StockLocate:        locate,
		TrackingNumber:     tracking,
		Timestamp:          time.Duration(t),
		PairedShares:       binary.BigEndian.Uint64(data[11:19]),
		ImbalanceShares:    binary.BigEndian.Uint64(data[19:27]),
		ImbalanceDirection: ImbalanceDirection(data[27]),
		Stock:              strings.TrimSpace(string(data[28:36])),
		FarPrice:           farP,
		NearPrice:          nearP,
		CurrentPrice:       curP,
		CrossType:          CrossType(data[48]),
		VariationIndicator: data[49],
	}, nil
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
		n.Stock, n.FarPrice, n.NearPrice,
		n.CurrentPrice, n.CrossType, n.VariationIndicator,
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
